package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Steph-business/annonce_de_vente/database"
	"github.com/Steph-business/annonce_de_vente/models"
)

// Liste toutes les annonces avec filtres facultatifs
func GetAllAnnonceAchat(c *gin.Context) {
	var achats []models.AnnonceAchat

	// Construire la requête dynamique
	query := database.DB.Preload("TypeCulture").Preload("User")

	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if statut := c.Query("statut"); statut != "" {
		query = query.Where("statut = ?", statut)
	}
	if typeCultureID := c.Query("type_culture_id"); typeCultureID != "" {
		query = query.Where("type_culture_id = ?", typeCultureID)
	}

	// Exécuter la requête
	if err := query.Find(&achats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []models.ListeAnnonceAchat
	for _, a := range achats {
		result = append(result, models.ListeAnnonceAchat{
			ID:                 a.ID.String(),
			Statut:             a.Statut,
			Prix:               a.Prix,
			Description:        a.Description,
			Quantite:           a.Quantite,
			UserNom:            a.User.Nom,
			TypeCultureLibelle: a.TypeCulture.Libelle,
		})
	}

	c.JSON(http.StatusOK, result)
}

// Afficher toutes les annonces d'achat d'un utilisateur
func GetAnnoncesAchatByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")

	// Vérifier que c'est un UUID valide
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var annonces []models.AnnonceAchat
	if err := database.DB.Where("user_id = ?", userID).Find(&annonces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des annonces : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, annonces)
}

type CreateAnnonceAchatInput struct {
	Statut        string  `json:"statut" binding:"required"`
	Prix          float64 `json:"prix_kg" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	TypeCultureID string  `json:"type_culture_id" binding:"required"`
	Quantite      float64 `json:"quantite" binding:"required"`
}

// Créer une nouvelle annonce_achat
func CreateAnnonceAchat(c *gin.Context) {
	var input CreateAnnonceAchatInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur de format JSON: " + err.Error()})
		return
	}

	typeCultureID, err := uuid.Parse(input.TypeCultureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID type culture invalide"})
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

	achats := models.AnnonceAchat{
		ID:            uuid.New(),
		UserID:        userID,
		Statut:        input.Statut,
		Prix:          input.Prix,
		Description:   input.Description,
		TypeCultureID: typeCultureID,
		Quantite:      input.Quantite,
	}

	if err := database.DB.Create(&achats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création : " + err.Error()})
		return
	}

	if err := database.DB.Preload("User").Preload("TypeCulture").First(&achats, "id = ?", achats.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur récupération après création : " + err.Error()})
		return
	}

	result := models.ListeAnnonceAchat{
		ID:                 achats.ID.String(),
		Statut:             achats.Statut,
		Description:        achats.Description,
		Prix:               achats.Prix,
		Quantite:           achats.Quantite,
		UserNom:            achats.User.Nom,
		TypeCultureLibelle: achats.TypeCulture.Libelle,
	}

	c.JSON(http.StatusCreated, result)
}

// Récupération d'une annonce par son Identifiant
func GetAnnonceAchatByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var achats models.AnnonceAchat
	if err := database.DB.First(&achats, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce_achat non trouvée"})
		return
	}

	result := models.ListeAnnonceAchat{
		ID:                 achats.ID.String(),
		Statut:             achats.Statut,
		Description:        achats.Description,
		Prix:               achats.Prix,
		Quantite:           achats.Quantite,
		UserNom:            achats.User.Nom,
		TypeCultureLibelle: achats.TypeCulture.Libelle,
	}

	c.JSON(http.StatusOK, result)
}

type UpdateAnnonceAchatInput struct {
	Statut        string  `json:"statut" binding:"required"`
	Prix          float64 `json:"prix_kg" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	TypeCultureID string  `json:"type_culture_id" binding:"required"`
	Quantite      float64 `json:"quantite" binding:"required"`
}

// Modifier une annonce existante
func UpdateAnnonceAchat(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Récupération de l'annonce existante
	var achats models.AnnonceAchat
	if err := database.DB.First(&achats, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce_achat non trouvée"})
		return
	}

	var input struct {
		Statut        string  `json:"statut"`
		Prix          float64 `json:"prix_kg"`
		Description   string  `json:"description"`
		TypeCultureID string  `json:"type_culture_id"`
		Quantite      float64 `json:"quantite"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur de format JSON: " + err.Error()})
		return
	}

	// Mise à jour des champs si fournis
	if input.Statut != "" {
		achats.Statut = input.Statut
	}
	if input.Prix != 0 {
		achats.Prix = input.Prix
	}
	if input.Description != "" {
		achats.Description = input.Description
	}
	if input.TypeCultureID != "" {
		typeCultureID, err := uuid.Parse(input.TypeCultureID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID type culture invalide"})
			return
		}
		achats.TypeCultureID = typeCultureID
	}
	if input.Quantite != 0 {
		achats.Quantite = input.Quantite
	}

	if err := database.DB.Save(&achats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour : " + err.Error()})
		return
	}

	// Recharger avec les relations
	if err := database.DB.Preload("User").Preload("TypeCulture").First(&achats, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du chargement des relations : " + err.Error()})
		return
	}

	result := models.ListeAnnonceAchat{
		ID:                 achats.ID.String(),
		Statut:             achats.Statut,
		Description:        achats.Description,
		Prix:               achats.Prix,
		Quantite:           achats.Quantite,
		UserNom:            achats.User.Nom,
		TypeCultureLibelle: achats.TypeCulture.Libelle,
	}

	c.JSON(http.StatusOK, result)
}

// Suppression d'une annonce par son ID
func DeleteAnnonceAchat(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'identifant est invalide"})
		return
	}

	if err := database.DB.Delete(&models.AnnonceAchat{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Annonce supprimée avec succès"})
}
