//go:build containers

package service

import (
	"context"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/helper"
	usercore "github.com/MetaDandy/cuent-ai-core/src/core/user"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/src/modules/project"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/fixtures"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProjectService_FindByID verifica que se obtiene un proyecto por ID
func TestProjectService_FindByID(t *testing.T) {
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
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(user.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	// Crear servicio
	repo := project.NewRepository(testDB.DB)
	svc := project.NewService(repo, nil)

	// Ejecutar test
	result, err := svc.FindByID(testProject.ID.String())

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testProject.ID.String(), result.ID)
	assert.Equal(t, testProject.Name, result.Name)
}

// TestProjectService_Create verifica la creación de un proyecto
func TestProjectService_Create(t *testing.T) {
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
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	repo := project.NewRepository(testDB.DB)
	userRepo := usercore.NewRepository(testDB.DB)
	svc := project.NewService(repo, userRepo)

	// Preparar proyecto
	testProject := fixtures.CreateTestProject(user.ID)
	input := &project.ProjectCreate{
		Name:        testProject.Name,
		Description: testProject.Description,
		UserId:      user.ID.String(),
	}

	// Ejecutar test
	result, err := svc.Create(input)

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)

	// Verificar que fue guardado en BD
	var savedProject model.Project
	findErr := testDB.DB.First(&savedProject, "id = ?", result.ID).Error
	assert.NoError(t, findErr)
}

// TestProjectService_Update verifica la actualización de un proyecto
func TestProjectService_Update(t *testing.T) {
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
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(user.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	repo := project.NewRepository(testDB.DB)
	userRepo := usercore.NewRepository(testDB.DB)
	svc := project.NewService(repo, userRepo)

	// Preparar actualización
	updatedName := "Updated Project"
	updatedDesc := "Updated description"
	input := &project.ProjectUpdate{
		Name:        &updatedName,
		Description: &updatedDesc,
	}

	// Ejecutar test
	result, err := svc.Update(testProject.ID.String(), input)

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, *input.Name, result.Name)

	// Verificar que fue actualizado en BD
	var updatedProject model.Project
	findErr := testDB.DB.First(&updatedProject, "id = ?", testProject.ID).Error
	assert.NoError(t, findErr)
	assert.Equal(t, "Updated Project", updatedProject.Name)
}

// TestProjectService_Delete verifica la eliminación de un proyecto
func TestProjectService_Delete(t *testing.T) {
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
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	testProject := fixtures.CreateTestProject(user.ID)
	err = testDB.DB.Create(testProject).Error
	require.NoError(t, err)

	repo := project.NewRepository(testDB.DB)
	userRepo := usercore.NewRepository(testDB.DB)
	svc := project.NewService(repo, userRepo)

	// Ejecutar test
	_, err = svc.SoftDelete(testProject.ID.String())

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que fue eliminado (soft delete)
	var deletedProject model.Project
	findErr := testDB.DB.First(&deletedProject, "id = ?", testProject.ID).Error
	assert.Error(t, findErr)

	// Verificar con Unscoped
	var unscopedProject model.Project
	unscopedErr := testDB.DB.Unscoped().First(&unscopedProject, "id = ?", testProject.ID).Error
	assert.NoError(t, unscopedErr)
}

// TestProjectService_GetAll verifica la obtención de todos los proyectos de un usuario
func TestProjectService_GetAll(t *testing.T) {
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
	user := fixtures.CreateTestUser()
	err = testDB.DB.Create(user).Error
	require.NoError(t, err)

	projects := fixtures.CreateMultipleTestProjects(user.ID, 5)
	for _, p := range projects {
		err = testDB.DB.Create(p).Error
		require.NoError(t, err)
	}

	repo := project.NewRepository(testDB.DB)
	userRepo := usercore.NewRepository(testDB.DB)
	svc := project.NewService(repo, userRepo)

	// Ejecutar test
	opts := &helper.FindAllOptions{
		Limit:  10,
		Offset: 0,
	}
	result, err := svc.FindAll(opts)

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(5), result.Total)
	assert.Len(t, result.Data, 5)
}
