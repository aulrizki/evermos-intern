package delivery

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/aulrizki/evermos-intern/internal/domain"
)

type TokoHandler struct {
	DB *gorm.DB
}

func NewTokoHandler(db *gorm.DB) *TokoHandler {
	return &TokoHandler{DB: db}
}

// GET /toko → ambil toko milik user
func (h *TokoHandler) GetMyToko(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var toko domain.Toko
	if err := h.DB.Where("user_id = ?", userID).First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "toko tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "success get toko",
		"data":    toko,
	})
}

// PUT /toko → update nama toko + foto (upload)
func (h *TokoHandler) UpdateToko(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var toko domain.Toko
	if err := h.DB.Where("user_id = ?", userID).First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "toko tidak ditemukan",
		})
	}

	// Ambil input text field
	namaToko := c.FormValue("nama_toko")

	// Ambil file foto (optional)
	file, err := c.FormFile("foto_toko")
	if err == nil && file != nil {
		uploadPath := os.Getenv("UPLOAD_PATH")
		if uploadPath == "" {
			uploadPath = "./uploads"
		}

		// Buat nama unik untuk file
		fileName := fmt.Sprintf("toko_%d_%d%s", userID, time.Now().Unix(), filepath.Ext(file.Filename))
		fullPath := filepath.Join(uploadPath, fileName)

		// Simpan file ke folder uploads
		if err := c.SaveFile(file, fullPath); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": "gagal menyimpan file",
			})
		}

		toko.UrlFoto = fullPath
	}

	if namaToko != "" {
		toko.NamaToko = namaToko
	}

	if err := h.DB.Save(&toko).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "gagal update toko",
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "toko berhasil diupdate",
		"data":    toko,
	})
}
