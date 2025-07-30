package models

type LiteAnnoncePrefinancement struct {
	ID              string  `json:"id" gorm:"type:uuid;primaryKey"`
	Statut          string  `json:"statut"`
	Description     string  `json:"description"`
	MontantPref     float64 `json:"montant_pref"`
	PrixKgPref      float64 `json:"prix_kg_pref"`
	Quantite        float64 `json:"quantite"`
	UserNom         string  `json:"nom"`
	TypeCultureLibelle  string  `json:"libelle"`
	ParcelleAdresse string  `json:"adresse"`
	ParcelleSuf     string  `json:"surface"`
}
