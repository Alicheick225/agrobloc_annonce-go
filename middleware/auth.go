package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims structure qui correspond au payload de votre token Node.js
type Claims struct {
	UserID   float64 `json:"user_id"`
	ProfilID float64 `json:"profil_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware vérifie le token JWT généré par votre API Node.js
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token requis"})
			c.Abort()
			return
		}

		// Format attendu: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Format de token invalide"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		// Utiliser le même secret que votre API Node.js
		jwtSecret := []byte("mon_projet_agrobloc_jwt")

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusForbidden, gin.H{"message": "Token invalide"})
			c.Abort()
			return
		}

		// Stocker les informations de l'utilisateur dans le contexte Gin
		c.Set("user_id", int(claims.UserID))
		c.Set("profil_id", int(claims.ProfilID))

		c.Next()
	}
}
