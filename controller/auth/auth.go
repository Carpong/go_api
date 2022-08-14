package auth

import (
	"fmt"
	"go/rest-api/database"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

type Register struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(c *gin.Context) {
	var json Register
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var CheckUser database.User
	database.Db.Where("username = ?", json.Username).First(&CheckUser)
	if CheckUser.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Exists"})
		return
	}

	HashPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	user := database.User{Username: json.Username, Password: string(HashPassword)}
	database.Db.Create(&user)
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Create Success", "user_id": user.ID})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Create Failed"})
	}
}

func LoginUser(c *gin.Context) {
	var json Register
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var CheckUser database.User
	database.Db.Where("username = ?", json.Username).First(&CheckUser)
	if CheckUser.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Exists"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(CheckUser.Password), []byte(json.Password))
	if err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"UserID": CheckUser.ID,
			"exp":    time.Now().Add(time.Minute * 15).Unix(),
		})
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Login Success", "token": tokenString})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Login Failed"})
	}
}
