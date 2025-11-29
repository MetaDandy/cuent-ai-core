//go:build containers

package fixtures

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTestProject crea un proyecto de prueba
func CreateTestProject(userID uuid.UUID) *model.Project {
	return &model.Project{
		ID:          uuid.New(),
		Name:        "Test Project",
		Description: "Test project description",
		Cuentokens:  "1000",
		State:       model.StatePending,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestProjectWithName crea un proyecto con nombre personalizado
func CreateTestProjectWithName(userID uuid.UUID, name string) *model.Project {
	return &model.Project{
		ID:          uuid.New(),
		Name:        name,
		Description: "Test project description",
		Cuentokens:  "1000",
		State:       model.StatePending,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestProjectWithState crea un proyecto con estado específico
func CreateTestProjectWithState(userID uuid.UUID, state model.State) *model.Project {
	return &model.Project{
		ID:          uuid.New(),
		Name:        "Test Project",
		Description: "Test project description",
		Cuentokens:  "1000",
		State:       state,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateMultipleTestProjects crea múltiples proyectos de prueba
func CreateMultipleTestProjects(userID uuid.UUID, count int) []*model.Project {
	projects := make([]*model.Project, count)
	for i := 0; i < count; i++ {
		projects[i] = &model.Project{
			ID:          uuid.New(),
			Name:        "Test Project",
			Description: "Test project description",
			Cuentokens:  "1000",
			State:       model.StatePending,
			UserID:      userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}
	return projects
}

// SeedProject inserta un proyecto de prueba en la BD
func SeedProject(userID uuid.UUID) func(*gorm.DB) error {
	return func(db *gorm.DB) error {
		project := CreateTestProject(userID)
		return db.Create(project).Error
	}
}

// SeedMultipleProjects inserta múltiples proyectos de prueba en la BD
func SeedMultipleProjects(userID uuid.UUID, count int) func(*gorm.DB) error {
	return func(db *gorm.DB) error {
		projects := CreateMultipleTestProjects(userID, count)
		return db.CreateInBatches(projects, 100).Error
	}
}

// CleanProjects elimina todos los proyectos
func CleanProjects(db *gorm.DB) error {
	return db.Exec("TRUNCATE TABLE projects CASCADE").Error
}
