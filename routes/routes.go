package routes


import (
	"github.com/gin-gonic/gin"
	"github.com/Steph-business/annonce_de_vente/controllers"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/annonces_vente", controllers.GetAllAnnonceVente)
	r.POST("/annonces_vente", controllers.CreateAnnonceVente)
	r.GET("/annonces_vente/:id", controllers.GetAnnonceByID)
	r.PUT("/annonces_vente/:id", controllers.UpdateAnnonceVente)
	r.DELETE("/annonces_vente/:id", controllers.DeleteAnnonceVente)


	r.GET("/annonces_achat", controllers.GetAllAnnonceAchat)
	r.POST("/annonces_achat", controllers.CreateAnnonceAchat)
	r.GET("/annonces_achat/:id", controllers.GetAnnonceAchatByID)
	r.PUT("/annonces_achat/:id", controllers.UpdateAnnonceAchat)
	r.DELETE("/annonces_achat/:id", controllers.DeleteAnnonceAchat)


	r.GET("/annonces_pref", controllers.GetAllAnnoncePref)
	r.POST("/annonces_pref", controllers.CreateAnnoncePref)
	r.GET("/annonces_pref/:id", controllers.GetAnnoncePrefByID)
	r.PUT("/annonces_pref/:id", controllers.UpdateAnnoncePref)
	r.DELETE("/annonces_pref/:id", controllers.DeleteAnnoncePref)

	return r
}
