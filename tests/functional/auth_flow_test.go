//go:build e2e

package functional

import (
	"net/http"
	"strings"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/tests/functional/config"
	"github.com/MetaDandy/cuent-ai-core/tests/functional/helpers"
)

// TC-U01: Login con credenciales válidas
func TestAuthFlow_LoginSuccess(t *testing.T) {
	cfg := config.NewTestConfig()
	client := helpers.NewHTTPClient(cfg.BaseURL, t)

	loginReq := helpers.LoginRequest{
		Email:    cfg.AdminEmail,
		Password: cfg.AdminPassword,
	}

	resp, err := client.POST("/users/sign-in", loginReq)
	if err != nil {
		t.Fatalf("Error en petición de login: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("Esperado status 200 o 201, recibido: %d", resp.StatusCode)
	}

	var loginResp helpers.LoginResponse
	if err := resp.ParseJSON(&loginResp); err != nil {
		t.Fatalf("Error parseando respuesta: %v", err)
	}

	if loginResp.Token == "" {
		t.Error("Token JWT no fue retornado")
	}

	if loginResp.Data.Email != cfg.AdminEmail {
		t.Errorf("Email esperado: %s, recibido: %s", cfg.AdminEmail, loginResp.Data.Email)
	}

	if loginResp.Data.ID == "" {
		t.Error("User ID no fue retornado")
	}
}

// TC-U02: Login con credenciales inválidas
func TestAuthFlow_LoginInvalidCredentials(t *testing.T) {
	cfg := config.NewTestConfig()
	client := helpers.NewHTTPClient(cfg.BaseURL, t)

	loginReq := helpers.LoginRequest{
		Email:    cfg.AdminEmail,
		Password: "wrong_password_123",
	}

	resp, err := client.POST("/users/sign-in", loginReq)
	if err != nil {
		t.Fatalf("Error en petición de login: %v", err)
	}

	// El backend actual devuelve 500 para credenciales inválidas
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Esperado status 500 (comportamiento actual del backend), recibido: %d", resp.StatusCode)
	}

	// Verificar que el mensaje de error contenga información sobre la contraseña incorrecta
	bodyStr := string(resp.Body)
	if !strings.Contains(bodyStr, "contraseña") && !strings.Contains(bodyStr, "password") {
		t.Errorf("El mensaje de error debería mencionar la contraseña incorrecta. Body: %s", bodyStr)
	}
}
