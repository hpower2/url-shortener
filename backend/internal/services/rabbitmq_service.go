package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hpower2/url-shortener/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// EmailMessage represents an email message in the queue
type EmailMessage struct {
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
	Type       string `json:"type"` // "otp" or "welcome"
	OTPCode    string `json:"otp_code,omitempty"`
	Purpose    string `json:"purpose,omitempty"`
	Retry      int    `json:"retry"`
	MaxRetries int    `json:"max_retries"`
}

// RabbitMQService interface defines the contract for RabbitMQ operations
type RabbitMQService interface {
	Connect() error
	Close() error
	PublishEmail(message *EmailMessage) error
	ConsumeEmails(handler func(*EmailMessage) error) error
	PublishDelayedEmail(message *EmailMessage, delay time.Duration) error
}

// rabbitMQService implements RabbitMQService interface
type rabbitMQService struct {
	config     *config.RabbitMQConfig
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewRabbitMQService creates a new RabbitMQ service
func NewRabbitMQService(config *config.RabbitMQConfig) RabbitMQService {
	return &rabbitMQService{
		config: config,
	}
}

// Connect establishes connection to RabbitMQ
func (s *rabbitMQService) Connect() error {
	var err error

	// Construct connection URL
	url := s.config.URL
	if url == "" {
		url = fmt.Sprintf("amqp://%s:%s@%s:%s/",
			s.config.Username, s.config.Password, s.config.Host, s.config.Port)
	}

	// Connect to RabbitMQ
	s.connection, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create channel
	s.channel, err = s.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	// Declare email queue
	_, err = s.channel.QueueDeclare(
		"email_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare email queue: %w", err)
	}

	// Declare delayed email queue
	_, err = s.channel.QueueDeclare(
		"email_delay_queue", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		amqp.Table{
			"x-message-ttl":             30000, // 30 seconds TTL
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": "email_queue",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare delayed email queue: %w", err)
	}

	log.Println("Connected to RabbitMQ successfully")
	return nil
}

// Close closes the RabbitMQ connection
func (s *rabbitMQService) Close() error {
	if s.channel != nil {
		if err := s.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	if s.connection != nil {
		if err := s.connection.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}
	return nil
}

// PublishEmail publishes an email message to the queue
func (s *rabbitMQService) PublishEmail(message *EmailMessage) error {
	if s.channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	// Set default values
	if message.MaxRetries == 0 {
		message.MaxRetries = 3
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = s.channel.Publish(
		"",            // exchange
		"email_queue", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Email message published to queue: %s", message.To)
	return nil
}

// PublishDelayedEmail publishes an email message with a delay
func (s *rabbitMQService) PublishDelayedEmail(message *EmailMessage, delay time.Duration) error {
	if s.channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = s.channel.Publish(
		"",                  // exchange
		"email_delay_queue", // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Expiration:   fmt.Sprintf("%d", delay.Milliseconds()),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish delayed message: %w", err)
	}

	log.Printf("Delayed email message published (delay: %v): %s", delay, message.To)
	return nil
}

// ConsumeEmails consumes email messages from the queue
func (s *rabbitMQService) ConsumeEmails(handler func(*EmailMessage) error) error {
	if s.channel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	// Set QoS to process one message at a time
	err := s.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := s.channel.Consume(
		"email_queue", // queue
		"",            // consumer
		false,         // auto-ack (we'll manually ack)
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Println("Starting email queue consumer...")

	for msg := range msgs {
		var emailMsg EmailMessage
		if err := json.Unmarshal(msg.Body, &emailMsg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			msg.Nack(false, false) // Reject message
			continue
		}

		log.Printf("Processing email message: %s", emailMsg.To)

		// Handle the message
		if err := handler(&emailMsg); err != nil {
			log.Printf("Failed to handle email message: %v", err)

			// Increment retry count
			emailMsg.Retry++

			// If max retries reached, reject the message
			if emailMsg.Retry >= emailMsg.MaxRetries {
				log.Printf("Max retries reached for email to %s, rejecting message", emailMsg.To)
				msg.Nack(false, false) // Reject without requeue
				continue
			}

			// Publish to delayed queue for retry
			delay := time.Duration(emailMsg.Retry*30) * time.Second // Exponential backoff
			if err := s.PublishDelayedEmail(&emailMsg, delay); err != nil {
				log.Printf("Failed to publish retry message: %v", err)
			} else {
				log.Printf("Scheduled retry %d/%d for email to %s (delay: %v)",
					emailMsg.Retry, emailMsg.MaxRetries, emailMsg.To, delay)
			}

			msg.Ack(false) // Acknowledge original message
		} else {
			log.Printf("Email message processed successfully: %s", emailMsg.To)
			msg.Ack(false) // Acknowledge successful processing
		}
	}

	return nil
}
