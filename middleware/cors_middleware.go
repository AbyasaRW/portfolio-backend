package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// buildAllowedOrigins mengembalikan map origin yang diizinkan
func buildAllowedOrigins() map[string]bool {
	allowed := map[string]bool{
		"https://abyasa-portfolio.vercel.app": true,
	}

	// Tambah origin dari env variable (comma-separated) untuk fleksibilitas
	if extra := os.Getenv("ALLOWED_ORIGINS"); extra != "" {
		for _, o := range strings.Split(extra, ",") {
			origin := strings.TrimSpace(o)
			if origin != "" {
				allowed[origin] = true
			}
		}
	}

	return allowed
}

func CORSMiddleware() fiber.Handler {
	allowedOrigins := buildAllowedOrigins()

	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		// Hanya set Access-Control-Allow-Origin jika origin ada di whitelist
		if allowedOrigins[origin] {
			c.Set("Access-Control-Allow-Origin", origin)
		}

		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}

		return c.Next()
	}
}
