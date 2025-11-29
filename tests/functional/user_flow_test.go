//go:build e2e

package functional

import (
	"net/http"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/tests/functional/config"
	"github.com/MetaDandy/cuent-ai-core/tests/functional/helpers"
)

// TC-U03: Registro de nuevo usuario con datos válidos
func TestUserFlow_SignUpSuccess(t *testing.T) {
	cfg := config.NewTestConfig()
	client := helpers.NewHTTPClient(cfg.BaseURL, t)
	testData := helpers.NewTestData()

	signupReq := helpers.SignUpRequest{
		Name:     "Test User E2E",
		Email:    helpers.GenerateUniqueEmail(),
		Password: testData.TestUserPass,
	}

	resp, err := client.POST("/users/sign-up", signupReq)
	if err != nil {
		t.Fatalf("Error en petición de registro: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Esperado status 200 o 201, recibido: %d. Body: %s", resp.StatusCode, string(resp.Body))
		return
	}

	var signupResp helpers.LoginResponse
	if err := resp.ParseJSON(&signupResp); err != nil {
		t.Fatalf("Error parseando respuesta: %v", err)
	}

	if signupResp.Token == "" {
		t.Error("Token JWT no fue retornado en el registro")
	}

	if signupResp.Data.Email != signupReq.Email {
		t.Errorf("Email esperado: %s, recibido: %s", signupReq.Email, signupResp.Data.Email)
	}
}

// TC-U04: Registro con datos duplicados (email ya existe)
func TestUserFlow_SignUpDuplicateEmail(t *testing.T) {
	cfg := config.NewTestConfig()
	client := helpers.NewHTTPClient(cfg.BaseURL, t)

	// Intentar registrar con el email del admin (que ya existe)
	signupReq := helpers.SignUpRequest{
		Name:     "Test Duplicate",
		Email:    cfg.AdminEmail,
		Password: "newpassword123",
	}

	resp, err := client.POST("/users/sign-up", signupReq)
	if err != nil {
		t.Fatalf("Error en petición de registro: %v", err)
	}

	// Debe fallar con error 400 o 500
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		t.Errorf("No debería permitir registrar email duplicado. Status: %d", resp.StatusCode)
	}
}

// TC-U05: Obtener perfil de usuario autenticado
func TestUserFlow_GetProfile(t *testing.T) {
	cfg := config.NewTestConfig()
	client, _ := helpers.SetupAuthenticatedClient(t, cfg)

	resp, err := client.GET("/users/profile")
	if err != nil {
		t.Fatalf("Error obteniendo perfil: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Esperado status 200 o 201, recibido: %d", resp.StatusCode)
		return
	}

	var profileWrapper struct {
		Data struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"data"`
	}

	if err := resp.ParseJSON(&profileWrapper); err != nil {
		t.Fatalf("Error parseando respuesta: %v", err)
	}

	if profileWrapper.Data.ID == "" {
		t.Error("User ID no fue retornado")
	}

	if profileWrapper.Data.Email != cfg.AdminEmail {
		t.Errorf("Email esperado: %s, recibido: %s", cfg.AdminEmail, profileWrapper.Data.Email)
	}
}
