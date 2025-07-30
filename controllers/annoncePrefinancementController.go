package controllers

import (
	"net/http"
	"strconv"

	"github.com/Steph-business/annonce_de_vente/database"
	"github.com/Steph-business/annonce_de_vente/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 🔹 Lister toutes les annonces de préfinancement avec relations
func GetAllAnnoncePref(c *gin.Context) {
	var annonces []models.AnnoncePrefinancement
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

	var result []models.LiteAnnoncePrefinancement
	for _, a := range annonces {
		result = append(result, models.LiteAnnoncePrefinancement{
			ID:                 a.ID.String(),
			Statut:             a.Statut,
			Description:        a.Description,
			MontantPref:        a.MontantPrefinancement,
			PrixKgPref:         a.Prix,
			Quantite:           a.Quantite,
			UserNom:            a.User.Nom,
			TypeCultureLibelle: a.TypeCulture.Libelle,
			ParcelleAdresse:    a.Parcelle.Adresse,
			ParcelleSuf:        a.Parcelle.Surface,
		})
	}

	c.JSON(http.StatusOK, result)
}

// Créer une nouvelle annonce de préfinancement
func CreateAnnoncePref(c *gin.Context) {
	var annonce models.AnnoncePrefinancement

	annonce.Statut = c.PostForm("statut")
	annonce.Description = c.PostForm("description")

	userID, err := uuid.Parse(c.PostForm("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}
	typeCultureID, err := uuid.Parse(c.PostForm("type_culture_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID type culture invalide"})
		return
	}
	parcelleID, err := uuid.Parse(c.PostForm("parcelle_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parcelle invalide"})
		return
	}
	annonce.UserID = userID
	annonce.TypeCultureID = typeCultureID
	annonce.ParcelleID = parcelleID

	var tc models.TypeCulture
	if err := database.DB.First(&tc, "id = ?", typeCultureID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type de culture introuvable"})
		return
	}
	var parcelle models.Parcelle
	if err := database.DB.First(&parcelle, "id = ?", parcelleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parcelle introuvable"})
		return
	}
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Utilisateur introuvable"})
		return
	}

	qStr := c.PostForm("quantite")

	quantite, err := strconv.ParseFloat(qStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantité invalide"})
		return
	}
	prixKgStr := c.PostForm("prix")
	prixKg, err := strconv.ParseFloat(prixKgStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prix par kg invalide"})
		return
	}
	if prixKg <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prix par kg doit être supérieur à zéro"})
		return
	}
	annonce.Prix = prixKg
	annonce.Quantite = quantite

	annonce.MontantPrefinancement = prixKg * quantite
	annonce.ID = uuid.New()

	if err := database.DB.Create(&annonce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création : " + err.Error()})
		return
	}
	database.DB.Preload("User").Preload("Parcelle").Preload("TypeCulture").First(&annonce)

	result := models.LiteAnnoncePrefinancement{
		ID:                 annonce.ID.String(),
		Statut:             annonce.Statut,
		Description:        annonce.Description,
		MontantPref:        annonce.MontantPrefinancement,
		PrixKgPref:         annonce.Prix,
		Quantite:           annonce.Quantite,
		UserNom:            user.Nom,
		TypeCultureLibelle: tc.Libelle,
		ParcelleAdresse:    parcelle.Adresse,
		ParcelleSuf:        parcelle.Surface,
	}

	c.JSON(http.StatusCreated, result)
}

// 🔹 Récupérer les annonces d’un utilisateur
func GetPrefinancementsByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	var annonces []models.AnnoncePrefinancement
	if err := database.DB.Where("user_id = ?", userID).Find(&annonces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, annonces)
}

// 🔹 Récupérer une annonce par ID
func GetAnnoncePrefByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var annonce models.AnnoncePrefinancement
	if err := database.DB.Preload("User").Preload("Parcelle").Preload("TypeCulture").First(&annonce, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce non trouvée"})
		return
	}

	result := models.LiteAnnoncePrefinancement{
		ID:                 annonce.ID.String(),
		Statut:             annonce.Statut,
		Description:        annonce.Description,
		MontantPref:        annonce.MontantPrefinancement,
		PrixKgPref:         annonce.Prix,
		Quantite:           annonce.Quantite,
		UserNom:            annonce.User.Nom,
		TypeCultureLibelle: annonce.TypeCulture.Libelle,
		ParcelleAdresse:    annonce.Parcelle.Adresse,
		ParcelleSuf:        annonce.Parcelle.Surface,
	}

	c.JSON(http.StatusOK, result)
}

// 🔹 Modifier une annonce existante
func UpdateAnnoncePref(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var annonce models.AnnoncePrefinancement
	if err := database.DB.First(&annonce, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Annonce non trouvée"})
		return
	}

	var input models.AnnoncePrefinancement
	body, err := c.GetRawData()
	if err != nil || len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide : corps vide"})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide : " + err.Error()})
		return
	}

	var tc models.TypeCulture
	if err := database.DB.First(&tc, "id = ?", input.TypeCultureID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type de culture introuvable"})
		return
	}
	if input.ParcelleID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ParcelleID manquant"})
		return
	}
	var parcelle models.Parcelle
	if err := database.DB.First(&parcelle, "id = ?", input.ParcelleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parcelle introuvable"})
		return
	}

	annonce.Statut = input.Statut
	annonce.Description = input.Description
	annonce.Quantite = input.Quantite
	annonce.Prix = input.Prix
	annonce.TypeCultureID = input.TypeCultureID
	annonce.ParcelleID = input.ParcelleID
	annonce.UserID = input.UserID
	annonce.MontantPrefinancement = input.Prix * input.Quantite

	if err := database.DB.Save(&annonce).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur mise à jour : " + err.Error()})
		return
	}

	database.DB.Preload("User").Preload("Parcelle").Preload("TypeCulture").First(&annonce)

	result := models.LiteAnnoncePrefinancement{
		ID:                 annonce.ID.String(),
		Statut:             annonce.Statut,
		Description:        annonce.Description,
		MontantPref:        annonce.MontantPrefinancement,
		PrixKgPref:         annonce.Prix,
		Quantite:           annonce.Quantite,
		UserNom:            annonce.User.Nom,
		TypeCultureLibelle: annonce.TypeCulture.Libelle,
		ParcelleAdresse:    annonce.Parcelle.Adresse,
		ParcelleSuf:        annonce.Parcelle.Surface,
	}

	c.JSON(http.StatusOK, result)
}

// 🔹 Supprimer une annonce
func DeleteAnnoncePref(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	if err := database.DB.Delete(&models.AnnoncePrefinancement{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur suppression : " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Annonce supprimée avec succès"})
}
