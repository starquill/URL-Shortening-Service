package generator

import (
	"context"
	"strings"
	"testing"
)

// MockChecker is a mock implementation of ExistsChecker for testing
type MockChecker struct {
	existingCodes map[string]bool
	checkCount    int
}

func NewMockChecker(existingCodes []string) *MockChecker {
	codes := make(map[string]bool)
	for _, code := range existingCodes {
		codes[code] = true
	}
	return &MockChecker{
		existingCodes: codes,
		checkCount:    0,
	}
}

func (m *MockChecker) ShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	m.checkCount++
	return m.existingCodes[shortCode], nil
}

func TestGenerateRandom(t *testing.T) {
	gen := NewShortCodeGenerator(nil, 7)

	code, err := gen.generateRandom()
	if err != nil {
		t.Fatalf("generateRandom failed: %v", err)
	}

	// Check length
	if len(code) != 7 {
		t.Errorf("Expected length 7, got %d", len(code))
	}

	// Check all characters are from charset
	for _, char := range code {
		if !strings.ContainsRune(charset, char) {
			t.Errorf("Generated code contains invalid character: %c", char)
		}
	}
}

func TestGenerateUnique(t *testing.T) {
	checker := NewMockChecker([]string{})
	gen := NewShortCodeGenerator(checker, 7)

	code, err := gen.Generate(context.Background())
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(code) != 7 {
		t.Errorf("Expected length 7, got %d", len(code))
	}

	// Should check existence once
	if checker.checkCount != 1 {
		t.Errorf("Expected 1 existence check, got %d", checker.checkCount)
	}
}

func TestGenerateWithCollision(t *testing.T) {
	// Pre-populate with many codes to increase collision probability
	existingCodes := []string{}
	for i := 0; i < 100; i++ {
		existingCodes = append(existingCodes, generateTestCode(i))
	}

	checker := NewMockChecker(existingCodes)
	gen := NewShortCodeGenerator(checker, 7)

	code, err := gen.Generate(context.Background())
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Should eventually find a unique code
	if checker.existingCodes[code] {
		t.Errorf("Generated code %s already exists", code)
	}
}

func TestGenerateMultipleUnique(t *testing.T) {
	checker := NewMockChecker([]string{})
	gen := NewShortCodeGenerator(checker, 7)

	generated := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code, err := gen.Generate(context.Background())
		if err != nil {
			t.Fatalf("Generate failed on iteration %d: %v", i, err)
		}

		if generated[code] {
			t.Errorf("Generated duplicate code: %s", code)
		}
		generated[code] = true

		// Add to mock checker to simulate real scenario
		checker.existingCodes[code] = true
	}
}

func TestValidateCustomCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{"valid lowercase", "abc123", false},
		{"valid uppercase", "ABC123", false},
		{"valid mixed", "aBc123", false},
		{"valid min length", "abc", false},
		{"too short", "ab", true},
		{"too long", "abcdefghijklmnopqrstu", true},
		{"invalid character space", "abc 123", true},
		{"invalid character dash", "abc-123", true},
		{"invalid character underscore", "abc_123", true},
		{"invalid character special", "abc@123", true},
		{"valid max length", "abcdefghijklmnopqrst", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCustomCode(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCustomCode(%q) error = %v, wantErr %v", tt.code, err, tt.wantErr)
			}
		})
	}
}

func TestDefaultLength(t *testing.T) {
	gen := NewShortCodeGenerator(nil, 0)
	if gen.length != defaultLength {
		t.Errorf("Expected default length %d, got %d", defaultLength, gen.length)
	}

	gen = NewShortCodeGenerator(nil, -5)
	if gen.length != defaultLength {
		t.Errorf("Expected default length %d for negative input, got %d", defaultLength, gen.length)
	}
}

// Helper function to generate test codes
func generateTestCode(seed int) string {
	code := ""
	for i := 0; i < 7; i++ {
		code += string(charset[(seed*7+i)%len(charset)])
	}
	return code
}
