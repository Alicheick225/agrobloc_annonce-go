package models

import (
	"time"

	"github.com/google/uuid"
)

type AnnoncePrefinancement struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid"`
	TypeCultureID uuid.UUID `json:"type_culture_id" gorm:"type:uuid"`
	ParcelleID    uuid.UUID `json:"parcelle_id" gorm:"type:uuid"`
	Statut        string    `json:"status" gorm:"column:status"`
	Description   string    `json:"description"`
	Montant       int       `json:"montant"`
	CreatedAt     time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt omitted for now
}

func (AnnoncePrefinancement) TableName() string {
	return "annonces_prefinancement"
}
