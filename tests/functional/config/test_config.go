//go:build e2e

package config

import (
    "os"
    "time"
)

// TestConfig contiene la configuración para las pruebas E2E
type TestConfig struct {
    BaseURL       string
    Timeout       time.Duration
    AdminEmail    string
    AdminPassword string
    TestUserEmail string
    TestUserPass  string
}

// NewTestConfig crea una nueva configuración de pruebas desde variables de entorno
func NewTestConfig() *TestConfig {
    return &TestConfig{
        BaseURL:       getEnv("TEST_API_URL", "http://localhost:8000/api/v1"),
        Timeout:       30 * time.Second,
        AdminEmail:    getEnv("TEST_ADMIN_EMAIL", "admin@gmail.com"),
        AdminPassword: getEnv("TEST_ADMIN_PASSWORD", "changeme123"),
        TestUserEmail: getEnv("TEST_USER_EMAIL", "testuser@example.com"),
        TestUserPass:  getEnv("TEST_USER_PASSWORD", "testpass123"),
    }
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
