package models

type ListeAnnonceAchat struct {
	ID                 string  `json:"id" gorm:"type:uuid;primaryKey"`
	Statut             string  `json:"statut"`
	Prix               float64 `json:"prix_kg"`
	Description        string  `json:"description"`
	Quantite           float64 `json:"quantite"`
	UserNom            string  `json:"nom"`
	TypeCultureLibelle string  `json:"libelle"`
}
 