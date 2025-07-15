# URL Shortener API v2.0

A **secure, multi-user URL shortener service** built with Go, featuring JWT authentication, advanced analytics, and comprehensive URL management.

## 🚀 Features

- **🔐 User Authentication** - JWT-based secure access
- **⚡ Fast URL shortening** with custom codes
- **⏰ URL expiration** with automatic handling
- **📊 Advanced analytics** with geographic tracking
- **📱 QR code generation** (native Go implementation)
- **🔒 User isolation** - Users only see their own URLs
- **📈 Real-time metrics** and click tracking
- **🚀 High performance** with Redis caching

## 🛠️ Tech Stack

- **Backend**: Go 1.21+ with Gin framework
- **Authentication**: JWT tokens with bcrypt password hashing
- **Database**: PostgreSQL with user isolation
- **Cache**: Redis for performance
- **QR Codes**: Native Go generation (github.com/skip2/go-qrcode)

## 🚀 Quick Start

### Using Docker Compose
```bash
git clone <repository>
cd url-shortener/backend
docker-compose up -d
```

### Manual Setup
```bash
# 1. Install dependencies
go mod download

# 2. Set up environment
cp env.example .env
# Edit .env with your settings

# 3. Set up PostgreSQL
createdb urlshortener
psql urlshortener < migrations/001_initial_schema.sql
psql urlshortener < migrations/002_add_users_table.sql

# 4. Start Redis
redis-server

# 5. Run the application
go run cmd/main.go
```

## 🔐 Authentication

All URL operations require authentication. Users can only access their own URLs.

### Register
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

# Response includes JWT token
{
  "user": {...},
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Using the Token
Include in all protected requests:
```bash
Authorization: Bearer <your-jwt-token>
```

## 📚 API Endpoints

### 🔓 Public Endpoints
```
POST /api/v1/auth/register     # User registration
POST /api/v1/auth/login        # User login
POST /api/v1/auth/logout       # User logout
GET  /:shortCode               # URL redirect (public)
GET  /health                   # Health check
```

### 🔒 Protected Endpoints (Require Authentication)

#### User Management
```
GET  /api/v1/profile                    # Get user profile
PUT  /api/v1/profile                    # Update profile
POST /api/v1/profile/change-password    # Change password
POST /api/v1/auth/refresh               # Refresh JWT token
```

#### URL Management
```
POST   /api/v1/urls                     # Create URL
GET    /api/v1/urls                     # Get user's URLs
GET    /api/v1/urls/:shortCode          # Get URL stats
PUT    /api/v1/urls/:shortCode          # Update URL
DELETE /api/v1/urls/:shortCode          # Delete URL
GET    /api/v1/urls/:shortCode/analytics # Get analytics
GET    /api/v1/urls/:shortCode/qr       # Generate QR code
```

## 💻 Usage Examples

### Create URL
```bash
curl -X POST http://localhost:15522/api/v1/urls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "url": "https://example.com/very/long/url",
    "custom_code": "my-link",
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

### Get User's URLs
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:15522/api/v1/urls?limit=10&offset=0"
```

### Get QR Code
```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:15522/api/v1/urls/my-link/qr -o qr-code.png
```

## ⚙️ Configuration

Configure via environment variables:

```bash
# Server
SERVER_PORT=15522

# Database
DB_HOST=localhost
DB_PORT=15521
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=urlshortener

# Redis
REDIS_HOST=localhost
REDIS_PORT=7655

# Application
BASE_URL=http://localhost:15522
JWT_SECRET=your-secret-key
```

## 🏗️ Database Schema

The service uses PostgreSQL with two main tables:
- `users` - User accounts with authentication
- `urls` - URLs with user ownership and analytics
- `click_events` - Detailed click tracking

## 🔒 Security Features

- **JWT Authentication** with secure token validation
- **Password Hashing** using bcrypt
- **User Isolation** - Complete data separation
- **Rate Limiting** - 100 requests/second
- **CORS Protection** with configurable origins
- **Security Headers** (XSS, CSRF protection)

## 📊 Analytics

Each URL tracks:
- Total and unique clicks
- Geographic data (country/city)
- Referrer information
- Click timestamps
- User agent details

## 🐳 Docker Support

```bash
# Build and run
docker-compose up -d

# The API will be available at http://localhost:15522
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License. 