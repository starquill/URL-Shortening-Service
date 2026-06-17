package generator

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	// Characters used for short code generation
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Default length of generated short codes
	defaultLength = 7
	// Maximum retry attempts for collision resolution
	maxRetries = 5
)

// ExistsChecker defines the interface for checking if a short code exists
type ExistsChecker interface {
	ShortCodeExists(ctx context.Context, shortCode string) (bool, error)
}

// ShortCodeGenerator handles generation of unique short codes
type ShortCodeGenerator struct {
	checker ExistsChecker
	length  int
}

// NewShortCodeGenerator creates a new short code generator
func NewShortCodeGenerator(checker ExistsChecker, length int) *ShortCodeGenerator {
	if length <= 0 {
		length = defaultLength
	}
	return &ShortCodeGenerator{
		checker: checker,
		length:  length,
	}
}

// Generate creates a unique short code that doesn't exist in the database
func (g *ShortCodeGenerator) Generate(ctx context.Context) (string, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		code, err := g.generateRandom()
		if err != nil {
			return "", fmt.Errorf("failed to generate random code: %w", err)
		}

		// Check if code already exists
		exists, err := g.checker.ShortCodeExists(ctx, code)
		if err != nil {
			return "", fmt.Errorf("failed to check if code exists: %w", err)
		}

		if !exists {
			return code, nil
		}

		// Collision detected, retry
	}

	return "", fmt.Errorf("failed to generate unique code after %d attempts", maxRetries)
}

// generateRandom creates a random string of specified length from charset
func (g *ShortCodeGenerator) generateRandom() (string, error) {
	result := make([]byte, g.length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < g.length; i++ {
		// Generate cryptographically secure random number
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

// ValidateCustomCode checks if a custom short code is valid
func ValidateCustomCode(code string) error {
	if len(code) < 3 {
		return fmt.Errorf("short code must be at least 3 characters long")
	}

	if len(code) > 20 {
		return fmt.Errorf("short code must be at most 20 characters long")
	}

	// Check if all characters are in allowed charset
	for _, char := range code {
		valid := false
		for _, allowed := range charset {
			if char == allowed {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("short code contains invalid character: %c", char)
		}
	}

	return nil
}
