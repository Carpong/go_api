package user

import (
	"fmt"
	"go/rest-api/database"
	"go/rest-api/models"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateFileName struct {
	FileName string `json:"filename" binding:"required"`
}

func Upload(c *gin.Context) {
	userId := c.MustGet("UserID").(float64)
	userStr := fmt.Sprintf("%v", userId)
	file, _ := c.FormFile("file")
	GenName := file.Filename
	Filetype := strings.Split(file.Filename, ".")

	if Filetype[1] == "jpg" || Filetype[1] == "jpeg" || Filetype[1] == "png" || Filetype[1] == "gif" || Filetype[1] == "tiff" {
		var CheckFile models.UploadFile
		database.DB.Where("user_id = ?", userId).First(&CheckFile)
		if CheckFile.ID > 0 {
			now := time.Now().Format("20220815121905")
			GenName = fmt.Sprintf("%v", now) + "_" + fmt.Sprintf("%v", rand.Intn(100)) + "_" + file.Filename
		}
		fileJson := models.UploadFile{UserId: userStr, FileName: file.Filename, FileNameGen: GenName}
		database.DB.Create(&fileJson)
		c.SaveUploadedFile(file, "public/images/"+GenName)
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Uploaded file success"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Support filetype image only"})
}

func UpdateFile(c *gin.Context) {
	userId := c.MustGet("UserID").(float64)
	userStr := fmt.Sprintf("%v", userId)
	var FileName models.UploadFile
	if err := database.DB.Where("id = ?", c.Param("id")).First(&FileName).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Record not found"})
		return
	}
	var input UpdateFileName
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	fileJson := models.UploadFile{UserId: userStr, FileName: input.FileName, FileNameGen: FileName.FileNameGen}
	database.DB.Model(&FileName).Updates(fileJson)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Update filename success", "data": input.FileName})
}

func DeleteFile(c *gin.Context) {
	var File models.UploadFile
	if err := database.DB.Where("id = ?", c.Param("id")).First(&File).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Record not found"})
		return
	}
	database.DB.Delete(&File)
	e := os.Remove("public/images/" + File.FileNameGen)
	if e != nil {
		log.Fatal(e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Delete file success"})
}

func Listfile(c *gin.Context) {
	var AllFile []models.UploadFile
	userId := c.MustGet("UserID").(float64)
	if err := database.DB.Where("user_id = ?", userId).Find(&AllFile).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Record not found!"})
		return
	}
	if len(AllFile) <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "No record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Read record success", "fileName": AllFile})
}
