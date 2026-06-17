# Makefile Quick Reference

Instead of typing long docker-compose commands, just use `make`:

## 🚀 Most Common Commands

```bash
make up        # Start everything
make down      # Stop everything
make rebuild   # Rebuild after code changes
make logs-app  # Watch app logs
make health    # Test if app is working
```

## 📋 All Commands

### Starting/Stopping
| Command | What it does |
|---------|--------------|
| `make up` | Start all services (app, postgres, redis) |
| `make down` | Stop all services (keeps your data) |
| `make restart` | Restart all services |
| `make clean` | ⚠️ Stop and DELETE all data |

### Development
| Command | What it does |
|---------|--------------|
| `make rebuild` | Rebuild app after changing Go code |
| `make health` | Test the /health endpoint |
| `make test` | Run health check + show status |

### Logs
| Command | What it does |
|---------|--------------|
| `make logs` | Show all logs (app + postgres + redis) |
| `make logs-app` | Show only app logs |
| `make logs-postgres` | Show only postgres logs |
| `make logs-redis` | Show only redis logs |

**Tip:** Press `Ctrl+C` to stop watching logs (services keep running)

### Status
| Command | What it does |
|---------|--------------|
| `make ps` | Show which containers are running |

### Access Containers
| Command | What it does |
|---------|--------------|
| `make shell-app` | Open shell in app container |
| `make shell-postgres` | Open shell in postgres container |
| `make shell-redis` | Open shell in redis container |
| `make db-shell` | Open PostgreSQL interactive terminal |
| `make redis-cli` | Open Redis CLI |

---

## 📖 Usage Examples

### Typical Workflow

```bash
# Morning: Start everything
make up

# After changing code in main.go
make rebuild

# Check if it's working
make health

# Watch logs for errors
make logs-app

# Evening: Stop everything
make down
```

### Debugging

```bash
# Check status
make ps

# Watch logs
make logs-app

# Get inside app container
make shell-app

# Access database directly
make db-shell
```

### Fresh Start

```bash
# Delete everything and start over
make clean
make up
```

---

## 💡 Tips

1. **Just type `make`** to see all commands
2. **Use `make help`** for detailed help
3. **`make down` is safe** - keeps your data
4. **`make clean` deletes data** - only use when you want a fresh start
5. **After editing Go code** - always run `make rebuild`

---

## 🆚 Make vs Docker Compose

| Instead of this long command... | Just type this! |
|----------------------------------|-----------------|
| `docker-compose up -d` | `make up` |
| `docker-compose down` | `make down` |
| `docker-compose up -d --build` | `make rebuild` |
| `docker-compose logs -f app` | `make logs-app` |
| `docker exec -it url_shortener_postgres psql -U postgres -d url_shortening_service` | `make db-shell` |
| `curl http://localhost:8080/health` | `make health` |

Much easier! 🎉
