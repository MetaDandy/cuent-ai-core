//go:build e2e

package functional

import (
    "net/http"
    "testing"

    "github.com/MetaDandy/cuent-ai-core/tests/functional/config"
    "github.com/MetaDandy/cuent-ai-core/tests/functional/helpers"
)

// TC-A01: Acceso sin autenticación a recursos protegidos
func TestAuthorization_AccessWithoutToken(t *testing.T) {
    cfg := config.NewTestConfig()
    client := helpers.NewHTTPClient(cfg.BaseURL, t)

    protectedEndpoints := []struct {
        method string
        path   string
        name   string
    }{
        {http.MethodGet, "/projects", "Listar proyectos"},
        {http.MethodGet, "/scripts", "Listar scripts"},
        {http.MethodGet, "/users/profile", "Obtener perfil"},
        {http.MethodPost, "/projects", "Crear proyecto"},
    }

    for _, endpoint := range protectedEndpoints {
        t.Run(endpoint.name, func(t *testing.T) {
            var resp *helpers.Response
            var err error

            switch endpoint.method {
            case http.MethodGet:
                resp, err = client.GET(endpoint.path)
            case http.MethodPost:
                resp, err = client.POST(endpoint.path, map[string]string{})
            }

            if err != nil {
                t.Errorf("Error en petición: %v", err)
                return
            }

            if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
                t.Errorf("Esperado 401/403, recibido: %d", resp.StatusCode)
            }
        })
    }
}

// TC-A02: Acceso con token inválido
func TestAuthorization_InvalidToken(t *testing.T) {
    cfg := config.NewTestConfig()
    client := helpers.NewHTTPClient(cfg.BaseURL, t)

    client.SetAuthToken("invalid_token_12345")

    resp, err := client.GET("/users/profile")
    if err != nil {
        t.Fatalf("Error en petición: %v", err)
    }

    if resp.StatusCode != http.StatusUnauthorized {
        t.Errorf("Esperado status 401, recibido: %d", resp.StatusCode)
    }
}
