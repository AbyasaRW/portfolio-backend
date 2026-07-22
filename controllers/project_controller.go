package controllers

import (
	"portfolio-backend/config"
	"portfolio-backend/models"

	"github.com/gofiber/fiber/v2"
)

func GetProjects(c *fiber.Ctx) error {
	var projects []models.Project
	config.DB.Find(&projects)
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

	config.DB.Create(&project)
	return c.Status(fiber.StatusCreated).JSON(project)
}

func UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project tidak ditemukan"})
	}

	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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
