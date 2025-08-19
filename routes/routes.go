package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Steph-business/annonce_de_vente/controllers"
	"github.com/Steph-business/annonce_de_vente/middleware"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	
	// Routes publiques (lecture)
	r.GET("/annonces_vente", controllers.GetAllAnnonceVente)
	r.GET("/annonces_vente/:id", controllers.GetAnnonceByID)
	
	r.GET("/annonces_achat", controllers.GetAllAnnonceAchat)
	r.GET("/annonces_achat/:id", controllers.GetAnnonceAchatByID)
	
	r.GET("/annonces_pref", controllers.GetAllAnnoncePref)
	r.GET("/annonces_pref/:id", controllers.GetAnnoncePrefByID)

	// Routes protégées (nécessitent authentification)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Annonces de vente
		protected.POST("/annonces_vente", controllers.CreateAnnonceVente)
		protected.PUT("/annonces_vente/:id", controllers.UpdateAnnonceVente)
		protected.DELETE("/annonces_vente/:id", controllers.DeleteAnnonceVente)

		// Annonces d'achat
		protected.POST("/annonces_achat", controllers.CreateAnnonceAchat)
		protected.PUT("/annonces_achat/:id", controllers.UpdateAnnonceAchat)
		protected.DELETE("/annonces_achat/:id", controllers.DeleteAnnonceAchat)

		// Annonces de préfinancement
		protected.POST("/annonces_pref", controllers.CreateAnnoncePref)
		protected.PUT("/annonces_pref/:id", controllers.UpdateAnnoncePref)
		protected.DELETE("/annonces_pref/:id", controllers.DeleteAnnoncePref)
	}

	return r
}
