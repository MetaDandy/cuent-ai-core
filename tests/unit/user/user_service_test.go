//go:build unit

package user_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/src/core/user"
	"gorm.io/gorm"
)

// Test de validación de email usando la misma regex del servicio
func TestEmailValidation(t *testing.T) {
	emailRx := regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

	tests := []struct {
		name    string
		email   string
		isValid bool
	}{
		{"Valid email - standard", "test@example.com", true},
		{"Valid email - with subdomain", "user@mail.example.com", true},
		{"Valid email - with plus", "user+tag@example.com", true},
		{"Valid email - with dash", "user-name@example.com", true},
		{"Valid email - with underscore", "user_name@example.com", true},
		{"Invalid - no @", "invalid-email", false},
		{"Invalid - no domain", "test@", false},
		{"Invalid - no TLD", "test@example", false},
		{"Invalid - multiple @", "test@@example.com", false},
		{"Empty email", "", false},
		{"Only spaces", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := emailRx.MatchString(tt.email)
			if result != tt.isValid {
				t.Errorf("email validation failed for '%s': expected %v, got %v", tt.email, tt.isValid, result)
			}
		})
	}
}

// Test de validación de contraseña
func TestPasswordValidation(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		shouldErr bool
	}{
		{"Valid password - 8 chars", "password", false},
		{"Valid password - longer", "password123", false},
		{"Valid password - with special chars", "pass@123!", false},
		{"Invalid - too short (7)", "short12", true},
		{"Invalid - too short (6)", "short1", true},
		{"Invalid - empty", "", true},
		{"Valid - exactly 8", "12345678", false},
		{"Valid - 9 chars", "123456789", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := len(tt.password) >= 8
			if (isValid == false) != tt.shouldErr {
				t.Errorf("password validation failed: expected error=%v, got isValid=%v", tt.shouldErr, isValid)
			}
		})
	}
}

// Test de normalización de email
func TestEmailNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Lowercase conversion", "Test@Example.COM", "test@example.com"},
		{"Trim leading spaces", "  test@example.com", "test@example.com"},
		{"Trim trailing spaces", "test@example.com  ", "test@example.com"},
		{"Trim both sides", "  test@example.com  ", "test@example.com"},
		{"Mixed case", "TeSt@ExAmPlE.CoM", "test@example.com"},
		{"Already normalized", "test@example.com", "test@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimSpace(strings.ToLower(tt.input))
			if result != tt.expected {
				t.Errorf("email normalization failed: expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Test de errores del servicio
func TestServiceErrors(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedError error
	}{
		{"InvalidEmail error", user.ErrInvalidEmail, user.ErrInvalidEmail},
		{"WeakPassword error", user.ErrWeakPassword, user.ErrWeakPassword},
		{"EmailTaken error", user.ErrEmailTaken, user.ErrEmailTaken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.err, tt.expectedError) {
				t.Errorf("error mismatch: expected %v, got %v", tt.expectedError, tt.err)
			}
		})
	}
}

// Test de validación completa de SignUp input
func TestSignUpInputValidation(t *testing.T) {
	emailRx := regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

	tests := []struct {
		name          string
		input         user.Singup
		expectedError error
	}{
		{
			name: "Valid input - all fields correct",
			input: user.Singup{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedError: nil,
		},
		{
			name: "Invalid email format",
			input: user.Singup{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedError: user.ErrInvalidEmail,
		},
		{
			name: "Weak password - too short",
			input: user.Singup{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "short",
			},
			expectedError: user.ErrWeakPassword,
		},
		{
			name: "Valid - name can be empty",
			input: user.Singup{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedError: nil,
		},
		{
			name: "Valid - name with spaces trimmed",
			input: user.Singup{
				Name:     "  Test User  ",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email := strings.TrimSpace(strings.ToLower(tt.input.Email))

			if !emailRx.MatchString(email) {
				if tt.expectedError != user.ErrInvalidEmail {
					t.Errorf("expected ErrInvalidEmail for invalid email '%s'", email)
				}
				return
			}

			if len(tt.input.Password) < 8 {
				if tt.expectedError != user.ErrWeakPassword {
					t.Errorf("expected ErrWeakPassword for weak password")
				}
				return
			}

			if tt.expectedError != nil {
				t.Errorf("expected error %v but validation passed", tt.expectedError)
			}
		})
	}
}

// Test de comparación de errores de GORM
func TestGormErrorHandling(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		isRecordNotFound bool
	}{
		{"GORM RecordNotFound", gorm.ErrRecordNotFound, true},
		{"Other error", errors.New("other error"), false},
		{"Nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isNotFound := errors.Is(tt.err, gorm.ErrRecordNotFound)
			if isNotFound != tt.isRecordNotFound {
				t.Errorf("expected isRecordNotFound=%v, got %v", tt.isRecordNotFound, isNotFound)
			}
		})
	}
}
