package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"github.com/Steph-business/annonce_de_vente/database"
	"github.com/Steph-business/annonce_de_vente/models"
)

//Liste toutes les annonces avec filtres facultatifs
func GetAllAnnonceAchat(c *gin.Context) {
	var achats []models.AnnonceAchat

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
	if err := query.Find(&achats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, achats)
}

// 🔹 Afficher toutes les annonces d'achat d'un utilisateur
func GetAnnoncesAchatByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")

	// Vérifier que c’est un UUID valide
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



//Créer une nouvelle annonce_achat
func CreateAnnonceAchat(c *gin.Context) {
	var achats models.AnnonceAchat
	if err := c.ShouldBindJSON(&achats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Soumission invalide : " + err.Error()})
		return
	}

	// Assigner un ID unique si nécessaire
	achats.ID = uuid.New()

	if err := database.DB.Create(&achats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création : " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, achats)
}

//Récupération d'une annonce par son Identifiant
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

	c.JSON(http.StatusOK, achats)
}


//  Modifier une annonce existante
func UpdateAnnonceAchat(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	//Récupération de l’annonce existante
	var achats models.AnnonceAchat
	if err := database.DB.First(&achats, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce_chat non trouvée"})
		return
	}

	//Mise à jour avec les données entrantes
	var updatedAnnonce models.AnnonceAchat
	if err := c.ShouldBindJSON(&updatedAnnonce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide : " + err.Error()})
		return
	}

	//Mise à jour individuel
	achats.Statut = updatedAnnonce.Statut
	achats.Quantite = updatedAnnonce.Quantite

	if err := database.DB.Save(&achats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, achats)
}



//Suppression d'une annonce par son ID
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

