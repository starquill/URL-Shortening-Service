import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [url, setUrl] = useState('');
  const [customCode, setCustomCode] = useState('');
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');
  const [urls, setUrls] = useState([]);
  const [loading, setLoading] = useState(false);

  const API_BASE = '';  // Empty because we're using proxy in package.json

  // Clear error after 5 seconds
  useEffect(() => {
    if (error) {
      const timer = setTimeout(() => setError(''), 5000);
      return () => clearTimeout(timer);
    }
  }, [error]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setResult(null);
    setLoading(true);

    try {
      const payload = { url };
      if (customCode.trim()) {
        payload.custom_code = customCode.trim();
      }

      const response = await fetch(`${API_BASE}/shorten`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      const data = await response.json();

      if (response.ok) {
        setResult(data);
        setUrl('');
        setCustomCode('');
        loadRecentUrls(); // Refresh list
      } else {
        setError(data.error || 'Failed to create short URL');
      }
    } catch (err) {
      setError('Failed to connect to server');
    } finally {
      setLoading(false);
    }
  };

  const loadRecentUrls = async () => {
    // Note: This is a simplified version
    // In production, you'd have an endpoint to list URLs
    // For now, we'll just keep track of created ones in state
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    alert('Copied to clipboard!');
  };

  const deleteUrl = async (shortCode) => {
    if (!window.confirm('Are you sure you want to delete this URL?')) {
      return;
    }

    try {
      const response = await fetch(`${API_BASE}/api/urls/${shortCode}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        setUrls(urls.filter(u => u.short_code !== shortCode));
        alert('URL deleted successfully');
      } else {
        alert('Failed to delete URL');
      }
    } catch (err) {
      alert('Failed to connect to server');
    }
  };

  const getUrlStats = async (shortCode) => {
    try {
      const response = await fetch(`${API_BASE}/api/urls/${shortCode}`);
      const data = await response.json();

      if (response.ok) {
        alert(`Stats for ${shortCode}:\n\nOriginal URL: ${data.url}\nAccess Count: ${data.access_count}\nCreated: ${new Date(data.created_at).toLocaleString()}`);
      } else {
        alert(data.error || 'Failed to fetch stats');
      }
    } catch (err) {
      alert('Failed to connect to server');
    }
  };

  return (
    <div className="App">
      <header className="header">
        <h1>🔗 URL Shortener</h1>
        <p>Transform long URLs into short, shareable links</p>
      </header>

      <div className="card">
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="url">Enter Long URL</label>
            <input
              type="url"
              id="url"
              placeholder="https://example.com/very/long/url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="customCode">Custom Short Code (Optional)</label>
            <input
              type="text"
              id="customCode"
              placeholder="my-custom-code (3-20 characters, alphanumeric only)"
              value={customCode}
              onChange={(e) => setCustomCode(e.target.value)}
              maxLength="20"
            />
          </div>

          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Creating...' : 'Shorten URL'}
          </button>
        </form>

        {error && (
          <div className="error">
            <strong>Error:</strong> {error}
          </div>
        )}

        {result && (
          <div className="result">
            <h3>✅ Short URL Created!</h3>
            <div className="short-url">
              <a href={result.short_url} target="_blank" rel="noopener noreferrer">
                {result.short_url}
              </a>
            </div>
            <button
              className="copy-btn"
              onClick={() => copyToClipboard(result.short_url)}
            >
              📋 Copy Link
            </button>
            <p style={{ marginTop: '10px', fontSize: '14px', color: '#666' }}>
              Original: {result.original_url}
            </p>
          </div>
        )}
      </div>

      <div className="card">
        <h2>How to Use</h2>
        <ol style={{ lineHeight: '1.8', color: '#555' }}>
          <li>Paste your long URL in the input field above</li>
          <li>(Optional) Enter a custom short code, or let the system generate one</li>
          <li>Click "Shorten URL" to create your short link</li>
          <li>Share your short URL anywhere!</li>
        </ol>

        <div className="stats" style={{ marginTop: '30px' }}>
          <div className="stat-card">
            <h3>Fast</h3>
            <p>Instant URL shortening with Redis caching</p>
          </div>
          <div className="stat-card">
            <h3>Secure</h3>
            <p>Safe redirects with access tracking</p>
          </div>
          <div className="stat-card">
            <h3>Simple</h3>
            <p>Clean interface, powerful API</p>
          </div>
        </div>
      </div>

      <div className="card">
        <h2>API Endpoints</h2>
        <div style={{ fontFamily: 'monospace', fontSize: '14px', color: '#666' }}>
          <p><strong>POST /shorten</strong> - Create short URL</p>
          <p><strong>GET /:code</strong> - Redirect to original URL</p>
          <p><strong>GET /api/urls/:code</strong> - Get URL details</p>
          <p><strong>PUT /api/urls/:code</strong> - Update URL</p>
          <p><strong>DELETE /api/urls/:code</strong> - Delete URL</p>
        </div>
      </div>
    </div>
  );
}

export default App;
