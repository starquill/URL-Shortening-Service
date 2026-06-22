package handler

import (
	"encoding/json"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func Home(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>URL Shortening Service</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            max-width: 800px;
            margin: 80px auto;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        .container {
            background: rgba(255, 255, 255, 0.1);
            padding: 40px;
            border-radius: 16px;
            backdrop-filter: blur(10px);
        }
        h1 { margin-top: 0; font-size: 2.5em; }
        h2 { margin-top: 32px; font-size: 1.8em; }
        .endpoint {
            background: rgba(255, 255, 255, 0.2);
            padding: 16px;
            margin: 12px 0;
            border-radius: 8px;
            font-family: 'Courier New', monospace;
        }
        .method {
            display: inline-block;
            padding: 4px 8px;
            background: rgba(255, 255, 255, 0.3);
            border-radius: 4px;
            font-weight: bold;
            margin-right: 8px;
            font-size: 13px;
        }
        a { color: #fff; text-decoration: underline; }
        .footer {
            margin-top: 40px;
            opacity: 0.8;
            font-size: 14px;
            border-top: 1px solid rgba(255, 255, 255, 0.2);
            padding-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>URL Shortening Service</h1>
        <p>A production-ready API for shortening URLs with Redis caching and PostgreSQL storage.</p>

        <h2>API Endpoints</h2>

        <div class="endpoint">
            <span class="method">POST</span> /shorten
            <br>Create a short URL
        </div>

        <div class="endpoint">
            <span class="method">GET</span> /:code
            <br>Redirect to original URL
        </div>

        <div class="endpoint">
            <span class="method">GET</span> /api/urls/:code
            <br>Get URL details and analytics
        </div>

        <div class="endpoint">
            <span class="method">PUT</span> /api/urls/:code
            <br>Update a URL
        </div>

        <div class="endpoint">
            <span class="method">DELETE</span> /api/urls/:code
            <br>Delete a URL
        </div>

        <div class="endpoint">
            <span class="method">GET</span> /health
            <br>Health check
        </div>

        <h2>Documentation</h2>
        <p>Full API documentation: <a href="https://github.com/starquill/URL-Shortening-Service" target="_blank">GitHub Repository</a></p>

        <div class="footer">
            Built with Go + PostgreSQL + Redis | Deployed on Render
        </div>
    </div>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
