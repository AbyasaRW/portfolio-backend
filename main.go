package main

import (
	"log"
	"os"

	"portfolio-backend/config"
	"portfolio-backend/models"
	"portfolio-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file tidak ditemukan, pakai environment variable sistem")
	}

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Project{}, &models.User{})

	router := gin.Default()
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}
