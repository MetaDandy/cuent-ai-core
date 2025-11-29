//go:build e2e

package helpers

import (
    "fmt"
    "time"

    "github.com/google/uuid"
)

// TestData contiene datos de prueba reutilizables
type TestData struct {
    TestUserEmail string
    TestUserPass  string
    ProjectName        string
    ProjectDescription string
    ScriptText string
}

// NewTestData crea datos de prueba con valores únicos
func NewTestData() *TestData {
    timestamp := time.Now().Unix()
    return &TestData{
        TestUserEmail: fmt.Sprintf("testuser_%d@example.com", timestamp),
        TestUserPass:  "testpass123",
        ProjectName:        fmt.Sprintf("Test Project %d", timestamp),
        ProjectDescription: "Proyecto de prueba E2E para validación funcional",
        ScriptText: "EDIPO REY - Escena de prueba para generación de audio y video con IA",
    }
}

// LoginRequest representa una petición de login
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// LoginResponse representa la respuesta de login
type LoginResponse struct {
    Data struct {
        ID    string `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    } `json:"data"`
    Token   string `json:"token"`
    Message string `json:"message"`
}

// SignUpRequest representa una petición de registro
type SignUpRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// CreateProjectRequest representa una petición de creación de proyecto
type CreateProjectRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    UserID      string `json:"user_id"`
}

// ProjectResponse representa un proyecto en la respuesta
type ProjectResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    State       string `json:"state"`
    UserID      string `json:"user_id"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

// CreateScriptRequest representa una petición de creación de script
type CreateScriptRequest struct {
    TextEntry string `json:"text_entry"`
    ProjectID string `json:"project_id"`
}

// ScriptResponse representa un script en la respuesta
type ScriptResponse struct {
    ID        string `json:"id"`
    TextEntry string `json:"text_entry"`
    ProjectID string `json:"project_id"`
    State     string `json:"state"`
    CreatedAt string `json:"created_at"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

// GenerateUniqueEmail genera un email único para pruebas
func GenerateUniqueEmail() string {
    return fmt.Sprintf("test_%s@example.com", uuid.New().String()[:8])
}

// GenerateProjectName genera un nombre único de proyecto
func GenerateProjectName() string {
    return fmt.Sprintf("Test Project %s", uuid.New().String()[:8])
}
