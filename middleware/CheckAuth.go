package middleware

import (
	"fmt"
	"go/rest-api/database"
	"go/rest-api/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
		header := c.Request.Header.Get("Authorization")
		tokenString := strings.ReplaceAll(header, "Bearer ", "")

		var CheckBlakclist models.JwtBlacklist
		database.DB.Where("Token = ?", tokenString).First(&CheckBlakclist)

		if CheckBlakclist.ID > 0 {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "error", "message": "Token is expired"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return hmacSampleSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("UserID", claims["UserID"])
		} else {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		}
	}
}
