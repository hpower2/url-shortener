package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents different types of errors
type ErrorCode string

const (
	// Client errors
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeInactive      ErrorCode = "URL_INACTIVE"
	ErrCodeExpired       ErrorCode = "URL_EXPIRED"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrCodeRateLimit     ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrCodeBadRequest    ErrorCode = "BAD_REQUEST"
	
	// Server errors
	ErrCodeInternal      ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase      ErrorCode = "DATABASE_ERROR"
	ErrCodeRedis         ErrorCode = "REDIS_ERROR"
	ErrCodeExternal      ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeTimeout       ErrorCode = "TIMEOUT_ERROR"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"-"`
	Err        error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// Predefined error constructors
func NewValidationError(message string, err error) *AppError {
	return NewAppError(ErrCodeValidation, message, http.StatusBadRequest, err)
}

func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(ErrCodeNotFound, message, http.StatusNotFound, err)
}

func NewInactiveError(message string, err error) *AppError {
	return NewAppError(ErrCodeInactive, message, http.StatusGone, err)
}

func NewExpiredError(message string, err error) *AppError {
	return NewAppError(ErrCodeExpired, message, http.StatusGone, err)
}

func NewAlreadyExistsError(message string, err error) *AppError {
	return NewAppError(ErrCodeAlreadyExists, message, http.StatusConflict, err)
}

func NewUnauthorizedError(message string, err error) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, http.StatusUnauthorized, err)
}

func NewForbiddenError(message string, err error) *AppError {
	return NewAppError(ErrCodeForbidden, message, http.StatusForbidden, err)
}

func NewRateLimitError(message string, err error) *AppError {
	return NewAppError(ErrCodeRateLimit, message, http.StatusTooManyRequests, err)
}

func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(ErrCodeBadRequest, message, http.StatusBadRequest, err)
}

func NewInternalError(message string, err error) *AppError {
	return NewAppError(ErrCodeInternal, message, http.StatusInternalServerError, err)
}

func NewDatabaseError(message string, err error) *AppError {
	return NewAppError(ErrCodeDatabase, message, http.StatusInternalServerError, err)
}

func NewRedisError(message string, err error) *AppError {
	return NewAppError(ErrCodeRedis, message, http.StatusInternalServerError, err)
}

func NewExternalServiceError(message string, err error) *AppError {
	return NewAppError(ErrCodeExternal, message, http.StatusBadGateway, err)
}

func NewTimeoutError(message string, err error) *AppError {
	return NewAppError(ErrCodeTimeout, message, http.StatusRequestTimeout, err)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from error
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}

// ErrorResponse represents the JSON error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents error details in response
type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// ToErrorResponse converts AppError to ErrorResponse
func (e *AppError) ToErrorResponse() ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		},
	}
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (v ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed with %d errors", len(v.Errors))
}

// NewValidationErrors creates new validation errors
func NewValidationErrors(errors []ValidationError) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    "Validation failed",
		StatusCode: http.StatusBadRequest,
		Details:    fmt.Sprintf("%d validation errors", len(errors)),
	}
} 