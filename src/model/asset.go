package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Asset struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Type       AudioLine `gorm:"type:audio_line;default:'TTS'"`
	Video_URL  string
	Audio_URL  string
	Line       string
	AudioState State   `gorm:"type:state;default:'PENDING'"`
	VideoState State   `gorm:"type:state;default:'PENDING'"`
	Duration   float64 `gorm:"not null"`
	Position   int     `gorm:"not null"`

	ScriptID uuid.UUID
	Script   Script

	GeneratedJobs []GeneratedJob `gorm:"foreignKey:AssetID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type AudioLine string

const (
	AudioTTS AudioLine = "TTS"
	AudioSFX AudioLine = "SFX"
)

func (a *AudioLine) Scan(v interface{}) error {
	if v == nil {
		*a = ""
		return nil
	}
	switch s := v.(type) {
	case string:
		*a = AudioLine(s)
	case []byte:
		*a = AudioLine(string(s))
	default:
		return fmt.Errorf("no se puede convertir %T a AudioLine", v)
	}
	return nil
}

func (a AudioLine) Value() (driver.Value, error) {
	return string(a), nil
}
