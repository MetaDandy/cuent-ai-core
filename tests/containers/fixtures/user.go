//go:build containers

package fixtures

import (
	"fmt"
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTestUser crea un usuario de prueba
func CreateTestUser() *model.User {
	return &model.User{
		ID:        uuid.New(),
		Name:      "Test User",
		Email:     "testuser@example.com",
		Password:  "hashedpassword123", // En tests reales, usar helper.HashPassword
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestUserWithEmail crea un usuario de prueba con email personalizado
func CreateTestUserWithEmail(email string) *model.User {
	return &model.User{
		ID:        uuid.New(),
		Name:      "Test User",
		Email:     email,
		Password:  "hashedpassword123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateMultipleTestUsers crea múltiples usuarios de prueba
func CreateMultipleTestUsers(count int) []*model.User {
	users := make([]*model.User, count)
	for i := 0; i < count; i++ {
		users[i] = &model.User{
			ID:        uuid.New(),
			Name:      "Test User",
			Email:     fmt.Sprintf("testuser%d@example.com", i),
			Password:  "hashedpassword123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	return users
}

// SeedTestUser inserta un usuario de prueba en la BD
func SeedTestUser(db *gorm.DB) error {
	user := CreateTestUser()
	return db.Create(user).Error
}

// SeedMultipleTestUsers inserta múltiples usuarios de prueba en la BD
func SeedMultipleTestUsers(count int) func(*gorm.DB) error {
	return func(db *gorm.DB) error {
		users := CreateMultipleTestUsers(count)
		return db.CreateInBatches(users, 100).Error
	}
}

// CleanUsers elimina todos los usuarios (útil para reset entre tests)
func CleanUsers(db *gorm.DB) error {
	return db.Exec("TRUNCATE TABLE users CASCADE").Error
}
