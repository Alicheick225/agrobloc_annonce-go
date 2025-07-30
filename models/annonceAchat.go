package models

import (
	"time"

	"github.com/google/uuid"
)

type AnnonceAchat struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid"`
	TypeCultureID uuid.UUID `json:"type_culture_id" gorm:"type:uuid"`
	Statut        string    `json:"statut"`
	Description   string    `json:"description"`
	Quantite      float64   `json:"quantite"`
	CreatedAt     time.Time `json:"cree_at"`
	// UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User        User        `json:"users" gorm:"foreignKey:UserID"`
	TypeCulture TypeCulture `json:"type_culture" gorm:"foreignKey:TypeCultureID"`
}

func (AnnonceAchat) TableName() string {
	return "annonces_achat"
}
