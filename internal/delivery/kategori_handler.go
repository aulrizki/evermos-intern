package delivery

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/aulrizki/evermos-intern/internal/domain"
)

type KategoriHandler struct {
	DB *gorm.DB
}

func NewKategoriHandler(db *gorm.DB) *KategoriHandler {
	return &KategoriHandler{DB: db}
}

// 🟩 POST /kategori → hanya admin
func (h *KategoriHandler) CreateKategori(c *fiber.Ctx) error {
	var req domain.Category
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "invalid request",
		})
	}

	if req.Nama == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "category name is required",
		})
	}

	if err := h.DB.Create(&req).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "failed to create category",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  true,
		"message": "category created successfully",
		"data":    req,
	})
}

// 🟩 GET /kategori → semua user boleh
func (h *KategoriHandler) GetAllKategori(c *fiber.Ctx) error {
	var kategori []domain.Category
	if err := h.DB.Find(&kategori).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "failed to get categories",
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "success get categories",
		"data":    kategori,
	})
}

// 🟩 PUT /kategori/:id → hanya admin
func (h *KategoriHandler) UpdateKategori(c *fiber.Ctx) error {
	id := c.Params("id")
	var req domain.Category
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "invalid request"})
	}

	var kategori domain.Category
	if err := h.DB.First(&kategori, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "category not found"})
	}

	kategori.Nama = req.Nama

	if err := h.DB.Save(&kategori).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "failed to update category"})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "category updated successfully",
		"data":    kategori,
	})
}

// 🟩 DELETE /kategori/:id → hanya admin
func (h *KategoriHandler) DeleteKategori(c *fiber.Ctx) error {
	id := c.Params("id")

	var kategori domain.Category
	if err := h.DB.First(&kategori, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "category not found"})
	}

	if err := h.DB.Delete(&kategori).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "failed to delete category"})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "category deleted successfully",
	})
}
