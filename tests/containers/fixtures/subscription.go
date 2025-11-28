//go:build containers

package fixtures

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTestSubscription crea una suscripción de prueba
func CreateTestSubscription(name string, tokens uint) *model.Subscription {
	return &model.Subscription{
		ID:         uuid.New(),
		Name:       name,
		Cuentokens: tokens,
		Price:      9.99,
		Duration:   time.Now().AddDate(0, 1, 0),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// CreateFreeSubscription crea la suscripción "Free" estándar
func CreateFreeSubscription() *model.Subscription {
	return &model.Subscription{
		ID:         uuid.New(),
		Name:       "Free",
		Cuentokens: 100,
		Price:      0,
		Duration:   time.Now().AddDate(0, 1, 0),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// CreatePremiumSubscription crea la suscripción "Premium" estándar
func CreatePremiumSubscription() *model.Subscription {
	return &model.Subscription{
		ID:         uuid.New(),
		Name:       "Premium",
		Cuentokens: 1000,
		Price:      29.99,
		Duration:   time.Now().AddDate(0, 1, 0),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// CreateUserSubscription crea una suscripción de usuario
func CreateUserSubscription(userID, subscriptionID uuid.UUID, status model.State) *model.UserSubscribed {
	now := time.Now()
	return &model.UserSubscribed{
		ID:              uuid.New(),
		UserID:          userID,
		SubscriptionID:  subscriptionID,
		Status:          status,
		TokensRemaining: 100,
		StartDate:       now,
		EndDate:         now.AddDate(0, 1, 0),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// SeedSubscriptions inserta suscripciones estándar en la BD
func SeedSubscriptions(db *gorm.DB) error {
	subscriptions := []model.Subscription{
		*CreateFreeSubscription(),
		*CreatePremiumSubscription(),
	}
	return db.Create(&subscriptions).Error
}

// CleanSubscriptions elimina todas las suscripciones
func CleanSubscriptions(db *gorm.DB) error {
	return db.Exec("TRUNCATE TABLE subscriptions CASCADE").Error
}
