package controllers

import (
	"portfolio-backend/config"
	"portfolio-backend/models"

	"github.com/gofiber/fiber/v2"
)

// GetProjects mengembalikan semua project, urut sesuai sort_order (lalu id).
// Endpoint ini publik: dipakai frontend untuk menampilkan portfolio.
func GetProjects(c *fiber.Ctx) error {
	projects := []models.Project{} // kosong dikirim sebagai [] (bukan null)
	config.DB.Order("sort_order ASC, id ASC").Find(&projects)
	return c.JSON(projects)
}

func GetProjectByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project tidak ditemukan"})
	}

	return c.JSON(project)
}

func CreateProject(c *fiber.Ctx) error {
	var project models.Project

	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	project.Normalize()

	// Validasi field wajib
	if project.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title tidak boleh kosong"})
	}

	config.DB.Create(&project)
	return c.Status(fiber.StatusCreated).JSON(project)
}

func UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project tidak ditemukan"})
	}

	// Simpan identitas asli supaya body request tidak bisa menimpa ID / created_at.
	original := project.Model

	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	project.Model = original
	project.Normalize()

	if project.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title tidak boleh kosong"})
	}

	config.DB.Save(&project)
	return c.JSON(project)
}

func DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project tidak ditemukan"})
	}

	config.DB.Delete(&project)
	return c.JSON(fiber.Map{"message": "Project berhasil dihapus"})
}
