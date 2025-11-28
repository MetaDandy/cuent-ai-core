//go:build containers

package repository

import (
	"context"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/fixtures"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProjectRepository_Create verifica que se crea un proyecto correctamente
func TestProjectRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)

	// Crear y guardar usuario
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	// Ejecutar test
	testProject := fixtures.CreateTestProject(user.ID)
	repo := project.NewRepository(testDB.DB)
	err = repo.Create(testProject)

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que fue guardado
	var savedProject model.Project
	findErr := testDB.DB.First(&savedProject, "id = ?", testProject.ID).Error
	assert.NoError(t, findErr)
	assert.Equal(t, testProject.Name, savedProject.Name)
}

// TestProjectRepository_FindById verifica que se encuentra un proyecto por ID
func TestProjectRepository_FindById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)

	// Crear user y project
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(user.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	// Ejecutar test
	repo := project.NewRepository(testDB.DB)
	found, err := repo.FindById(testProject.ID.String())

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, testProject.ID, found.ID)
	assert.Equal(t, testProject.Name, found.Name)
}

// TestProjectRepository_Update verifica que se actualiza un proyecto correctamente
func TestProjectRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)

	// Crear user y project
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(user.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	// Actualizar proyecto
	testProject.Name = "Updated Project"
	testProject.State = model.StateActive

	repo := project.NewRepository(testDB.DB)
	err = repo.Update(testProject)

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que fue actualizado
	var updatedProject model.Project
	findErr := testDB.DB.First(&updatedProject, "id = ?", testProject.ID).Error
	assert.NoError(t, findErr)
	assert.Equal(t, "Updated Project", updatedProject.Name)
	assert.Equal(t, model.StateActive, updatedProject.State)
}

// TestProjectRepository_SoftDelete verifica que se borra suavemente un proyecto
func TestProjectRepository_SoftDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)

	// Crear user y project
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(user.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	// Ejecutar soft delete
	repo := project.NewRepository(testDB.DB)
	err = repo.SoftDelete(testProject.ID.String())

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que no aparece en búsqueda normal
	var deletedProject model.Project
	findErr := testDB.DB.First(&deletedProject, "id = ?", testProject.ID).Error
	assert.Error(t, findErr)

	// Verificar que aparece con Unscoped
	findUnscopedErr := testDB.DB.Unscoped().First(&deletedProject, "id = ?", testProject.ID).Error
	assert.NoError(t, findUnscopedErr)
	assert.NotNil(t, deletedProject.DeletedAt.Time)
}

// TestProjectRepository_Restore verifica que se restaura un proyecto eliminado
func TestProjectRepository_Restore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)

	// Crear user y project
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(user.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	// Soft delete
	repo := project.NewRepository(testDB.DB)
	err = repo.SoftDelete(testProject.ID.String())
	require.NoError(t, err)

	// Ejecutar restore
	err = repo.Restore(testProject.ID.String())

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que aparece en búsqueda normal nuevamente
	var restoredProject model.Project
	findErr := testDB.DB.First(&restoredProject, "id = ?", testProject.ID).Error
	assert.NoError(t, findErr)
	assert.False(t, restoredProject.DeletedAt.Valid, "DeletedAt should be NULL after restore")
}

// TestProjectRepository_FindAll verifica que se obtienen todos los proyectos
func TestProjectRepository_FindAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)

	// Crear user y projects
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	projects := fixtures.CreateMultipleTestProjects(user.ID, 5)
	for _, p := range projects {
		err = testDB.DB.Create(p).Error
		require.NoError(t, err)
	}

	// Ejecutar test
	repo := project.NewRepository(testDB.DB)
	opts := &helper.FindAllOptions{
		Limit:  10,
		Offset: 0,
	}
	found, total, err := repo.FindAll(opts)

	// Verificaciones
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, found, 5)
}

// TestProjectRepository_FindAll_Pagination verifica la paginación
func TestProjectRepository_FindAll_Pagination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)

	// Crear user y 15 projects
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	projects := fixtures.CreateMultipleTestProjects(user.ID, 15)
	for _, p := range projects {
		err = testDB.DB.Create(p).Error
		require.NoError(t, err)
	}

	repo := project.NewRepository(testDB.DB)

	// Página 1: 10 resultados
	opts := &helper.FindAllOptions{
		Limit:  10,
		Offset: 0,
	}
	found, total, err := repo.FindAll(opts)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, found, 10)

	// Página 2: 5 resultados
	opts = &helper.FindAllOptions{
		Limit:  10,
		Offset: 10,
	}
	found, total, err = repo.FindAll(opts)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, found, 5)
}
