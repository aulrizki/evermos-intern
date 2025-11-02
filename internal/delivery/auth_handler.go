package delivery

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"

    "github.com/aulrizki/evermos-intern/internal/domain"
    "github.com/aulrizki/evermos-intern/internal/usecase"
    "github.com/aulrizki/evermos-intern/internal/utils"
)

type AuthHandler struct {
    DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
    return &AuthHandler{DB: db}
}

type RegisterRequest struct {
    Nama     string `json:"nama" validate:"required"`
    NoTelp   string `json:"no_telp" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
    }

    user, err := usecase.Register(h.DB, req.Nama, req.NoTelp, req.Password, req.Email)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(201).JSON(fiber.Map{
        "message": "register berhasil",
        "data":    user,
    })
}

type LoginRequest struct {
    NoTelp   string `json:"no_telp"`
    Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
    }

    var user domain.User
    if err := h.DB.Where("no_telp = ?", req.NoTelp).First(&user).Error; err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
    }

    if bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(req.Password)) != nil {
        return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
    }

    // 🟢 Tambahkan ini untuk memastikan nilai terbaru dari database
    h.DB.First(&user, user.ID)
    // Cetak ke console untuk memastikan
    println("Login as:", user.Nama, " | IsAdmin:", user.IsAdmin)

    token, err := utils.GenerateJWT(user.ID, user.Email, user.IsAdmin)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "token generation failed"})
    }

    return c.JSON(fiber.Map{
        "message": "login success",
        "token":   token,
    })
}

