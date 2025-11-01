package delivery

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"

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
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "user not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "success get user profile",
		"data": fiber.Map{
			"id":          user.ID,
			"nama":        user.Nama,
			"no_telp":     user.NoTelp,
			"email":       user.Email,
			"tanggal_lahir": user.TanggalLahir,
			"jenis_kelamin": user.JenisKelamin,
			"tentang":     user.Tentang,
			"pekerjaan":   user.Pekerjaan,
			"is_admin":    user.IsAdmin,
			"created_at":  user.CreatedAt,
		},
	})
}

type UpdateUserRequest struct {
	Nama          string    `json:"nama"`
	TanggalLahir  *time.Time `json:"tanggal_lahir"`
	JenisKelamin  string    `json:"jenis_kelamin"`
	Tentang       string    `json:"tentang"`
	Pekerjaan     string    `json:"pekerjaan"`
	Email         string    `json:"email"`
}

// PUT /user/profile
// PUT /user/profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)

    var req UpdateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  false,
            "message": "invalid request",
        })
    }

    var user domain.User
    if err := h.DB.First(&user, userID).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "status":  false,
            "message": "user not found",
        })
    }

    // 🟩 BAGIAN INI ADALAH "UPDATE SEBELUM SAVE"
    // hanya update field yang dikirim dari body
    user.Nama          = req.Nama
    user.TanggalLahir  = req.TanggalLahir
    user.JenisKelamin  = req.JenisKelamin
    user.Tentang       = req.Tentang
    user.Pekerjaan     = req.Pekerjaan
    user.Email         = req.Email
    // 🟩 Sampai sini kita baru "mengubah nilai dalam struct user" di memori

    // 🟦 Setelah diubah, baru disimpan ke database:
    if err := h.DB.Save(&user).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  false,
            "message": "failed to update user",
        })
    }

    // 🟨 Response setelah berhasil update
    return c.JSON(fiber.Map{
        "status":  true,
        "message": "profile updated successfully",
        "data":    user,
    })
}

