# Docker Commands Reference

## Start Everything
```bash
docker-compose up -d
```

## Stop Everything
```bash
docker-compose down
```

## Stop Everything + Delete Data
```bash
docker-compose down -v  # Warning: deletes database and cache data!
```

## View Logs
```bash
docker-compose logs -f           # All services
docker-compose logs -f app       # Just the Go app
docker-compose logs -f postgres  # Just Postgres
docker-compose logs -f redis     # Just Redis
```

## Check Status
```bash
docker-compose ps
```

## Rebuild After Code Changes
```bash
docker-compose up -d --build
```

## Access Containers
```bash
docker exec -it url_shortener_app sh         # Go app shell
docker exec -it url_shortener_postgres sh    # Postgres shell
docker exec -it url_shortener_redis sh       # Redis shell
```

## Direct Database Access
```bash
docker exec -it url_shortener_postgres psql -U postgres -d url_shortening_service
```

## Direct Redis Access
```bash
docker exec -it url_shortener_redis redis-cli
```

## Restart Individual Service
```bash
docker-compose restart app
docker-compose restart postgres
docker-compose restart redis
```

## URLs
- **App**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Postgres**: localhost:5432
- **Redis**: localhost:6379
