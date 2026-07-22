package controllers

import (
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

// SetupAdminRoutes serve admin panel HTML langsung dari Go backend.
// adminFS harus di-pass dari main package karena go:embed tidak support path ../
func SetupAdminRoutes(app *fiber.App, adminFS fs.FS) {
	app.Use("/admin", filesystem.New(filesystem.Config{
		Root:   http.FS(adminFS),
		Index:  "index.html",
		Browse: false,
	}))
}
