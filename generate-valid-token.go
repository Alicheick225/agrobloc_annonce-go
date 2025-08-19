package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken() {
	// Données de test
	claims := jwt.MapClaims{
		"user_id":   123.0,
		"profil_id": 456.0,
		"exp":       time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 jours
	}

	// Secret (doit correspondre exactement à votre middleware)
	secret := "mon_projet_agrobloc_jwt"

	// Créer le token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Printf("Erreur lors de la génération du token: %v\n", err)
		return
	}

	fmt.Println("=== TOKEN JWT VALIDE (Go) ===")
	fmt.Println(tokenString)
	fmt.Println("==============================")

	// Pour tester le décodage
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Printf("Erreur lors de la validation du token: %v\n", err)
		return
	}

	fmt.Println("Token validé avec succès!")
	fmt.Printf("Claims: %+v\n", parsedToken.Claims)
}

func main() {
	generateToken()
}
