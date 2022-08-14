package user

import (
	"go/rest-api/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReadAll(c *gin.Context) {
	var users []database.User
	database.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "User": users})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("UserID").(float64)
	var user database.User
	database.Db.First(&user, userId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": user})
}
