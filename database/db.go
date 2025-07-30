package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	// Charger les variables d'environnement depuis le fichier .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Erreur de chargement du fichier .env")
	}

	// Récupération des valeurs d’environnement
	DB_HOST := os.Getenv("host")
	DB_PORT := os.Getenv("port")
	DB_USER := os.Getenv("user")
	DB_PASSWORD := os.Getenv("password")
	DB_NAME := os.Getenv("dbname")

	// Chaîne de connexion PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	// Connexion avec désactivation du protocole préparé
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, //  IMPORTANT : désactive les requêtes préparées
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Erreur de connexion à la base de données :", err)
	}

	fmt.Println("Connexion à la base de données réussie")
}
