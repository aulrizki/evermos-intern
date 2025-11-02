package delivery

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/aulrizki/evermos-intern/internal/domain"
)

type AlamatHandler struct {
	DB *gorm.DB
}

func NewAlamatHandler(db *gorm.DB) *AlamatHandler {
	return &AlamatHandler{DB: db}
}

// 🟩 POST /alamat
func (h *AlamatHandler) CreateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req domain.Alamat
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "invalid request"})
	}

	req.UserID = userID

	if err := h.DB.Create(&req).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "failed to create address"})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  true,
		"message": "address created successfully",
		"data":    req,
	})
}

// 🟩 GET /alamat
func (h *AlamatHandler) GetAllAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var alamat []domain.Alamat
	if err := h.DB.Where("user_id = ?", userID).Find(&alamat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "failed to get addresses"})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "success get all addresses",
		"data":    alamat,
	})
}

// 🟩 GET /alamat/:id
func (h *AlamatHandler) GetAlamatByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var alamat domain.Alamat
	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).First(&alamat).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "address not found"})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "success get address",
		"data":    alamat,
	})
}

// 🟩 PUT /alamat/:id
func (h *AlamatHandler) UpdateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var req domain.Alamat
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "invalid request"})
	}

	var alamat domain.Alamat
	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).First(&alamat).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "address not found"})
	}

	alamat.JudulAlamat = req.JudulAlamat
	alamat.NamaPenerima = req.NamaPenerima
	alamat.NoTelp = req.NoTelp
	alamat.DetailAlamat = req.DetailAlamat

	if err := h.DB.Save(&alamat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "failed to update address"})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "address updated successfully",
		"data":    alamat,
	})
}

// 🟩 DELETE /alamat/:id
func (h *AlamatHandler) DeleteAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var alamat domain.Alamat
	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).First(&alamat).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "address not found"})
	}

	if err := h.DB.Delete(&alamat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "failed to delete address"})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "address deleted successfully",
	})
}
