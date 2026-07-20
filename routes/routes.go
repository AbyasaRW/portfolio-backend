package routes

import (
	"portfolio-backend/controllers"
	"portfolio-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Auth routes — publik, siapa aja boleh akses
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// Project routes — GET tetap publik (biar portofolio bisa ditampilin ke visitor)
		api.GET("/projects", controllers.GetProjects)
		api.GET("/projects/:id", controllers.GetProjectByID)

		// POST/PUT/DELETE diproteksi — cuma admin (yang punya token) yang boleh
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			protected.POST("/projects", controllers.CreateProject)
			protected.PUT("/projects/:id", controllers.UpdateProject)
			protected.DELETE("/projects/:id", controllers.DeleteProject)
		}
	}
}
