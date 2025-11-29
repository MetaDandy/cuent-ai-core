//go:build e2e

package helpers

import (
    "testing"

    "github.com/MetaDandy/cuent-ai-core/tests/functional/config"
)

// SetupAuthenticatedClient inicia sesión con el usuario admin configurado y devuelve
// un cliente con el token de autenticación listo para usar en pruebas.
func SetupAuthenticatedClient(t *testing.T, cfg *config.TestConfig) (*HTTPClient, string) {
    client := NewHTTPClient(cfg.BaseURL, t)

    loginReq := LoginRequest{
        Email:    cfg.AdminEmail,
        Password: cfg.AdminPassword,
    }

    resp, err := client.POST("/users/sign-in", loginReq)
    if err != nil {
        t.Fatalf("Error en login: %v", err)
    }

    var loginResp LoginResponse
    if err := resp.ParseJSON(&loginResp); err != nil {
        t.Fatalf("Error parseando login: %v", err)
    }

    client.SetAuthToken(loginResp.Token)
    return client, loginResp.Data.ID
}
