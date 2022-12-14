package database

import (
	"go/rest-api/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect() {
	dsn := os.Getenv("MYSQL_DNS")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.UploadFile{})
	DB.AutoMigrate(&models.JwtBlacklist{})
}
