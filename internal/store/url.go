package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URL struct {
	ID          int64     `json:"id"`
	URL         string    `json:"url"`
	ShortCode   string    `json:"short_code"`
	AccessCount int64     `json:"access_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type URLStore struct {
	pool *pgxpool.Pool
}

func NewURLStore(pool *pgxpool.Pool) *URLStore {
	return &URLStore{pool: pool}
}

// Create inserts a new URL with the given short code
func (s *URLStore) Create(ctx context.Context, url, shortCode string) (*URL, error) {
	query := `
		INSERT INTO urls (url, short_code)
		VALUES ($1, $2)
		RETURNING id, url, short_code, access_count, created_at, updated_at
	`

	var u URL
	err := s.pool.QueryRow(ctx, query, url, shortCode).Scan(
		&u.ID,
		&u.URL,
		&u.ShortCode,
		&u.AccessCount,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetByShortCode retrieves a URL by its short code
func (s *URLStore) GetByShortCode(ctx context.Context, shortCode string) (*URL, error) {
	query := `
		SELECT id, url, short_code, access_count, created_at, updated_at
		FROM urls
		WHERE short_code = $1
	`

	var u URL
	err := s.pool.QueryRow(ctx, query, shortCode).Scan(
		&u.ID,
		&u.URL,
		&u.ShortCode,
		&u.AccessCount,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// Update updates the URL for a given short code
func (s *URLStore) Update(ctx context.Context, shortCode, newURL string) (*URL, error) {
	query := `
		UPDATE urls
		SET url = $1, updated_at = NOW()
		WHERE short_code = $2
		RETURNING id, url, short_code, access_count, created_at, updated_at
	`

	var u URL
	err := s.pool.QueryRow(ctx, query, newURL, shortCode).Scan(
		&u.ID,
		&u.URL,
		&u.ShortCode,
		&u.AccessCount,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// Delete removes a URL by its short code
func (s *URLStore) Delete(ctx context.Context, shortCode string) error {
	query := `DELETE FROM urls WHERE short_code = $1`
	_, err := s.pool.Exec(ctx, query, shortCode)
	return err
}

// IncrementAccessCount increments the access count for a short code
func (s *URLStore) IncrementAccessCount(ctx context.Context, shortCode string) error {
	query := `
		UPDATE urls
		SET access_count = access_count + 1
		WHERE short_code = $1
	`
	_, err := s.pool.Exec(ctx, query, shortCode)
	return err
}

// ShortCodeExists checks if a short code already exists
func (s *URLStore) ShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`

	var exists bool
	err := s.pool.QueryRow(ctx, query, shortCode).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
