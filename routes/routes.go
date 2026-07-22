package routes

import (
	"portfolio-backend/controllers"
	"portfolio-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth routes — publik, siapa aja boleh akses
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// Project routes — GET tetap publik (biar portofolio bisa ditampilin ke visitor)
	api.Get("/projects", controllers.GetProjects)
	api.Get("/projects/:id", controllers.GetProjectByID)

	// POST/PUT/DELETE diproteksi — cuma admin (yang punya token) yang boleh
	protected := api.Group("/", middleware.AuthRequired())
	protected.Post("/projects", controllers.CreateProject)
	protected.Put("/projects/:id", controllers.UpdateProject)
	protected.Delete("/projects/:id", controllers.DeleteProject)
}
