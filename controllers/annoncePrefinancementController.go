package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Steph-business/annonce_de_vente/database"
	"github.com/Steph-business/annonce_de_vente/models"
)

// Liste toutes les annonces avec filtres facultatifs
func GetAllAnnoncePref(c *gin.Context) {
	var prefinancements []models.AnnoncePrefinancement

	// Récupérer les filtres
	userID := c.Query("user_id")
	statut := c.Query("statut")
	typeCultureID := c.Query("type_culture_id")

	// Construire la requête dynamique
	query := database.DB

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if statut != "" {
		query = query.Where("statut = ?", statut)
	}
	if typeCultureID != "" {
		query = query.Where("type_culture_id = ?", typeCultureID)
	}

	// Exécuter la requête
	if err := query.Find(&prefinancements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prefinancements)
}

// Créer une nouvelle annonce_achat
func CreateAnnoncePref(c *gin.Context) {
	var prefinancements models.AnnoncePrefinancement
	if err := c.ShouldBindJSON(&prefinancements); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Soumission invalide : " + err.Error()})
		return
	}

	// Assigner un ID unique si nécessaire
	prefinancements.ID = uuid.New()

	if err := database.DB.Create(&prefinancements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création : " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, prefinancements)
}

// Récupération d'une annonce par son Identifiant
func GetAnnoncePrefByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var prefinancements models.AnnoncePrefinancement
	if err := database.DB.First(&prefinancements, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce_achat non trouvée"})
		return
	}

	c.JSON(http.StatusOK, prefinancements)
}

// Modifier une annonce existante
func UpdateAnnoncePref(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	//Récupération de l’annonce existante
	var prefinancements models.AnnoncePrefinancement
	if err := database.DB.First(&prefinancements, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce_chat non trouvée"})
		return
	}

	//Mise à jour avec les données entrantes
	var updatedAnnonce models.AnnoncePrefinancement
	if err := c.ShouldBindJSON(&updatedAnnonce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide : " + err.Error()})
		return
	}

	//Mise à jour individuel
	prefinancements.Statut = updatedAnnonce.Statut
	prefinancements.Montant = updatedAnnonce.Montant
	prefinancements.Description = updatedAnnonce.Description

	if err := database.DB.Save(&prefinancements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, prefinancements)
}

// Suppression d'une annonce par son ID
func DeleteAnnoncePref(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'identifant est invalide"})
		return
	}

	if err := database.DB.Delete(&models.AnnoncePrefinancement{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Annonce supprimée avec succès"})
}
