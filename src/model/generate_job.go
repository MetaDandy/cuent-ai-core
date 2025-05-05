package model

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GeneratedJob struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Provider      Provider  `gorm:"type:provider;default:'ELEVENLAB'"`
	Model         string
	Token_Spent   string
	Input_Chats   string
	State         State `gorm:"type:state;default:'PENDING'"`
	Finished_At   time.Time
	Stared_At     time.Time
	Error_Message string
	Cost          float64

	AssetID uuid.UUID
	Asset   Asset

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Provider string

const (
	ProviderOpenAI    Provider = "OPENAI"
	ProviderGemini    Provider = "GEMINI"
	ProviderElevenlab Provider = "ELEVENLAB"
)

func (p *Provider) Scan(v interface{}) error    { *p = Provider(v.(string)); return nil }
func (p Provider) Value() (driver.Value, error) { return string(p), nil }
