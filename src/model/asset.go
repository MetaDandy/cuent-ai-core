package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Asset struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Type       string
	Video_URL  string
	Audio_URL  string
	Line       string
	AudioState State          `gorm:"type:state;default:'PENDING'"`
	VideoState State          `gorm:"type:state;default:'PENDING'"`
	Duration   datatypes.Time `gorm:"type:time"`

	ScriptID uuid.UUID
	Script   Script

	GeneratedJobs []GeneratedJob `gorm:"foreignKey:AssetID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
