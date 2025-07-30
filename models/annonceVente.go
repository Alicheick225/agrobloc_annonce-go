package models

import (
	"time"

	"github.com/google/uuid"
)

type AnnonceVente struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid"`
	TypeCultureID uuid.UUID `json:"type_culture_id" gorm:"type:uuid"`
	ParcelleID    uuid.UUID `json:"parcelle_id" gorm:"type:uuid"`
	Photo         string    `json:"photo"`
	Statut        string    `json:"statut"`
	Quantite      float64   `json:"quantite"`
	PrixKg        float64   `json:"prix_kg"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"créé_a"`
	// UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt omitted for now

	// Relations
	User        User        `json:"users" gorm:"foreignKey:UserID"`
	TypeCulture TypeCulture `json:"type_culture" gorm:"foreignKey:TypeCultureID"`
	Parcelle    Parcelle    `json:"parcelle" gorm:"foreignKey:ParcelleID"`
}

func (AnnonceVente) TableName() string {
	return "annonces_vente"
}
