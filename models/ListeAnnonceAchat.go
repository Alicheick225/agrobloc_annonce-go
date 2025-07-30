package models

type ListeAnnonceAchat struct {
	ID                  string  `json:"id" gorm:"type:uuid;primaryKey"`
	Statut              string    `json:"statut"`
	Description         string    `json:"description"`
	Quantite            float64   `json:"quantite"`
	UserNom             string  `json:"nom"`
	TypeCultureLibelle  string  `json:"libelle"`
}