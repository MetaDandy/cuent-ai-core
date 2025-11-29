//go:build containers

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/fixtures"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProjectHandler_CreateProject verifica el endpoint POST /projects
func TestProjectHandler_CreateProject(t *testing.T) {
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
	body := bytes.NewBufferString(`{
		"name": "New Project",
		"description": "Test project",
		"user_id": "` + testUser.ID.String() + `"
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/projects", body)
	req.Header.Set("Content-Type", "application/json")
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

	// Verificar que el proyecto fue creado en BD
	var createdProject model.Project
	findErr := testDB.DB.First(&createdProject, "name = ?", "New Project").Error
	assert.NoError(t, findErr)
	assert.Equal(t, testUser.ID, createdProject.UserID)
}

// TestProjectHandler_GetProject verifica el endpoint GET /projects/:id
func TestProjectHandler_GetProject(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear usuario y proyecto
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(testUser.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Generar JWT
	token, err := helper.GenerateJwt(testUser.ID.String(), testUser.Email)
	require.NoError(t, err)

	// Preparar request
	req, _ := http.NewRequest("GET", "/api/v1/projects/"+testProject.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Parsear respuesta
	var resp helper.Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
}

// TestProjectHandler_GetProjects verifica el endpoint GET /projects
func TestProjectHandler_GetProjects(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear usuario y múltiples proyectos
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	projects := fixtures.CreateMultipleTestProjects(testUser.ID, 3)
	for _, p := range projects {
		err = testDB.DB.Create(p).Error
		require.NoError(t, err)
	}

	app := setupFiberApp(t, testDB)

	// Generar JWT
	token, err := helper.GenerateJwt(testUser.ID.String(), testUser.Email)
	require.NoError(t, err)

	// Preparar request
	req, _ := http.NewRequest("GET", "/api/v1/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Parsear respuesta
	var resp helper.Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	assert.NoError(t, err)
}

// TestProjectHandler_UpdateProject verifica el endpoint PATCH /projects/:id
func TestProjectHandler_UpdateProject(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear usuario y proyecto
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(testUser.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Generar JWT
	token, err := helper.GenerateJwt(testUser.ID.String(), testUser.Email)
	require.NoError(t, err)

	// Preparar request
	body := bytes.NewBufferString(`{
		"name": "Updated Project",
		"description": "Updated description",
		"cuentokens": "2000"
	}`)

	req, _ := http.NewRequest("PATCH", "/api/v1/projects/"+testProject.ID.String(), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Verificar que fue actualizado en BD
	var updatedProject model.Project
	findErr := testDB.DB.First(&updatedProject, "id = ?", testProject.ID).Error
	assert.NoError(t, findErr)
	assert.Equal(t, "Updated Project", updatedProject.Name)
}

// TestProjectHandler_DeleteProject verifica el endpoint DELETE /projects/:id
func TestProjectHandler_DeleteProject(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear usuario y proyecto
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(testUser.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	app := setupFiberApp(t, testDB)

	// Generar JWT
	token, err := helper.GenerateJwt(testUser.ID.String(), testUser.Email)
	require.NoError(t, err)

	// Preparar request
	req, _ := http.NewRequest("DELETE", "/api/v1/projects/"+testProject.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	// Verificar que fue eliminado (soft delete) en BD
	var deletedProject model.Project
	findErr := testDB.DB.First(&deletedProject, "id = ?", testProject.ID).Error
	assert.Error(t, findErr)
}

// TestProjectHandler_CreateProject_MissingToken verifica que rechaza sin JWT
func TestProjectHandler_CreateProject_MissingToken(t *testing.T) {
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

	// Preparar request sin Authorization header
	body := bytes.NewBufferString(`{
		"name": "New Project",
		"description": "Test project",
		"cuentokens": "1000"
	}`)

	req, _ := http.NewRequest("POST", "/api/v1/projects", body)
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar test
	res, err := app.Test(req, -1)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}
