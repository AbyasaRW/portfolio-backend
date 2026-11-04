package routes

import (
	"os"
	"portfolio-backend/controllers"
	"portfolio-backend/middleware"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth routes — publik, siapa aja boleh akses
	if strings.EqualFold(os.Getenv("ENABLE_REGISTRATION"), "true") {
		api.Post("/register", controllers.Register)
	} else {
		api.Post("/register", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNotFound)
		})
	}
	api.Post("/login", controllers.Login)

	// Info website (publik): dipakai Admin Panel untuk thumbnail & tautan "Lihat website".
	api.Get("/site", controllers.GetSiteInfo)

	// Project routes — GET tetap publik (biar portofolio bisa ditampilin ke visitor)
	api.Get("/projects", controllers.GetProjects)
	api.Get("/projects/:id", controllers.GetProjectByID)

	// Pengalaman & Tools — GET publik (dibaca website), perubahan hanya admin.
	api.Get("/experiences", controllers.GetExperiences)
	api.Get("/tools", controllers.GetTools)

	// POST/PUT/DELETE diproteksi — cuma admin (yang punya token) yang boleh
	protected := api.Group("/", middleware.AuthRequired())
	protected.Post("/projects", controllers.CreateProject)
	protected.Put("/projects/:id", controllers.UpdateProject)
	protected.Delete("/projects/:id", controllers.DeleteProject)

	protected.Post("/experiences", controllers.CreateExperience)
	protected.Put("/experiences/:id", controllers.UpdateExperience)
	protected.Delete("/experiences/:id", controllers.DeleteExperience)

	protected.Post("/tools", controllers.CreateTool)
	protected.Put("/tools/:id", controllers.UpdateTool)
	protected.Delete("/tools/:id", controllers.DeleteTool)
}
