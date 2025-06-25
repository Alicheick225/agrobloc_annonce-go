package main

import (
	"github.com/Steph-business/annonce_de_vente/database"
	"github.com/Steph-business/annonce_de_vente/routes"
)

func main() {
	database.InitDB()
	r := routes.SetupRoutes()
	r.Run(":8080")
}
