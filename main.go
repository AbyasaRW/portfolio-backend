package main

import (
	"embed"
	"io/fs"
	"log"
	"os"

	"portfolio-backend/config"
	"portfolio-backend/controllers"
	"portfolio-backend/middleware"
	"portfolio-backend/models"
	"portfolio-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

//go:embed admin/*
var adminFiles embed.FS

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

	// Serve embedded admin panel di /admin
	adminSubFS, _ := fs.Sub(adminFiles, "admin")
	controllers.SetupAdminRoutes(app, adminSubFS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}

