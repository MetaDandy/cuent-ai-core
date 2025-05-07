package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Script struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Prompt_Tokens     uint32    `gorm:"not null"`
	Completion_Tokens uint32    `gorm:"not null"`
	Total_Tokens      uint32    `gorm:"not null"`
	State             State     `gorm:"type:state;default:'PENDING'"`
	Text_Entry        string    `gorm:"not null"`
	Processed_Text    string    `gorm:"not null"`
	Total_Cost        float64   `gorm:"type:numeric(10,4);not null"`
	Mixed_Audio       string
	Mixed_Media       string

	ProjectID uuid.UUID
	Project   Project

	Assets []Asset `gorm:"foreignKey:ScriptID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
