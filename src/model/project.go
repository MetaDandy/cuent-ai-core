package model

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Name        string    `gorm:"not null"`
	Description string
	Cuentokens  string `gorm:"not null"`
	State       State  `gorm:"type:state;default:'PENDING'"`

	UserID uuid.UUID
	User   User

	Scripts []Script `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type State string

const (
	StatePending     State = "PENDING"
	StateActive      State = "ACTIVE"
	StateFinished    State = "FINISHED"
	StateRegenerated State = "REGENERATED"
	StateError       State = "ERROR"
)

func (s *State) Scan(v interface{}) error    { *s = State(v.(string)); return nil }
func (s State) Value() (driver.Value, error) { return string(s), nil }
