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

	Status State `gorm:"type:state;default:'PENDING'"`

	UserID         uuid.UUID
	User           User
	SubscriptionID uuid.UUID
	Subscription   Subscription

	Payments []Payment `gorm:"foreignKey:UserSuscribedID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
