package delivery

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/aulrizki/evermos-intern/internal/domain"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// GET /user/profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var user domain.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "success get user profile",
		"data": fiber.Map{
			"id":            user.ID,
			"nama":          user.Nama,
			"no_telp":       user.NoTelp,
			"email":         user.Email,
			"tanggal_lahir": user.TanggalLahir,
			"jenis_kelamin": user.JenisKelamin,
			"tentang":       user.Tentang,
			"pekerjaan":     user.Pekerjaan,
			"is_admin":      user.IsAdmin,
			"created_at":    user.CreatedAt,
			"updated_at":    user.UpdatedAt,
		},
	})
}

// Struct untuk request update
type UpdateUserRequest struct {
	Nama          string     `json:"nama"`
	TanggalLahir  *time.Time `json:"tanggal_lahir"`
	JenisKelamin  string     `json:"jenis_kelamin"`
	Tentang       string     `json:"tentang"`
	Pekerjaan     string     `json:"pekerjaan"`
	Email         string     `json:"email"`
}

// PUT /user/profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid request body",
		})
	}

	var user domain.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "user not found",
		})
	}

	// ✅ Validasi: jika email diubah, pastikan belum digunakan user lain
	if req.Email != "" && req.Email != user.Email {
		var existing domain.User
		if err := h.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existing).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "gunakan email lain yang belum digunakan",
			})
		}
	}

	// Update field yang diinput
	if req.Nama != "" {
		user.Nama = req.Nama
	}
	if req.TanggalLahir != nil {
		user.TanggalLahir = req.TanggalLahir
	}
	if req.JenisKelamin != "" {
		user.JenisKelamin = req.JenisKelamin
	}
	if req.Tentang != "" {
		user.Tentang = req.Tentang
	}
	if req.Pekerjaan != "" {
		user.Pekerjaan = req.Pekerjaan
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := h.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "gagal update profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "profile berhasil diupdate",
		"data": fiber.Map{
			"id":            user.ID,
			"nama":          user.Nama,
			"no_telp":       user.NoTelp,
			"email":         user.Email,
			"tanggal_lahir": user.TanggalLahir,
			"jenis_kelamin": user.JenisKelamin,
			"tentang":       user.Tentang,
			"pekerjaan":     user.Pekerjaan,
			"is_admin":      user.IsAdmin,
			"updated_at":    user.UpdatedAt,
		},
	})
}
