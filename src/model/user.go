package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Name             string    `gorm:"not null"`
	Email            string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password         string    `gorm:"type:varchar(100)"`
	StripeCustomerID string

	Projects           []Project        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UsersSubscriptions []UserSubscribed `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
