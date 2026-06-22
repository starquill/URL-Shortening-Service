# --- frontend build stage ---
FROM node:18-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# --- backend build stage ---
FROM golang:1.25-alpine AS backend-build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# --- run stage ---
FROM alpine:3.20
WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

COPY --from=backend-build /app/server .
COPY --from=frontend-build /app/frontend/build ./frontend/build
COPY config.yaml .
COPY config.prod.yaml .
COPY migrations ./migrations

EXPOSE 8080
CMD ["./server"]
