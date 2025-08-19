package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Steph-business/annonce_de_vente/database"
	"github.com/Steph-business/annonce_de_vente/models"
)

// Fonction utilitaire : transforme une AnnonceVente en ListeAnnonceVente
func toAnnonceDTO(a models.AnnonceVente) models.ListeAnnonceVente {
	return models.ListeAnnonceVente{
		ID:                 a.ID.String(),
		Photo:              a.Photo,
		Statut:             a.Statut,
		Description:        a.Description,
		PrixKg:             a.PrixKg,
		Quantite:           a.Quantite,
		UserNom:            a.User.Nom,
		TypeCultureLibelle: a.TypeCulture.Libelle,
		ParcelleAdresse:    a.Parcelle.Adresse,
	}
}

// validateParcelleID vérifie que le parcelleID est valide et existe en base
// Returns true if parcelleID is valid or empty (for optional updates)
func validateParcelleID(c *gin.Context, parcelleIDStr string) (uuid.UUID, bool) {
	if parcelleIDStr == "" {
		// Allow empty ParcelleID - will use existing value
		return uuid.Nil, true
	}

	log.Printf("Validation ParcelleID : valeur reçue '%s'\n", parcelleIDStr)
	parcelleID, err := uuid.Parse(parcelleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ParcelleID invalide - format UUID attendu"})
		log.Printf("Validation ParcelleID échouée : parsing UUID invalide pour %s\n", parcelleIDStr)
		return uuid.Nil, false
	}
	var parcelle models.Parcelle
	if err := database.DB.First(&parcelle, "id = ?", parcelleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parcelle introuvable avec l'ID fourni"})
		log.Printf("Validation ParcelleID échouée : parcelle introuvable pour ID %s\n", parcelleID.String())
		return uuid.Nil, false
	}
	log.Printf("Validation ParcelleID réussie pour ID %s\n", parcelleID.String())
	return parcelleID, true
}

// Liste toutes les annonces avec filtres facultatifs
func GetAllAnnonceVente(c *gin.Context) {
	var annonces []models.AnnonceVente

	query := database.DB.Preload("User").Preload("TypeCulture").Preload("Parcelle")

	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if statut := c.Query("statut"); statut != "" {
		query = query.Where("statut = ?", statut)
	}
	if typeCultureID := c.Query("type_culture_id"); typeCultureID != "" {
		query = query.Where("type_culture_id = ?", typeCultureID)
	}

	if err := query.Find(&annonces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []models.ListeAnnonceVente
	for _, a := range annonces {
		result = append(result, toAnnonceDTO(a))
	}

	c.JSON(http.StatusOK, result)
}

type CreateAnnonceVenteInput struct {
	Statut        string `json:"statut" binding:"required"`
	Description   string `json:"description" binding:"required"`
	TypeCultureID string `json:"type_culture_id" binding:"required"`
	ParcelleID    string `json:"parcelle_id" binding:"required"`
	Quantite      string `json:"quantite" binding:"required"`
	PrixKg        string `json:"prix_kg" binding:"required"`
	Photo         string `json:"photo"` // Optional, handle separately if needed
}

// Créer une nouvelle annonce
func CreateAnnonceVente(c *gin.Context) {
	var input CreateAnnonceVenteInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données JSON invalides ou manquantes"})
		return
	}

	// Récupérer l'ID utilisateur depuis le contexte (mis par le middleware)
	userIDFloat, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Convertir l'ID utilisateur en UUID
	userID, err := uuid.Parse(fmt.Sprintf("%.0f", userIDFloat))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}
	typeCultureID, err := uuid.Parse(input.TypeCultureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID type culture invalide"})
		return
	}
	parcelleID, err := uuid.Parse(input.ParcelleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parcelle invalide"})
		return
	}

	// Validate parcelle existence
	_, ok := validateParcelleID(c, input.ParcelleID)
	if !ok {
		return
	}

	quantite, err := strconv.ParseFloat(input.Quantite, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantité invalide"})
		return
	}

	prixKg, err := strconv.ParseFloat(input.PrixKg, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prix au kg invalide"})
		return
	}

	annonce := models.AnnonceVente{
		ID:            uuid.New(),
		Statut:        input.Statut,
		Description:   input.Description,
		UserID:        userID,
		TypeCultureID: typeCultureID,
		ParcelleID:    parcelleID,
		Quantite:      quantite,
		PrixKg:        prixKg,
		Photo:         input.Photo,
	}

	if err := database.DB.Create(&annonce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création : " + err.Error()})
		return
	}

	database.DB.Preload("User").Preload("TypeCulture").Preload("Parcelle").First(&annonce)
	result := toAnnonceDTO(annonce)

	c.JSON(http.StatusCreated, result)
}

// Récupérer une annonce par ID
func GetAnnonceByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var annonce models.AnnonceVente
	if err := database.DB.Preload("User").Preload("TypeCulture").Preload("Parcelle").
		First(&annonce, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce non trouvée"})
		return
	}

	result := toAnnonceDTO(annonce)
	c.JSON(http.StatusOK, result)
}

// Modifier une annonce existante
func UpdateAnnonceVente(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var annonce models.AnnonceVente
	if err := database.DB.First(&annonce, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce non trouvée"})
		return
	}

	// Structure pour recevoir les données JSON
	var input struct {
		Statut        string `json:"statut"`
		Description   string `json:"description"`
		TypeCultureID string `json:"type_culture_id"`
		ParcelleID    string `json:"parcelle_id"`
		Quantite      string `json:"quantite"`
		PrixKg        string `json:"prix_kg"`
		Photo         string `json:"photo"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données JSON invalides ou manquantes"})
		return
	}

	// Log the received payload for debugging
	log.Printf("Update request payload: %+v", input)

	// Mise à jour des champs si fournis
	if input.Statut != "" {
		annonce.Statut = input.Statut
	}
	if input.Description != "" {
		annonce.Description = input.Description
	}
	if input.TypeCultureID != "" {
		typeCultureID, err := uuid.Parse(input.TypeCultureID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID type culture invalide"})
			return
		}
		annonce.TypeCultureID = typeCultureID
	}
	if input.ParcelleID != "" {
		parcelleID, ok := validateParcelleID(c, input.ParcelleID)
		if !ok {
			return
		}
		// Only update if a valid parcelleID was provided
		if parcelleID != uuid.Nil {
			annonce.ParcelleID = parcelleID
		}
	}
	if input.Quantite != "" {
		quantite, err := strconv.ParseFloat(input.Quantite, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Quantité invalide"})
			return
		}
		annonce.Quantite = quantite
	}
	if input.PrixKg != "" {
		prixKg, err := strconv.ParseFloat(input.PrixKg, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Prix au kg invalide"})
			return
		}
		annonce.PrixKg = prixKg
	}
	if input.Photo != "" {
		annonce.Photo = input.Photo
	}

	// Sauvegarde
	if err := database.DB.Save(&annonce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur mise à jour : " + err.Error()})
		return
	}

	database.DB.Preload("User").Preload("TypeCulture").Preload("Parcelle").First(&annonce)
	result := toAnnonceDTO(annonce)
	c.JSON(http.StatusOK, result)
}

// Supprimer une annonce par ID
func DeleteAnnonceVente(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	if err := database.DB.Delete(&models.AnnonceVente{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Annonce supprimée avec succès"})
}
