# --- build stage ---
    FROM golang:1.25-alpine AS build
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

    COPY --from=build /app/server .
    COPY config.yaml .
    COPY config.prod.yaml .
    COPY migrations ./migrations

    EXPOSE 8080
    CMD ["./server"]