//go:build e2e

package functional

import (
    "net/http"
    "testing"

    "github.com/MetaDandy/cuent-ai-core/tests/functional/config"
    "github.com/MetaDandy/cuent-ai-core/tests/functional/helpers"
)

// TC-P01: Crear proyecto con datos válidos
func TestProjectFlow_CreateProject(t *testing.T) {
    cfg := config.NewTestConfig()
    client, userID := helpers.SetupAuthenticatedClient(t, cfg)
    testData := helpers.NewTestData()

    createReq := helpers.CreateProjectRequest{
        Name:        helpers.GenerateProjectName(),
        Description: testData.ProjectDescription,
        UserID:      userID,
    }

    resp, err := client.POST("/projects", createReq)
    if err != nil {
        t.Fatalf("Error creando proyecto: %v", err)
    }

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        t.Errorf("Esperado status 200 o 201, recibido: %d", resp.StatusCode)
        return
    }

    var wrapper struct {
        Data helpers.ProjectResponse `json:"data"`
    }
    if err := resp.ParseJSON(&wrapper); err != nil {
        t.Fatalf("Error parseando respuesta: %v", err)
    }

    projectResp := wrapper.Data

    if projectResp.ID == "" {
        t.Error("Project ID no fue retornado")
    }

    if projectResp.Name != createReq.Name {
        t.Errorf("Nombre esperado: %s, recibido: %s", createReq.Name, projectResp.Name)
    }
}
