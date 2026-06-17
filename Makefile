.PHONY: help up down restart rebuild logs logs-app logs-postgres logs-redis ps clean shell-app shell-postgres shell-redis db-shell redis-cli test health

# Default target: show help
help:
	@echo "URL Shortening Service - Make Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  up              Start all services in background"
	@echo "  down            Stop all services (keeps data)"
	@echo "  restart         Restart all services"
	@echo "  rebuild         Rebuild and restart app (after code changes)"
	@echo "  logs            Show logs for all services (follow mode)"
	@echo "  logs-app        Show app logs (follow mode)"
	@echo "  logs-postgres   Show postgres logs (follow mode)"
	@echo "  logs-redis      Show redis logs (follow mode)"
	@echo "  ps              Show running containers status"
	@echo "  clean           Stop and remove everything including volumes (⚠️  deletes data!)"
	@echo "  shell-app       Open shell in app container"
	@echo "  shell-postgres  Open shell in postgres container"
	@echo "  shell-redis     Open shell in redis container"
	@echo "  db-shell        Open PostgreSQL interactive terminal"
	@echo "  redis-cli       Open Redis CLI"
	@echo "  test            Run health check test"
	@echo "  health          Check health endpoint"
	@echo ""
	@echo "Examples:"
	@echo "  make up         # Start everything"
	@echo "  make logs-app   # Watch app logs"
	@echo "  make rebuild    # After changing Go code"
	@echo "  make db-shell   # Access database"

# Start all services
up:
	@echo "🚀 Starting all services..."
	docker-compose up -d
	@echo "✅ Services started!"
	@make ps

# Stop all services (keeps data)
down:
	@echo "🛑 Stopping all services..."
	docker-compose down
	@echo "✅ Services stopped (data preserved)"

# Restart all services
restart:
	@echo "🔄 Restarting all services..."
	docker-compose restart
	@echo "✅ Services restarted"
	@make ps

# Rebuild and restart (after code changes)
rebuild:
	@echo "🔨 Rebuilding app..."
	docker-compose up -d --build
	@echo "✅ App rebuilt and restarted"
	@make ps

# Show all logs (follow mode)
logs:
	docker-compose logs -f

# Show app logs
logs-app:
	docker-compose logs -f app

# Show postgres logs
logs-postgres:
	docker-compose logs -f postgres

# Show redis logs
logs-redis:
	docker-compose logs -f redis

# Show container status
ps:
	@docker-compose ps

# Clean everything (deletes data!)
clean:
	@echo "⚠️  WARNING: This will delete all data!"
	@echo "Press Ctrl+C to cancel, or wait 3 seconds to continue..."
	@sleep 3
	@echo "🧹 Cleaning everything..."
	docker-compose down -v
	@echo "✅ Everything cleaned (data deleted)"

# Open shell in app container
shell-app:
	@echo "🐚 Opening shell in app container..."
	docker exec -it url_shortener_app sh

# Open shell in postgres container
shell-postgres:
	@echo "🐚 Opening shell in postgres container..."
	docker exec -it url_shortener_postgres sh

# Open shell in redis container
shell-redis:
	@echo "🐚 Opening shell in redis container..."
	docker exec -it url_shortener_redis sh

# Open PostgreSQL interactive terminal
db-shell:
	@echo "🗃️  Opening PostgreSQL shell..."
	docker exec -it url_shortener_postgres psql -U postgres -d url_shortening_service

# Open Redis CLI
redis-cli:
	@echo "🔴 Opening Redis CLI..."
	docker exec -it url_shortener_redis redis-cli

# Run health check
health:
	@echo "🏥 Checking health endpoint..."
	@curl -s http://localhost:8080/health | jq '.' || echo "❌ Health check failed"

# Run tests (placeholder for future)
test:
	@echo "🧪 Running tests..."
	@echo "Health Check:"
	@curl -s http://localhost:8080/health | jq '.'
	@echo ""
	@echo "Container Status:"
	@make ps
