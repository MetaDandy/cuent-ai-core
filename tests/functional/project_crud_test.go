//go:build e2e

package functional

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/tests/functional/config"
	"github.com/MetaDandy/cuent-ai-core/tests/functional/helpers"
)

// TC-P02: Listar todos los proyectos del usuario autenticado
func TestProjectFlow_ListProjects(t *testing.T) {
	cfg := config.NewTestConfig()
	client, _ := helpers.SetupAuthenticatedClient(t, cfg)

	resp, err := client.GET("/projects")
	if err != nil {
		t.Fatalf("Error listando proyectos: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Esperado status 200, recibido: %d", resp.StatusCode)
		return
	}

	var projectsWrapper struct {
		Data []helpers.ProjectResponse `json:"data"`
	}

	if err := resp.ParseJSON(&projectsWrapper); err != nil {
		t.Fatalf("Error parseando respuesta: %v", err)
	}

	// Puede estar vacío o tener proyectos, pero debe ser un array válido
	t.Logf("Total de proyectos encontrados: %d", len(projectsWrapper.Data))
}

// TC-P03: Actualizar un proyecto existente
func TestProjectFlow_UpdateProject(t *testing.T) {
	cfg := config.NewTestConfig()
	client, userID := helpers.SetupAuthenticatedClient(t, cfg)

	// Primero crear un proyecto
	createReq := helpers.CreateProjectRequest{
		Name:        helpers.GenerateProjectName(),
		Description: "Descripción original",
		UserID:      userID,
	}

	createResp, err := client.POST("/projects", createReq)
	if err != nil {
		t.Fatalf("Error creando proyecto: %v", err)
	}

	var createWrapper struct {
		Data helpers.ProjectResponse `json:"data"`
	}
	if err := createResp.ParseJSON(&createWrapper); err != nil {
		t.Fatalf("Error parseando respuesta de creación: %v", err)
	}
	projectID := createWrapper.Data.ID

	// Ahora actualizar el proyecto
	updateReq := map[string]interface{}{
		"name":        "Proyecto Actualizado E2E",
		"description": "Descripción modificada en test E2E",
	}

	updateResp, err := client.PATCH(fmt.Sprintf("/projects/%s", projectID), updateReq)
	if err != nil {
		t.Fatalf("Error actualizando proyecto: %v", err)
	}

	if updateResp.StatusCode != http.StatusOK {
		t.Errorf("Esperado status 200, recibido: %d. Body: %s", updateResp.StatusCode, string(updateResp.Body))
		return
	}

	var updateWrapper struct {
		Data helpers.ProjectResponse `json:"data"`
	}
	if err := updateResp.ParseJSON(&updateWrapper); err != nil {
		t.Fatalf("Error parseando respuesta de actualización: %v", err)
	}

	if updateWrapper.Data.Name != "Proyecto Actualizado E2E" {
		t.Errorf("Nombre no fue actualizado correctamente. Esperado: 'Proyecto Actualizado E2E', Recibido: '%s'", updateWrapper.Data.Name)
	}
}

// TC-P04: Obtener un proyecto por ID
func TestProjectFlow_GetProjectByID(t *testing.T) {
	cfg := config.NewTestConfig()
	client, userID := helpers.SetupAuthenticatedClient(t, cfg)

	// Crear un proyecto primero
	createReq := helpers.CreateProjectRequest{
		Name:        helpers.GenerateProjectName(),
		Description: "Proyecto para obtener por ID",
		UserID:      userID,
	}

	createResp, err := client.POST("/projects", createReq)
	if err != nil {
		t.Fatalf("Error creando proyecto: %v", err)
	}

	var createWrapper struct {
		Data helpers.ProjectResponse `json:"data"`
	}
	if err := createResp.ParseJSON(&createWrapper); err != nil {
		t.Fatalf("Error parseando respuesta: %v", err)
	}
	projectID := createWrapper.Data.ID

	// Obtener el proyecto por ID
	resp, err := client.GET(fmt.Sprintf("/projects/%s", projectID))
	if err != nil {
		t.Fatalf("Error obteniendo proyecto: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Esperado status 200, recibido: %d", resp.StatusCode)
		return
	}

	var projectWrapper struct {
		Data helpers.ProjectResponse `json:"data"`
	}
	if err := resp.ParseJSON(&projectWrapper); err != nil {
		t.Fatalf("Error parseando respuesta: %v", err)
	}

	if projectWrapper.Data.ID != projectID {
		t.Errorf("ID del proyecto no coincide. Esperado: %s, Recibido: %s", projectID, projectWrapper.Data.ID)
	}

	if projectWrapper.Data.Name != createReq.Name {
		t.Errorf("Nombre del proyecto no coincide. Esperado: %s, Recibido: %s", createReq.Name, projectWrapper.Data.Name)
	}
}

// TC-P05: Eliminación suave (soft delete) de un proyecto
func TestProjectFlow_SoftDeleteProject(t *testing.T) {
	cfg := config.NewTestConfig()
	client, userID := helpers.SetupAuthenticatedClient(t, cfg)

	// Crear un proyecto para eliminar
	createReq := helpers.CreateProjectRequest{
		Name:        helpers.GenerateProjectName(),
		Description: "Proyecto a eliminar",
		UserID:      userID,
	}

	createResp, err := client.POST("/projects", createReq)
	if err != nil {
		t.Fatalf("Error creando proyecto: %v", err)
	}

	var createWrapper struct {
		Data helpers.ProjectResponse `json:"data"`
	}
	if err := createResp.ParseJSON(&createWrapper); err != nil {
		t.Fatalf("Error parseando respuesta: %v", err)
	}
	projectID := createWrapper.Data.ID

	// Eliminar el proyecto
	deleteResp, err := client.DELETE(fmt.Sprintf("/projects/%s", projectID))
	if err != nil {
		t.Fatalf("Error eliminando proyecto: %v", err)
	}

	if deleteResp.StatusCode != http.StatusNoContent && deleteResp.StatusCode != http.StatusOK {
		t.Errorf("Esperado status 204 o 200, recibido: %d. Body: %s", deleteResp.StatusCode, string(deleteResp.Body))
	}
}

// TC-P06: Crear proyecto con datos incompletos (debe fallar)
func TestProjectFlow_CreateProjectInvalidData(t *testing.T) {
	cfg := config.NewTestConfig()
	client, _ := helpers.SetupAuthenticatedClient(t, cfg)

	// Intentar crear proyecto sin nombre (campo requerido)
	invalidReq := map[string]interface{}{
		"description": "Proyecto sin nombre",
	}

	resp, err := client.POST("/projects", invalidReq)
	if err != nil {
		t.Fatalf("Error en petición: %v", err)
	}

	// Debe fallar con 400 o 500
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		t.Errorf("No debería permitir crear proyecto sin datos válidos. Status: %d", resp.StatusCode)
	}

	t.Logf("Error esperado recibido con status: %d", resp.StatusCode)
}
