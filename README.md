# URL Shortening Service

A RESTful URL shortener built in Go. Create, manage, and track short links for long URLs.

## Requirements

- Go 1.21 or later

The server starts on `http://localhost:8080` by default.

## API

| Method | Path | Description |
| ------ | ---- | ----------- |
| POST | `/shorten` | Create a short URL |
| GET | `/shorten/:shortCode` | Get original URL |
| PUT | `/shorten/:shortCode` | Update a short URL |
| DELETE | `/shorten/:shortCode` | Delete a short URL |
| GET | `/shorten/:shortCode/stats` | Get access stats |
