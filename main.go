package main

import (
	"fmt"
	AuthController "go/rest-api/controller/auth"
	UserController "go/rest-api/controller/user"
	"go/rest-api/database"
	"go/rest-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	database.InitDB()
	r := gin.Default()
	r.POST("/register", AuthController.RegisterUser)
	r.POST("/login", AuthController.LoginUser)
	authorized := r.Group("/users", middleware.Auth())
	authorized.GET("/readall", UserController.ReadAll)
	authorized.GET("/profile", UserController.Profile)
	authorized.GET("/fileall", UserController.Listfile)
	authorized.POST("/upload", UserController.Upload)
	authorized.PATCH("/update/:id", UserController.UpdateFile)
	authorized.DELETE("/delete/:id", UserController.DeleteFile)
	r.Run()
}
