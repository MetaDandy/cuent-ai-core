//go:build containers

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/cmd/api"
	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/fixtures"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/setup"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupFiberApp crea una aplicación Fiber configurada con BD de test
func setupFiberApp(t *testing.T, testDB *setup.PostgresTestDB) *fiber.App {
	// Set JWT_SECRET for tests
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-only-do-not-use-in-production")

	// Usar la BD de test en config
	config.DB = testDB.DB

	// Crear la aplicación Fiber
	app := fiber.New()

	// Setup container con BD de test
	container := src.SetupContainer()

	// Setup API
	api.SetupApi(app, container)

	return app
}

// TestUserHandler_SignUp verifica el endpoint POST /users/sign-up
func TestUserHandler_SignUp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Seed suscripción Free
	err = fixtures.SeedSubscriptions(testDB.DB)
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Preparar request
	body := bytes.NewBufferString(`{
		"name": "Test User",
		"email": "testuser@example.com",
		"password": "password123"
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/users/sign-up", body)
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Parsear respuesta
	var resp helper.Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)

	// Verificar que el usuario fue creado en BD
	var createdUser model.User
	findErr := testDB.DB.First(&createdUser, "email = ?", "testuser@example.com").Error
	assert.NoError(t, findErr)
}

// TestUserHandler_SignUp_InvalidEmail verifica validación de email
func TestUserHandler_SignUp_InvalidEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	err = fixtures.SeedSubscriptions(testDB.DB)
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Email inválido
	body := bytes.NewBufferString(`{
		"name": "Test User",
		"email": "invalid-email",
		"password": "password123"
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/users/sign-up", body)
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

// TestUserHandler_SignIn verifica el endpoint POST /users/sign-in
func TestUserHandler_SignIn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear usuario con contraseña hasheada
	hashedPassword, err := helper.HashPassword("password123")
	require.NoError(t, err)

	testUser := &model.User{
		ID:       fixtures.CreateTestUser().ID,
		Name:     "Test User",
		Email:    "signin@example.com",
		Password: string(hashedPassword),
	}
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Preparar request
	body := bytes.NewBufferString(`{
		"email": "signin@example.com",
		"password": "password123"
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/users/sign-in", body)
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Parsear respuesta
	var resp helper.Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

// TestUserHandler_SignIn_InvalidPassword verifica contraseña incorrecta
func TestUserHandler_SignIn_InvalidPassword(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear usuario
	hashedPassword, err := helper.HashPassword("correctpassword")
	require.NoError(t, err)

	testUser := &model.User{
		ID:       fixtures.CreateTestUser().ID,
		Name:     "Test User",
		Email:    "signin@example.com",
		Password: string(hashedPassword),
	}
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Contraseña incorrecta
	body := bytes.NewBufferString(`{
		"email": "signin@example.com",
		"password": "wrongpassword"
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/users/sign-in", body)
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

// TestUserHandler_GetProfile verifica el endpoint GET /users/profile con JWT
func TestUserHandler_GetProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear usuario
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Generar JWT
	token, err := helper.GenerateJwt(testUser.ID.String(), testUser.Email)
	require.NoError(t, err)

	// Preparar request con Authorization header
	req, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Parsear respuesta
	var resp helper.Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
}

// TestUserHandler_GetProfile_MissingToken verifica que rechaza sin JWT
func TestUserHandler_GetProfile_MissingToken(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Request sin Authorization header
	req, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// TestUserHandler_FindById verifica el endpoint GET /users/:id
func TestUserHandler_FindById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear usuario
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Generar JWT
	token, err := helper.GenerateJwt(testUser.ID.String(), testUser.Email)
	require.NoError(t, err)

	// Preparar request
	req, _ := http.NewRequest("GET", "/api/v1/users/"+testUser.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	// Parsear respuesta
	var resp helper.Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
}
