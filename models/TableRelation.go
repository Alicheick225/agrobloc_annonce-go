package models

import "github.com/google/uuid"

type User struct {
	ID  uuid.UUID `json:"id" gorm:"type:uuid"`
	Nom string    `json:"nom"`
}

func (User) TableName() string {
	return "users" // Correspond à la table dans la base de données
}

type TypeCulture struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid"`
	Libelle string    `json:"libelle"`
}

func (TypeCulture) TableName() string {
	return "type_culture" //  Correspond à la table dans la base de données
}

type Parcelle struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Adresse string    `json:"adresse"`
	Surface string    `json:"surface"`
}

func (Parcelle) TableName() string {
	return "parcelle" //  Correspond à la table dans la base de données
}
