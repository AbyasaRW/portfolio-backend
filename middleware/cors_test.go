// middleware/cors_test.go
package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"portfolio-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func setupCORSApp() *fiber.App {
	app := fiber.New()
	app.Use(middleware.CORSMiddleware())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	return app
}

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	app := setupCORSApp()

	// Vercel origin harus diizinkan (hardcoded di middleware)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://abyasa-portfolio.vercel.app")

	res, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "https://abyasa-portfolio.vercel.app", res.Header.Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_BlockedOrigin(t *testing.T) {
	app := setupCORSApp()

	// Origin yang tidak di whitelist harus diblokir (tidak ada header ACAO)
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://evil-site.com")

	res, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Empty(t, res.Header.Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_OptionsPreFlight(t *testing.T) {
	app := setupCORSApp()

	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "https://abyasa-portfolio.vercel.app")

	res, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)
	require.Equal(t, "https://abyasa-portfolio.vercel.app", res.Header.Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_EnvOrigin(t *testing.T) {
	// Set ALLOWED_ORIGINS env var untuk tambah origin ekstra
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	// Buat app baru supaya buildAllowedOrigins() dipanggil ulang
	app := fiber.New()
	app.Use(middleware.CORSMiddleware())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	res, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "http://localhost:3000", res.Header.Get("Access-Control-Allow-Origin"))
}

