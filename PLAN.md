# Project Plan

## Workflow (every step)

1. Build the step (with AI if you want).
2. **Break it** — change or remove something, run it, watch it fail.
3. **Understand it** — trace the flow and explain what each part does in your own words.
4. **Commit** — one commit per completed step, only after steps 2–3.

Example commit messages:

- `phase1: step 1 — project setup and health check`
- `phase1: step 2 — docker compose for app, postgres, redis`

## Phase 1 — Backend

### Step 1: Project Setup
- `go mod init`
- Install dependencies: `chi`, `viper`, `pgx`, `golang-migrate`, `go-redis`, `rs/cors`
- Set up project folder structure (`cmd/`, `internal/`, `migrations/`, `config/`)
- Create `config/config.yaml` with server, database, Redis, and CORS settings
- Wire up `viper` to load config from YAML + env var overrides
- Add a `GET /health` endpoint to verify the server runs

### Step 2: Docker Compose (local dev)
- Write `Dockerfile` with a multi-stage build for the Go binary
- Write `docker-compose.yml` to run the app, Postgres, and Redis together
- Verify all three containers start and can talk to each other

### Step 3: Database — PostgreSQL
- Write migration files (`001_create_urls.up.sql`, `001_create_urls.down.sql`)
- Schema: `id`, `url`, `short_code` (unique index), `access_count`, `created_at`, `updated_at`
- Integrate `golang-migrate` to run migrations on startup
- Implement the store layer (CRUD + increment access count) using `pgx`

### Step 4: Cache — Redis
- Implement a cache layer using `go-redis`
- Cache `shortCode → url` on read with a configurable TTL
- Invalidate cache on update and delete
- Strategy: check Redis first, fall back to Postgres on miss, populate cache on miss

### Step 5: Short Code Generator
- Generate random 7-character alphanumeric codes using `crypto/rand`
- Ensure uniqueness by checking the DB before saving
- Retry on collision

### Step 6: API Handlers
- Implement all 5 endpoints:
  - `POST /shorten` — validate URL, generate short code, save to DB, return `201`
  - `GET /shorten/:shortCode` — check cache → DB, return `200` or `404`
  - `PUT /shorten/:shortCode` — update URL, invalidate cache, return `200` or `404`
  - `DELETE /shorten/:shortCode` — delete from DB, invalidate cache, return `204` or `404`
  - `GET /shorten/:shortCode/stats` — return URL + `accessCount`, return `200` or `404`
- Configure CORS using `rs/cors` (allowed origins, methods, headers)
- Add request validation and consistent JSON error responses

### Step 7: Tests
- Unit test the short code generator and URL validator
- Handler tests using `net/http/httptest`
- Integration tests against a real Postgres + Redis instance (via Docker)

---

## Phase 2 — Frontend (UI)

### Step 1: Static File Setup
- Serve static HTML/CSS/JS from Go using `net/http` file server
- Create pages: shorten form, result page, stats view

### Step 2: Shorten Form
- Input field for long URL
- Call `POST /shorten` on submit
- Display the returned short link

### Step 3: Redirect
- On visiting `/{shortCode}`, fetch original URL via `GET /shorten/:shortCode`
- Redirect using `window.location`

### Step 4: Stats Page
- Input a short code
- Call `GET /shorten/:shortCode/stats`
- Display original URL and access count

### Step 5: CORS Verification
- Ensure the frontend works when served from a different origin than the API
- Update `config.yaml` allowed origins accordingly

---

## Phase 3 — Deploy (FREE - Fly.io)

### Step 1: Containerize ✅
- Finalize multi-stage `Dockerfile` (build stage → minimal runtime image)
- Test the production Docker image locally
- Create `fly.toml` configuration

### Step 2: Fly.io Setup
- Install Fly CLI (`brew install flyctl`)
- Sign up for free account (no credit card required)
- Authenticate with `flyctl auth login`

### Step 3: Launch App
- Run `flyctl launch` to create app
- Provision PostgreSQL database (free tier)
- Provision Redis instance (free tier)
- Configure environment variables

### Step 4: Deploy
- Run `flyctl deploy` to build and deploy
- Verify migrations run automatically
- Test all endpoints on live URL

### Step 5: CI/CD with GitHub Actions (Optional)
- Write `.github/workflows/deploy.yml` — on merge to `main`:
  - Build Docker image
  - Deploy to Fly.io using `flyctl deploy`

### Step 6: Final Checks
- Verify all 5 API endpoints work on the live URL
- Test redirect functionality
- Check cache behavior in production
- Monitor logs via `flyctl logs`

---

## Phase 3 Alternative — Deploy (AWS) - SKIPPED (Costs Money)

If you want to deploy to AWS instead (not free):
- See `AWS_DEPLOYMENT.md` for complete guide
- Estimated cost: $70-75/month (can be reduced with stop/start strategy)
- Steps: ECR → RDS → ElastiCache → ECS Fargate → ALB → CI/CD
