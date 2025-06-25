package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"github.com/Steph-business/annonce_de_vente/database"
	"github.com/Steph-business/annonce_de_vente/models"
)

//Liste toutes les annonces
// func GetAllAnnonceVente(c *gin.Context) {
// 	var annonces []models.AnnonceVente
// 	if err := database.DB.Find(&annonces).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, annonces)
// }

//Liste toutes les annonces avec filtres facultatifs
func GetAllAnnonceVente(c *gin.Context) {
	var annonces []models.AnnonceVente

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
	if err := query.Find(&annonces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, annonces)
}


//Créer une nouvelle annonce
func CreateAnnonceVente(c *gin.Context) {
	var annonce models.AnnonceVente
	if err := c.ShouldBindJSON(&annonce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Soumission invalide : " + err.Error()})
		return
	}

	// Assigner un ID unique si nécessaire
	annonce.ID = uuid.New()

	if err := database.DB.Create(&annonce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création : " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, annonce)
}

//Récupération d'une annonce par son Identifiant
func GetAnnonceByID(c *gin.Context) {
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

	c.JSON(http.StatusOK, annonce)
}


//  Modifier une annonce existante
func UpdateAnnonceVente(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	//Récupération de l’annonce existante
	var annonce models.AnnonceVente
	if err := database.DB.First(&annonce, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce non trouvée"})
		return
	}

	//Mise à jour avec les données entrantes
	var updatedAnnonce models.AnnonceVente
	if err := c.ShouldBindJSON(&updatedAnnonce); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide : " + err.Error()})
		return
	}

	//Mise à jour individuel
	annonce.Photo = updatedAnnonce.Photo
	annonce.Statut = updatedAnnonce.Statut
	annonce.Quantite = updatedAnnonce.Quantite
	annonce.PrixKg = updatedAnnonce.PrixKg
	annonce.ParcelleID = updatedAnnonce.ParcelleID
	annonce.TypeCultureID = updatedAnnonce.TypeCultureID

	if err := database.DB.Save(&annonce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, annonce)
}



//Suppression d'une annonce par son ID
func DeleteAnnonceVente(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'identifant est invalide"})
		return
	}

	if err := database.DB.Delete(&models.AnnonceVente{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Annonce supprimée avec succès"})
}

