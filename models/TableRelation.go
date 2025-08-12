package models

import "github.com/google/uuid"

type User struct {
	ID  uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Nom string    `json:"nom"`
	// Email    string    `json:"email" gorm:"uniqueIndex"`
	// Phone    float64    `json:"numero_tel" gorm:"uniqueIndex"`
	// Password string    `json:"-"`
}

func (User) TableName() string {
	return "users"
}

type TypeCulture struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid"`
	Libelle string    `json:"libelle"`
}

func (TypeCulture) TableName() string {
	return "type_culture"
}

type Parcelle struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Adresse string    `json:"adresse"`
	Surface string    `json:"surface"`
}

func (Parcelle) TableName() string {
	return "parcelle"
}
