# Frontend Setup Guide

## Quick Start

### Option 1: Development Mode (Recommended for development)

```bash
# Terminal 1: Start backend
make up

# Terminal 2: Start frontend dev server
cd frontend
npm install  # First time only
npm start
```

Frontend will open at http://localhost:3000 (or 3001 if 3000 is busy)

### Option 2: Production Build (Served by Go backend)

```bash
# Build React app
cd frontend
npm install  # First time only
npm run build

# Start backend (it will serve the React build)
cd ..
make up
```

Frontend will be served at http://localhost:8080

## Features

✅ **Shorten URLs** - Paste long URL, get short link  
✅ **Custom Codes** - Choose your own short code  
✅ **Copy to Clipboard** - One-click copy  
✅ **Beautiful UI** - Gradient design, responsive  
✅ **Error Handling** - Clear error messages  
✅ **API Docs** - Built-in API documentation  

## Tech Stack

- React 18
- Pure CSS3 (no frameworks)
- Fetch API for HTTP
- React Hooks (useState, useEffect)

## Development

```bash
cd frontend

# Install dependencies
npm install

# Start dev server (hot reload)
npm start

# Build for production
npm run build

# The build goes to frontend/build/
# Backend serves this automatically
```

## API Proxy

The frontend uses a proxy in `package.json`:

```json
"proxy": "http://localhost:8080"
```

This means:
- `/shorten` → proxied to `http://localhost:8080/shorten`
- `/api/urls/:code` → proxied to `http://localhost:8080/api/urls/:code`

No CORS issues in development!

## Screenshots

### Main Form
- Large input for long URL
- Optional custom code field
- "Shorten URL" button

### Result Display
- Shows short URL with copy button
- Original URL displayed
- Success message

### Info Cards
- How to Use guide
- Feature highlights (Fast, Secure, Simple)
- API endpoint documentation

## Customization

### Change Colors

Edit `frontend/src/index.css` and `frontend/src/App.css`:

```css
/* Main gradient */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

/* Primary button */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
```

### Change Base URL

If backend runs on different port, update `package.json`:

```json
"proxy": "http://localhost:YOURPORT"
```

## Troubleshooting

**Port 3000 already in use:**
```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9

# Or use different port
PORT=3001 npm start
```

**Backend not responding:**
```bash
# Check backend is running
curl http://localhost:8080/health

# Restart backend
make restart
```

**Build errors:**
```bash
# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

## Production Deployment

The React build is static files that can be:

1. **Served by Go backend** (current setup)
2. **Deployed to CDN** (S3, Cloudflare Pages, Netlify)
3. **Served by Nginx** (reverse proxy)

For AWS deployment (Phase 3), the React build will be:
- Built in CI/CD pipeline
- Copied to Docker image
- Served by Go backend in ECS

## File Structure

```
frontend/
├── public/
│   ├── index.html       # HTML template
│   └── favicon.ico      # Tab icon
├── src/
│   ├── App.js           # Main component
│   ├── App.css          # Component styles
│   ├── index.js         # React entry point
│   └── index.css        # Global styles
├── package.json         # Dependencies & scripts
└── README.md            # Frontend-specific docs
```

## Next Steps

- Add URL list/history feature
- Add statistics dashboard
- Add URL edit functionality
- Add analytics charts
- Add QR code generation
