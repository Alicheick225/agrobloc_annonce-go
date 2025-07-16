package models

import (
	"github.com/google/uuid"
)

type AnnoncePrefinancement struct {
	ID                    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID                uuid.UUID `json:"user_id" gorm:"type:uuid"`
	TypeCultureID         uuid.UUID `json:"type_culture_id " gorm:"type:uuid"`
	ParcelleID            uuid.UUID `json:"parcelle_id" gorm:"type:uuid"`
	Statut                string    `json:"statut" gorm:"column:statut"`
	Description           string    `json:"description"`
	MontantPrefinancement float64   `json:"montant_pref"  gorm:"column:montant_pref"`
	Prix                  float64   `json:"prix_kg_pref"  gorm:"column:prix_kg_pref"`
	Quantite              float64   `json:"quantite"  gorm:"column:quantite"`

	// UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt omitted for now

	// Relations
	User        User        `json:"users" gorm:"foreignKey:UserID"`
	TypeCulture TypeCulture `json:"type_culture" gorm:"foreignKey:TypeCultureID"`
	Parcelle    Parcelle    `json:"parcelle" gorm:"foreignKey:ParcelleID"`
}

func (AnnoncePrefinancement) TableName() string {
	return "annonces_prefinancement"

}
