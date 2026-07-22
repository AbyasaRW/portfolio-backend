package main

import (
	"log"
	"os"

	"portfolio-backend/config"
	"portfolio-backend/middleware"
	"portfolio-backend/models"
	"portfolio-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file tidak ditemukan, pakai environment variable sistem")
	}

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Project{}, &models.User{})
	SeedDatabase()

	app := fiber.New()
	app.Use(middleware.CORSMiddleware())
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
