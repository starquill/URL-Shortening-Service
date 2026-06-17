# URL Shortener - React Frontend

Simple, clean React UI for the URL Shortener service.

## Features

- ✅ Shorten URLs with auto-generated or custom codes
- ✅ Copy short URLs to clipboard
- ✅ Beautiful gradient UI
- ✅ Responsive design
- ✅ Real-time error handling
- ✅ API documentation display

## Quick Start

```bash
# Install dependencies
npm install

# Start development server (port 3000)
npm start

# Build for production
npm build
```

## Requirements

- Backend API running on http://localhost:8080
- Node.js 14+ and npm

## Usage

1. Start the backend: `make up` (from project root)
2. Start the frontend: `npm start` (from frontend directory)
3. Open http://localhost:3000

## API Proxy

The frontend proxies API requests to http://localhost:8080 (configured in package.json).

## Tech Stack

- React 18
- CSS3 (no frameworks, pure CSS)
- Fetch API for HTTP requests
