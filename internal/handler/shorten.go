package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/starquill/URL-Shortening-Service/internal/cache"
	"github.com/starquill/URL-Shortening-Service/internal/generator"
	"github.com/starquill/URL-Shortening-Service/internal/store"
)

type URLHandler struct {
	urlStore      *store.URLStore
	cache         *cache.RedisCache
	generator     *generator.ShortCodeGenerator
	baseURL       string
}

func NewURLHandler(urlStore *store.URLStore, cache *cache.RedisCache, gen *generator.ShortCodeGenerator, baseURL string) *URLHandler {
	return &URLHandler{
		urlStore:  urlStore,
		cache:     cache,
		generator: gen,
		baseURL:   baseURL,
	}
}

// Request/Response types
type ShortenRequest struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code,omitempty"`
}

type ShortenResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLResponse struct {
	ID          int64  `json:"id"`
	URL         string `json:"url"`
	ShortCode   string `json:"short_code"`
	AccessCount int64  `json:"access_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateURLRequest struct {
	URL string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// POST /shorten - Create short URL
func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate URL
	if req.URL == "" {
		respondError(w, http.StatusBadRequest, "url is required")
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		respondError(w, http.StatusBadRequest, "invalid url format")
		return
	}

	ctx := r.Context()
	var shortCode string
	var err error

	// Use custom code or generate one
	if req.CustomCode != "" {
		// Validate custom code
		if err := generator.ValidateCustomCode(req.CustomCode); err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check if custom code already exists
		exists, err := h.urlStore.ShortCodeExists(ctx, req.CustomCode)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to check code availability")
			return
		}

		if exists {
			respondError(w, http.StatusConflict, "custom code already exists")
			return
		}

		shortCode = req.CustomCode
	} else {
		// Generate random code
		shortCode, err = h.generator.Generate(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to generate short code")
			return
		}
	}

	// Save to database
	urlRecord, err := h.urlStore.Create(ctx, req.URL, shortCode)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create short url")
		return
	}

	// Cache the short code -> URL mapping
	if err := h.cache.Set(ctx, shortCode, req.URL); err != nil {
		// Log error but don't fail the request
		// Cache miss will just hit the database
	}

	respondJSON(w, http.StatusCreated, ShortenResponse{
		ShortCode:   urlRecord.ShortCode,
		ShortURL:    h.baseURL + "/" + urlRecord.ShortCode,
		OriginalURL: urlRecord.URL,
	})
}

// GET /:code - Redirect to original URL
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		respondError(w, http.StatusBadRequest, "short code is required")
		return
	}

	ctx := r.Context()

	// Try cache first
	originalURL, err := h.cache.Get(ctx, shortCode)
	if err == nil {
		// Cache hit - increment count asynchronously and redirect
		go h.urlStore.IncrementAccessCount(context.Background(), shortCode)
		http.Redirect(w, r, originalURL, http.StatusFound)
		return
	}

	// Cache miss - get from database
	urlRecord, err := h.urlStore.GetByShortCode(ctx, shortCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "short code not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to retrieve url")
		return
	}

	// Update cache
	h.cache.Set(ctx, shortCode, urlRecord.URL)

	// Increment access count asynchronously
	go h.urlStore.IncrementAccessCount(context.Background(), shortCode)

	// Redirect
	http.Redirect(w, r, urlRecord.URL, http.StatusFound)
}

// GET /api/urls/:code - Get URL details
func (h *URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		respondError(w, http.StatusBadRequest, "short code is required")
		return
	}

	ctx := r.Context()
	urlRecord, err := h.urlStore.GetByShortCode(ctx, shortCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "short code not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to retrieve url")
		return
	}

	respondJSON(w, http.StatusOK, URLResponse{
		ID:          urlRecord.ID,
		URL:         urlRecord.URL,
		ShortCode:   urlRecord.ShortCode,
		AccessCount: urlRecord.AccessCount,
		CreatedAt:   urlRecord.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   urlRecord.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// PUT /api/urls/:code - Update URL
func (h *URLHandler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		respondError(w, http.StatusBadRequest, "short code is required")
		return
	}

	var req UpdateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.URL == "" {
		respondError(w, http.StatusBadRequest, "url is required")
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		respondError(w, http.StatusBadRequest, "invalid url format")
		return
	}

	ctx := r.Context()

	// Update in database
	urlRecord, err := h.urlStore.Update(ctx, shortCode, req.URL)
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "short code not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update url")
		return
	}

	// Invalidate cache
	h.cache.Delete(ctx, shortCode)

	respondJSON(w, http.StatusOK, URLResponse{
		ID:          urlRecord.ID,
		URL:         urlRecord.URL,
		ShortCode:   urlRecord.ShortCode,
		AccessCount: urlRecord.AccessCount,
		CreatedAt:   urlRecord.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   urlRecord.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// DELETE /api/urls/:code - Delete URL
func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		respondError(w, http.StatusBadRequest, "short code is required")
		return
	}

	ctx := r.Context()

	// Delete from database
	err := h.urlStore.Delete(ctx, shortCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			respondError(w, http.StatusNotFound, "short code not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete url")
		return
	}

	// Invalidate cache
	h.cache.Delete(ctx, shortCode)

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
