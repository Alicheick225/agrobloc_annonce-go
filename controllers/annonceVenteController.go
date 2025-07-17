package controllers

import (
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

// validateParcelleID vérifie que le parcelleID est non vide, valide et existe en base
func validateParcelleID(c *gin.Context, parcelleIDStr string) (uuid.UUID, bool) {
	if parcelleIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ParcelleID manquant"})
		log.Println("Validation ParcelleID échouée : champ vide")
		return uuid.Nil, false
	}
	log.Printf("Validation ParcelleID : valeur reçue '%s'\n", parcelleIDStr)
	parcelleID, err := uuid.Parse(parcelleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ParcelleID invalide"})
		log.Printf("Validation ParcelleID échouée : parsing UUID invalide pour %s\n", parcelleIDStr)
		return uuid.Nil, false
	}
	var parcelle models.Parcelle
	if err := database.DB.First(&parcelle, "id = ?", parcelleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parcelle introuvable"})
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

// Créer une nouvelle annonce
func CreateAnnonceVente(c *gin.Context) {
	var annonce models.AnnonceVente

	// Récupère les champs du formulaire
	annonce.Statut = c.PostForm("statut")
	annonce.Description = c.PostForm("description")
	annonce.UserID, _ = uuid.Parse(c.PostForm("user_id"))
	annonce.TypeCultureID, _ = uuid.Parse(c.PostForm("type_culture_id"))
	parcelleIDStr := c.PostForm("parcelle_id")
	parcelleID, ok := validateParcelleID(c, parcelleIDStr)
	if !ok {
		return
	}
	annonce.ParcelleID = parcelleID

	quantiteStr := c.PostForm("quantite")
	quantite, err := strconv.ParseFloat(quantiteStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantité invalide"})
		return
	}
	annonce.Quantite = quantite

	prixKgStr := c.PostForm("prix_kg")
	prixKg, err := strconv.ParseFloat(prixKgStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prix par kg invalide"})
		return
	}
	annonce.PrixKg = prixKg

	// Gérer l'image
	file, err := c.FormFile("photo")
	if err == nil {
		// Génère un nom de fichier unique
		filename := file.Filename
		filepath := "static/" + filename

		// Sauvegarde le fichier sur le disque
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de sauvegarde de l'image : " + err.Error()})
			return
		}

		// Enregistre le chemin d’accès
		annonce.Photo = "" + filepath
	}

	// ID unique
	annonce.ID = uuid.New()

	// Enregistrement
	if err := database.DB.Create(&annonce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création : " + err.Error()})
		return
	}

	// Recharge avec relations
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

	// Mise à jour des champs (form-data)
	annonce.Statut = c.PostForm("statut")
	annonce.Description = c.PostForm("description")
	annonce.TypeCultureID, _ = uuid.Parse(c.PostForm("type_culture_id"))
	parcelleIDStr := c.PostForm("parcelle_id")
	parcelleID, ok := validateParcelleID(c, parcelleIDStr)
	if !ok {
		return
	}
	annonce.ParcelleID = parcelleID

	annonce.Quantite, _ = strconv.ParseFloat(c.PostForm("quantite"), 64)
	annonce.PrixKg, _ = strconv.ParseFloat(c.PostForm("prix_kg"), 64)

	// Vérifie si une nouvelle image est envoyée
	file, err := c.FormFile("photo")
	if err == nil {
		// Génère un nom de fichier
		filename := file.Filename
		filepath := "static/" + filename

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur upload image : " + err.Error()})
			return
		}

		// Remplace le chemin image
		annonce.Photo = "" + filepath
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
