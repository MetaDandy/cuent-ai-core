package model

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey;"`
	UserID                string
	StripeSessionID       string
	StripePaymentIntentID string
	Amount                int
	Currency              string
	Status                State `gorm:"type:state;default:'PENDING'"`

	UserSuscribed   UserSubscribed
	UserSuscribedID string

	CreatedAt time.Time
	UpdatedAt time.Time
}
