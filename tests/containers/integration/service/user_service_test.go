//go:build containers

package service

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

// TestUserService_FindById verifica que se obtiene un usuario por ID
func TestUserService_FindById(t *testing.T) {
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

	// Crear servicio
	repo := user.NewRepository(testDB.DB)
	svc := user.NewService(repo)

	// Ejecutar test
	result, err := svc.FindById(testUser.ID.String())

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.Email, result.Email)
}

// TestUserService_SignUp verifica el flujo de registro
func TestUserService_SignUp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)
	require.NoError(t, err)

	// Pre-poblar suscripción "Free"
	err = fixtures.SeedSubscriptions(testDB.DB)
	require.NoError(t, err)

	repo := user.NewRepository(testDB.DB)
	svc := user.NewService(repo)

	// Preparar input
	input := &user.Singup{
		Name:     "New User",
		Email:    "newuser@example.com",
		Password: "securepassword123",
	}

	// Ejecutar test
	result, token, err := svc.SignUp(input)

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, token)
	assert.Equal(t, input.Email, result.Email)

	// Verificar que el usuario fue creado en BD
	var createdUser model.User
	findErr := testDB.DB.First(&createdUser, "email = ?", input.Email).Error
	assert.NoError(t, findErr)

	// Verificar que la suscripción Free fue asignada
	var userSub model.UserSubscribed
	subErr := testDB.DB.
		Preload("Subscription").
		Where("user_id = ?", createdUser.ID).
		First(&userSub).Error
	assert.NoError(t, subErr)
	assert.Equal(t, "Free", userSub.Subscription.Name)
}

// TestUserService_SignUp_InvalidEmail verifica validación de email
func TestUserService_SignUp_InvalidEmail(t *testing.T) {
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

	repo := user.NewRepository(testDB.DB)
	svc := user.NewService(repo)

	// Email inválido
	input := &user.Singup{
		Name:     "New User",
		Email:    "invalid-email",
		Password: "securepassword123",
	}

	// Ejecutar test
	result, token, err := svc.SignUp(input)

	// Verificaciones
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, token)
}

// TestUserService_SignUp_WeakPassword verifica validación de contraseña
func TestUserService_SignUp_WeakPassword(t *testing.T) {
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

	repo := user.NewRepository(testDB.DB)
	svc := user.NewService(repo)

	// Contraseña débil (< 8 caracteres)
	input := &user.Singup{
		Name:     "New User",
		Email:    "test@example.com",
		Password: "short",
	}

	// Ejecutar test
	result, token, err := svc.SignUp(input)

	// Verificaciones
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, token)
}

// TestUserService_SignUp_EmailTaken verifica que no permite email duplicado
func TestUserService_SignUp_EmailTaken(t *testing.T) {
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

	// Crear usuario existente
	existingUser := fixtures.CreateTestUserWithEmail("taken@example.com")
	err = testDB.DB.Create(existingUser).Error
	require.NoError(t, err)

	repo := user.NewRepository(testDB.DB)
	svc := user.NewService(repo)

	// Intentar registrar con email existente
	input := &user.Singup{
		Name:     "Another User",
		Email:    "taken@example.com",
		Password: "securepassword123",
	}

	// Ejecutar test
	result, token, err := svc.SignUp(input)

	// Verificaciones
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, token)
}

// TestUserService_SignIn verifica el flujo de inicio de sesión
func TestUserService_SignIn(t *testing.T) {
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

	repo := user.NewRepository(testDB.DB)
	svc := user.NewService(repo)

	// Ejecutar test
	input := &user.Signin{
		Email:    "signin@example.com",
		Password: "password123",
	}
	result, token, err := svc.Signin(input)

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, token)
	assert.Equal(t, testUser.Email, result.Email)
}

// TestUserService_SignIn_InvalidPassword verifica contraseña incorrecta
func TestUserService_SignIn_InvalidPassword(t *testing.T) {
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

	repo := user.NewRepository(testDB.DB)
	svc := user.NewService(repo)

	// Intentar con contraseña incorrecta
	input := &user.Signin{
		Email:    "signin@example.com",
		Password: "wrongpassword",
	}
	result, token, err := svc.Signin(input)

	// Verificaciones
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Empty(t, token)
}

// TestUserService_FindAll verifica la obtención de todos los usuarios
func TestUserService_FindAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := setup.SetupTestDB(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	config.Migrate(testDB.DB)

	// Crear múltiples usuarios
	users := fixtures.CreateMultipleTestUsers(3)
	for _, u := range users {
		err = testDB.DB.Create(u).Error
		require.NoError(t, err)
	}

	repo := user.NewRepository(testDB.DB)
	svc := user.NewService(repo)

	// Ejecutar test
	opts := &helper.FindAllOptions{
		Limit:  10,
		Offset: 0,
	}
	result, err := svc.FindAll(opts)

	// Verificaciones
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(3), result.Total)
	assert.Len(t, result.Data, 3)
}
