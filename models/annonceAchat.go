	package models

import (
	"github.com/google/uuid"
	"time"
)

type AnnonceAchat struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid"`
	TypeCultureID  uuid.UUID `json:"type_culture_id" gorm:"type:uuid"`
	Statut         string    `json:"statut"`
	Quantite       int       `json:"quantite"`
	CreatedAt time.Time `json:"créé_a"`
	// UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt omitted for now
}

func (AnnonceAchat) TableName() string {
	return "annonces_achat"
}

