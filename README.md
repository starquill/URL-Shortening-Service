# URL Shortening Service

A production-ready URL shortener built with Go, PostgreSQL, Redis, and React. Features cryptographically secure short code generation, caching, analytics, and a modern web interface.

**🌐 Live Demo:** [https://url-shortening-service-dbam.onrender.com](https://url-shortening-service-dbam.onrender.com)

> **Note:** First request may take ~30 seconds to wake up (free tier limitation)

## 🚀 Features

- **Short URL Generation**: Crypto-secure 7-character codes with collision detection
- **Custom Codes**: User-defined short codes (3-20 characters)
- **Caching Layer**: Redis-powered caching with automatic invalidation
- **Access Tracking**: Real-time click analytics for shortened URLs
- **REST API**: 5 production-ready endpoints
- **Modern UI**: React-based frontend with copy-to-clipboard
- **Production Ready**: Docker containerized, fully tested, deployed

## 🏗️ Architecture

```
┌─────────────┐
│   React UI  │
└──────┬──────┘
       │ REST API
┌──────▼──────────────────┐
│   Go Backend (Chi)      │
│  ┌────────────────────┐ │
│  │   API Handlers     │ │
│  └────────┬───────────┘ │
│  ┌────────▼───────────┐ │
│  │  Cache (Redis)     │ │  ← 24h TTL, cache invalidation
│  └────────┬───────────┘ │
│  ┌────────▼───────────┐ │
│  │  Store (Postgres)  │ │  ← Persistent storage
│  └────────────────────┘ │
└─────────────────────────┘
```

## 🛠️ Tech Stack

**Backend:**
- Go 1.25 - High-performance compiled language
- Chi Router - Lightweight HTTP router
- PostgreSQL 16 - Relational database with migrations
- Redis 7 - In-memory caching
- golang-migrate - Database migration management
- Viper - Configuration management

**Frontend:**
- React 18 - Modern UI library
- Fetch API - HTTP client
- CSS3 - Gradient styling

**DevOps:**
- Docker - Multi-stage containerization
- Docker Compose - Local orchestration
- Makefile - 16+ development commands
- GitHub Actions Ready - CI/CD pipeline

**Testing:**
- 18 tests (6 unit + 12 E2E)
- httptest - API testing
- Real database integration tests

## 📦 Installation

### Prerequisites
- Docker & Docker Compose
- Go 1.25+ (for local development)
- Make (optional, for convenience commands)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/starquill/URL-Shortening-Service.git
cd URL-Shortening-Service

# Start all services (PostgreSQL + Redis + API)
make up

# In another terminal, start the frontend
cd frontend
npm install
npm start

# Your app is now running!
# Backend:  http://localhost:8080
# Frontend: http://localhost:3000
```

## 🎯 API Endpoints

### Create Short URL
```bash
POST /shorten
Content-Type: application/json

{
  "url": "https://example.com",
  "custom_code": "mycode" // optional
}

Response: {
  "short_code": "abc123",
  "short_url": "http://localhost:8080/abc123",
  "original_url": "https://example.com"
}
```

### Redirect to Original URL
```bash
GET /:code
Response: 302 Redirect to original URL
```

### Get URL Details
```bash
GET /api/urls/:code
Response: {
  "id": 1,
  "url": "https://example.com",
  "short_code": "abc123",
  "access_count": 42,
  "created_at": "2026-06-22T10:00:00Z",
  "updated_at": "2026-06-22T10:00:00Z"
}
```

### Update URL
```bash
PUT /api/urls/:code
Content-Type: application/json

{
  "url": "https://updated-url.com"
}
```

### Delete URL
```bash
DELETE /api/urls/:code
Response: 204 No Content
```

### Health Check
```bash
GET /health
Response: {
  "status": "ok",
  "timestamp": "2026-06-22T10:00:00Z"
}
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run E2E tests only (requires running services)
make test-e2e

# Run with coverage
go test ./... -cover
```

## 🐳 Docker Commands

```bash
# Start services
make up

# Stop services (keeps data)
make down

# Rebuild after code changes
make rebuild

# View logs
make logs          # all services
make logs-app      # app only
make logs-postgres # database only
make logs-redis    # cache only

# Database access
make db-shell      # PostgreSQL shell
make redis-cli     # Redis CLI

# Clean everything (deletes data!)
make clean
```

## 📁 Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── cache/           # Redis cache implementation
│   ├── config/          # Configuration management
│   ├── database/        # PostgreSQL connection & migrations
│   ├── generator/       # Short code generation
│   ├── handler/         # HTTP request handlers
│   └── store/           # Database operations (CRUD)
├── migrations/          # SQL migration files
├── test/
│   └── e2e/            # End-to-end API tests
├── frontend/           # React UI
├── docker-compose.yml  # Local development orchestration
├── Dockerfile          # Multi-stage production build
└── Makefile           # Development commands
```

## 🔧 Configuration

Configuration is managed via `config.yaml` and environment variables:

```yaml
server:
  port: 8080
  base_url: http://localhost:8080

database:
  url: postgres://user:pass@host:5432/dbname

redis:
  url: localhost:6379
  password: ""  # optional
  ttl: 24h

cors:
  allowed_origins:
    - http://localhost:3000
```

Environment variables override config file (e.g., `DATABASE_URL`, `REDIS_URL`).

## 🚢 Deployment

This application is deployed on [Render.com](https://render.com) with:
- Web Service (Go backend + React frontend)
- PostgreSQL database (1GB free tier)
- Redis Cloud (30MB free tier)
- Automatic deployments from GitHub
- Zero-downtime rolling updates

**🌐 Live Demo:** [https://url-shortening-service-dbam.onrender.com](https://url-shortening-service-dbam.onrender.com)

**Architecture in Production:**
```
GitHub Push → Render Webhook → Docker Build → Deploy → Live
```

## 📊 Key Features & Implementation

### Cryptographically Secure Short Codes
- Uses `crypto/rand` for unpredictable codes
- 62-character set (a-z, A-Z, 0-9)
- 7 characters = 3.5 trillion possible combinations
- Collision detection with retry (up to 5 attempts)

### Caching Strategy
1. **Cache-first reads**: Check Redis before PostgreSQL
2. **Write-through**: Cache on database miss
3. **Cache invalidation**: On update/delete operations
4. **TTL**: 24-hour expiration

### Access Count Tracking
- Asynchronous increment (doesn't block redirect)
- Uses goroutines for non-blocking updates
- Tracked on every redirect (not on API reads)

## 🎯 Development

### Run Locally Without Docker

```bash
# Start PostgreSQL and Redis manually
# Update config.yaml with connection strings

# Run migrations
migrate -path migrations -database "postgres://..." up

# Start the server
go run cmd/server/main.go

# In another terminal, start frontend
cd frontend && npm start
```

### Add a New Migration

```bash
# Create migration files
migrate create -ext sql -dir migrations -seq <migration_name>

# Edit the .up.sql and .down.sql files
# Restart the app - migrations run automatically
```

## 🤝 Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new features
4. Ensure all tests pass
5. Submit a pull request

## 📝 License

MIT License - see LICENSE file for details

## 👤 Author

**Ajay Gitala**
- GitHub: [@starquill](https://github.com/starquill)

## 🙏 Acknowledgments

- Built with Go's excellent standard library
- Chi router for minimal HTTP routing
- PostgreSQL for reliable data storage
- Redis for blazing-fast caching
