package models

type ListeAnnonceVente struct {
	ID                 string  `json:"id" gorm:"type:uuid;primaryKey"`
	Photo              string  `json:"photo"`
	Statut             string  `json:"statut"`
	Description        string  `json:"description"`
	Quantite           float64 `json:"quantite"`
	PrixKg             float64 `json:"prix_kg"`
	UserNom            string  `json:"nom"`
	TypeCultureLibelle string  `json:"libelle"`
	ParcelleAdresse    string  `json:"adresse"`
	// ParcelleSuf     string  `json:"surface"`
}
