//go:build containers

package repository

import (
	"context"
	"testing"

	"github.com/MetaDandy/cuent-ai-core/config"
	"github.com/MetaDandy/cuent-ai-core/helper"
	"github.com/MetaDandy/cuent-ai-core/src/core/user"
	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/fixtures"
	"github.com/MetaDandy/cuent-ai-core/tests/containers/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserRepository_FindByEmail verifica que se encuentra un usuario por email
func TestUserRepository_FindByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	// Migrar esquema
	config.Migrate(testDB.DB)

	// Crear y guardar usuario
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	// Crear repositorio
	repo := user.NewRepository(testDB.DB)

	// Ejecutar test
	found, err := repo.FindByEmail(testUser.Email)

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, testUser.Email, found.Email)
	assert.Equal(t, testUser.ID, found.ID)
}

// TestUserRepository_FindByEmail_NotFound verifica que retorna error cuando no existe
func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	repo := user.NewRepository(testDB.DB)

	// Ejecutar test
	found, err := repo.FindByEmail("nonexistent@example.com")

	// Verificaciones
	assert.Error(t, err)
	// Cuando hay error record not found, el usuario devuelto es vacío
	if found != nil {
		assert.Equal(t, "", found.Email)
	}
}

// TestUserRepository_FindById verifica que se encuentra un usuario por ID
func TestUserRepository_FindById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear y guardar usuario
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	repo := user.NewRepository(testDB.DB)

	// Ejecutar test
	found, err := repo.FindById(testUser.ID.String())

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, testUser.ID, found.ID)
}

// TestUserRepository_Create verifica que se crea un usuario correctamente
func TestUserRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	repo := user.NewRepository(testDB.DB)
	testUser := fixtures.CreateTestUser()

	// Ejecutar test
	err = repo.Create(testUser)

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que fue guardado
	var savedUser model.User
	findErr := testDB.DB.First(&savedUser, "id = ?", testUser.ID).Error
	assert.NoError(t, findErr)
	assert.Equal(t, testUser.Email, savedUser.Email)
}

// TestUserRepository_Update verifica que se actualiza un usuario correctamente
func TestUserRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear y guardar usuario
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	repo := user.NewRepository(testDB.DB)

	// Actualizar usuario
	testUser.Name = "Updated Name"
	err = repo.Update(testUser)

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que fue actualizado
	var updatedUser model.User
	findErr := testDB.DB.First(&updatedUser, "id = ?", testUser.ID).Error
	assert.NoError(t, findErr)
	assert.Equal(t, "Updated Name", updatedUser.Name)
}

// TestUserRepository_SoftDelete verifica que se borra suavemente un usuario
func TestUserRepository_SoftDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear y guardar usuario
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	repo := user.NewRepository(testDB.DB)

	// Ejecutar soft delete
	err = repo.SoftDelete(testUser.ID.String())

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que no aparece en búsqueda normal
	var deletedUser model.User
	findErr := testDB.DB.First(&deletedUser, "id = ?", testUser.ID).Error
	assert.Error(t, findErr)

	// Verificar que aparece con Unscoped
	findUnscopedErr := testDB.DB.Unscoped().First(&deletedUser, "id = ?", testUser.ID).Error
	assert.NoError(t, findUnscopedErr)
	assert.NotNil(t, deletedUser.DeletedAt.Time)
}

// TestUserRepository_Restore verifica que se restaura un usuario eliminado
func TestUserRepository_Restore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Crear, guardar y borrar usuario
	testUser := fixtures.CreateTestUser()
	err = testDB.DB.Create(testUser).Error
	require.NoError(t, err)

	repo := user.NewRepository(testDB.DB)
	err = repo.SoftDelete(testUser.ID.String())
	require.NoError(t, err)

	// Ejecutar restore
	err = repo.Restore(testUser.ID.String())

	// Verificaciones
	assert.NoError(t, err)

	// Verificar que aparece en búsqueda normal nuevamente
	var restoredUser model.User
	findErr := testDB.DB.First(&restoredUser, "id = ?", testUser.ID).Error
	assert.NoError(t, findErr)
	assert.False(t, restoredUser.DeletedAt.Valid, "DeletedAt should be NULL after restore")
}

// TestUserRepository_FindAll verifica que se obtienen todos los usuarios
func TestUserRepository_FindAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	repo := user.NewRepository(testDB.DB)

	// Crear y guardar múltiples usuarios
	users := fixtures.CreateMultipleTestUsers(5)
	for _, u := range users {
		err = testDB.DB.Create(u).Error
		require.NoError(t, err)
	}

	// Ejecutar test
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
