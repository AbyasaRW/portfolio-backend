package controllers

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// DefaultSiteURL dipakai bila SITE_URL tidak diisi (alamat `npm run dev` frontend).
const DefaultSiteURL = "http://localhost:3000"

// GetSiteInfo memberi tahu Admin Panel alamat website portfolio.
//
// Kenapa perlu: gambar project disimpan sebagai path relatif (mis. /Projects/x.png) yang
// letaknya di folder public/ FRONTEND, bukan di backend. Admin Panel memakai alamat ini
// untuk menampilkan thumbnail dan tautan "Lihat website" dengan benar.
// Endpoint ini publik dan hanya berisi satu URL, tidak ada data rahasia.
func GetSiteInfo(c *fiber.Ctx) error {
	siteURL := strings.TrimRight(strings.TrimSpace(os.Getenv("SITE_URL")), "/")
	if siteURL == "" {
		siteURL = DefaultSiteURL
	}
	return c.JSON(fiber.Map{"site_url": siteURL})
}
