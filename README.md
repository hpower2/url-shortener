# URL Shortener 

## THIS CODE 95% GENERATED BY CURSOR PROMPT - CLAUDE-4-SONNET

- Sample Dashboard : https://short.irvineafri.com
- Sample Shortener (custom link) : https://s.iafri.com/japan-song-2
- Sample Generated Links : https://s.iafri.com/nch4OOA2

A modern, full-stack URL shortener application with user authentication, analytics, and QR code generation.

## 🚀 Features

- **🔐 User Authentication** - JWT-based secure registration and login
- **⚡ Fast URL Shortening** - Generate short links with optional custom codes
- **📊 Advanced Analytics** - Track clicks, geographic data, and referrer information
- **📱 QR Code Generation** - Create QR codes for shortened URLs
- **⏰ URL Expiration** - Set expiration dates for links
- **🎯 User Dashboard** - Manage all your shortened URLs in one place
- **📈 Real-time Metrics** - View detailed click statistics and analytics
- **🔒 User Isolation** - Each user can only access their own URLs
- **🚀 High Performance** - Redis caching for fast redirects

## 🛠️ Tech Stack

### Backend
- **Go 1.24+** with Gin framework
- **PostgreSQL** for data persistence
- **Redis** for caching and performance
- **JWT** for authentication
- **RabbitMQ** for email queuing (optional)
- **Docker** for containerization

### Frontend
- **React 18** with TypeScript
- **Vite** for build tooling
- **Tailwind CSS** for styling
- **React Router** for navigation
- **React Query** for state management
- **Heroicons** for icons
- **Recharts** for analytics visualization

## 📁 Project Structure

```
url-shortener/
├── backend/                 # Go API server
│   ├── cmd/                # Application entry point
│   ├── internal/           # Internal packages
│   │   ├── handlers/       # HTTP handlers
│   │   ├── services/       # Business logic
│   │   ├── repository/     # Data access layer
│   │   ├── models/         # Data models
│   │   ├── middleware/     # HTTP middleware
│   │   └── config/         # Configuration
│   ├── migrations/         # Database migrations
│   ├── env.example         # Environment variables template
│   └── Dockerfile          # Backend container
├── frontend/               # React application
│   ├── src/
│   │   ├── components/     # Reusable components
│   │   ├── pages/          # Page components
│   │   ├── contexts/       # React contexts
│   │   └── services/       # API services
│   └── Dockerfile          # Frontend container
└── docker-compose.yml      # Application containers
```

## 🚀 Quick Start

### Prerequisites

Before running the application, you need to set up the following services:

1. **PostgreSQL Database**
   - Create a database for the URL shortener
   - Run the migration scripts from `backend/migrations/`

2. **Redis Server**
   - Install and run Redis server
   - Default configuration should work for development

3. **RabbitMQ (Optional)**
   - Only required if you want email notifications
   - Can be skipped for basic functionality

### Setup Instructions

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd url-shortener
   ```

2. **Set up your database**
   ```bash
   # Create PostgreSQL database
   createdb urlshortener
   
   # Run migrations
   psql urlshortener < backend/migrations/001_initial_schema.sql
   psql urlshortener < backend/migrations/003_add_otp_and_link_limits.sql
   ```

3. **Configure environment variables**
   ```bash
   # Copy the example environment file
   cp backend/env.example backend/.env
   
   # Edit backend/.env with your database and Redis settings
   nano backend/.env
   ```

4. **Start the application**
   ```bash
   # Build and run both frontend and backend
   docker-compose up -d
   ```

5. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:15522

## ⚙️ Configuration

### Required Environment Variables

Edit `backend/.env` with your settings:

```bash
# Server Configuration
SERVER_PORT=8080

# Database Configuration (Required)
DB_HOST=your-database-host        # e.g., localhost or your DB server IP
DB_PORT=5432
DB_USER=your-db-username
DB_PASSWORD=your-db-password
DB_NAME=your-database-name
DB_SSLMODE=disable               # or require for production

# Redis Configuration (Required)
REDIS_HOST=your-redis-host       # e.g., localhost or your Redis server IP
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password  # leave empty if no password

# Application Configuration
BASE_URL=http://localhost:15522   # Your backend URL
FRONTEND_URL=http://localhost:3000 # Your frontend URL
JWT_SECRET=your-jwt-secret-key    # Use a strong secret key

# RabbitMQ Configuration (Optional - for email features)
RABBITMQ_URL=amqp://user:password@rabbitmq-host:5672/
RABBITMQ_HOST=your-rabbitmq-host
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=your-rabbitmq-user
RABBITMQ_PASSWORD=your-rabbitmq-password
```

### Database Setup Guide

#### PostgreSQL Setup

1. **Install PostgreSQL** (if not already installed)
   ```bash
   # Ubuntu/Debian
   sudo apt update
   sudo apt install postgresql postgresql-contrib
   
   # macOS with Homebrew
   brew install postgresql
   
   # Start PostgreSQL service
   sudo systemctl start postgresql  # Linux
   brew services start postgresql   # macOS
   ```

2. **Create database and user**
   ```bash
   sudo -u postgres psql
   
   # In PostgreSQL shell:
   CREATE DATABASE urlshortener;
   CREATE USER urluser WITH PASSWORD 'your-password';
   GRANT ALL PRIVILEGES ON DATABASE urlshortener TO urluser;
   \q
   ```

3. **Run migrations**
   ```bash
   psql -h localhost -U urluser -d urlshortener < backend/migrations/001_initial_schema.sql
   psql -h localhost -U urluser -d urlshortener < backend/migrations/003_add_otp_and_link_limits.sql
   ```

#### Redis Setup

1. **Install Redis** (if not already installed)
   ```bash
   # Ubuntu/Debian
   sudo apt update
   sudo apt install redis-server
   
   # macOS with Homebrew
   brew install redis
   
   # Start Redis service
   sudo systemctl start redis-server  # Linux
   brew services start redis          # macOS
   ```

2. **Test Redis connection**
   ```bash
   redis-cli ping
   # Should return: PONG
   ```

## 🏗️ Database Schema

The application uses PostgreSQL with the following main tables:

- **users** - User accounts with authentication
- **urls** - Shortened URLs with user ownership
- **click_events** - Detailed click tracking for analytics
- **otp_verifications** - OTP codes for email verification

Migration files are located in `backend/migrations/` and should be run in order.

## 📚 API Documentation

### Authentication Endpoints

```bash
POST /api/v1/auth/register    # User registration
POST /api/v1/auth/login       # User login
POST /api/v1/auth/logout      # User logout
POST /api/v1/auth/refresh     # Refresh JWT token
```

### URL Management Endpoints

```bash
POST   /api/v1/urls                     # Create short URL
GET    /api/v1/urls                     # Get user's URLs
GET    /api/v1/urls/:shortCode          # Get URL statistics
PUT    /api/v1/urls/:shortCode          # Update URL
DELETE /api/v1/urls/:shortCode          # Delete URL
GET    /api/v1/urls/:shortCode/analytics # Get detailed analytics
GET    /api/v1/urls/:shortCode/qr       # Generate QR code
```

### Public Endpoints

```bash
GET /:shortCode    # Redirect to original URL
GET /health        # Health check
```

## 💻 Usage Examples

### Create a Short URL

```bash
curl -X POST http://localhost:15522/api/v1/urls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "url": "https://example.com/very/long/url",
    "custom_code": "my-link",
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

### Get QR Code

```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
  http://localhost:15522/api/v1/urls/my-link/qr -o qr-code.png
```

## 🔒 Security Features

- **JWT Authentication** with secure token validation
- **Password Hashing** using bcrypt
- **User Isolation** - Complete data separation between users
- **Rate Limiting** - Prevents abuse
- **CORS Protection** with configurable origins
- **Security Headers** (XSS, CSRF protection)
- **Input Validation** and sanitization

## 📊 Analytics Features

Each shortened URL tracks:
- Total and unique clicks
- Geographic data (country/city)
- Referrer information
- Device and browser details
- Click timestamps for trend analysis

## 🐳 Docker Commands

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Rebuild and restart
docker-compose down && docker-compose up -d --build

# View running containers
docker-compose ps
```

## 🚨 Troubleshooting

### Common Issues

1. **Database connection failed**
   - Verify PostgreSQL is running
   - Check database credentials in `.env`
   - Ensure database exists and migrations are applied

2. **Redis connection failed**
   - Verify Redis server is running
   - Check Redis host and port in `.env`
   - Test connection with `redis-cli ping`

3. **Frontend can't reach backend**
   - Verify backend is running on port 15522
   - Check CORS settings in backend configuration
   - Ensure `VITE_API_URL` is set correctly

4. **Authentication issues**
   - Verify JWT_SECRET is set in `.env`
   - Check if user registration is working
   - Ensure database has proper user table structure

### Development Tips

- Use `docker-compose logs backend` to view backend logs
- Use `docker-compose logs frontend` to view frontend logs
- Backend runs on port 15522, frontend on port 3000
- Database migrations must be run manually before first startup

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details. 