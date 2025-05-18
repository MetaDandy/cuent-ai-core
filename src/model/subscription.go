package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Name       string    `gorm:"not null"`
	Cuentokens uint
	Duration   time.Time

	// Poner un precio luego de la monetización

	UsersSubscriptions []UserSubscribed `gorm:"foreignKey:SubscriptionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
