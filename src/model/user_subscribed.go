package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSubscribed struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;"`
	TokensRemaining uint      `gorm:"column:total_cuentokens;not null"`
	StartDate       time.Time
	EndDate         time.Time
	CanceledAt      *time.Time

	UserID         uuid.UUID
	User           User
	SubscriptionID uuid.UUID
	Subscription   Subscription

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
