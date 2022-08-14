package user

import (
	"fmt"
	"go/rest-api/database"
	"go/rest-api/models"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateFileName struct {
	FileName string `json:"filename" binding:"required"`
}

func ReadAll(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "User": users})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("UserID").(float64)
	var user models.User
	database.DB.First(&user, userId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": user})
}

func Upload(c *gin.Context) {
	userId := c.MustGet("UserID").(float64)
	userStr := fmt.Sprintf("%v", userId)
	file, _ := c.FormFile("file")
	GenName := file.Filename

	var CheckFile models.UploadFile
	database.DB.Where("user_id = ?", userId).First(&CheckFile)
	if CheckFile.ID > 0 {
		now := time.Now().Format("20220815121905")
		GenName = fmt.Sprintf("%v", now) + "_" + fmt.Sprintf("%v", rand.Intn(100)) + "_" + file.Filename
	}
	fileJson := models.UploadFile{UserId: userStr, FileName: file.Filename, FileNameGen: GenName}
	database.DB.Create(&fileJson)
	c.SaveUploadedFile(file, "public/images/"+GenName)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Uploaded File Success"})

}

func Listfile(c *gin.Context) {
	var AllFile []models.UploadFile
	userId := c.MustGet("UserID").(float64)
	database.DB.Where("user_id = ?", userId).Find(&AllFile)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User File Success", "fileName": AllFile})
}

func DeleteFile(c *gin.Context) {
	var File models.UploadFile
	if err := database.DB.Where("id = ?", c.Param("id")).First(&File).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	database.DB.Delete(&File)
	e := os.Remove("public/images/" + File.FileNameGen)
	if e != nil {
		log.Fatal(e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func UpdateFile(c *gin.Context) {
	userId := c.MustGet("UserID").(float64)
	userStr := fmt.Sprintf("%v", userId)
	var FileName models.UploadFile
	if err := database.DB.Where("id = ?", c.Param("id")).First(&FileName).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	var input UpdateFileName
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileJson := models.UploadFile{UserId: userStr, FileName: input.FileName, FileNameGen: FileName.FileNameGen}
	database.DB.Model(&FileName).Updates(fileJson)
	c.JSON(http.StatusOK, gin.H{"data": input.FileName})
}
