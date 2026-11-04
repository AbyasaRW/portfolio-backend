package controllers

import (
	"portfolio-backend/config"
	"portfolio-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Resource adalah syarat sebuah model bisa memakai CRUD generik di bawah:
// bisa dirapikan (Normalize), divalidasi (Problem), dan identitasnya bisa dilindungi (Base).
type Resource[T any] interface {
	*T
	Normalize()
	Problem() string
	Base() *gorm.Model
}

// Daftar yang tampil di website. Urutan: sort_order, lalu id.
const contentOrder = "sort_order ASC, id ASC"

// listAll mengembalikan semua data. Daftar kosong dikirim sebagai [] (bukan null).
func listAll[T any, P Resource[T]](c *fiber.Ctx) error {
	items := []T{}
	config.DB.Order(contentOrder).Find(&items)
	return c.JSON(items)
}

func createOne[T any, P Resource[T]](c *fiber.Ctx) error {
	var item T
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	p := P(&item)
	p.Normalize()
	if msg := p.Problem(); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}
	p.Base().ID = 0 // ID selalu dibuat database, bukan dari body

	config.DB.Create(&item)
	return c.Status(fiber.StatusCreated).JSON(item)
}

func updateOne[T any, P Resource[T]](c *fiber.Ctx) error {
	var item T
	if err := config.DB.First(&item, c.Params("id")).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}

	// Simpan identitas asli supaya body request tidak bisa menimpa ID / created_at.
	p := P(&item)
	original := *p.Base()

	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	*p.Base() = original
	p.Normalize()
	if msg := p.Problem(); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	config.DB.Save(&item)
	return c.JSON(item)
}

func deleteOne[T any, P Resource[T]](c *fiber.Ctx) error {
	var item T
	if err := config.DB.First(&item, c.Params("id")).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data tidak ditemukan"})
	}
	config.DB.Delete(&item)
	return c.JSON(fiber.Map{"message": "Data berhasil dihapus"})
}

// ---- Handler siap pakai untuk routes ----

func GetExperiences(c *fiber.Ctx) error   { return listAll[models.Experience](c) }
func CreateExperience(c *fiber.Ctx) error { return createOne[models.Experience](c) }
func UpdateExperience(c *fiber.Ctx) error { return updateOne[models.Experience](c) }
func DeleteExperience(c *fiber.Ctx) error { return deleteOne[models.Experience](c) }

func GetTools(c *fiber.Ctx) error   { return listAll[models.Tool](c) }
func CreateTool(c *fiber.Ctx) error { return createOne[models.Tool](c) }
func UpdateTool(c *fiber.Ctx) error { return updateOne[models.Tool](c) }
func DeleteTool(c *fiber.Ctx) error { return deleteOne[models.Tool](c) }