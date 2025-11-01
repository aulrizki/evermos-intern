package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/aulrizki/evermos-intern/config"
	"github.com/aulrizki/evermos-intern/internal/domain"
	"github.com/aulrizki/evermos-intern/internal/delivery"
	"github.com/aulrizki/evermos-intern/internal/middleware"
)

func main() {
	// 1️⃣ Inisialisasi koneksi database
	db := config.InitDB()

	// 2️⃣ Jalankan migrasi model ke database
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Toko{},
		&domain.Alamat{},
		&domain.Category{},
		&domain.Produk{},
		&domain.FotoProduk{},
		&domain.Trx{},
		&domain.DetailTrx{},
		&domain.LogProduk{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 3️⃣ Inisialisasi Fiber app
	app := fiber.New()

	// 4️⃣ Inisialisasi handler
	authHandler := delivery.NewAuthHandler(db)
	userHandler := delivery.NewUserHandler(db)

	// 5️⃣ Routing group untuk /auth
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// 6️⃣ Routing group untuk /user
	user := app.Group("/user", middleware.JWTProtected())
	user.Get("/profile", userHandler.GetProfile)
	user.Put("/profile", userHandler.UpdateProfile)

	// 7️⃣ Contoh route test JWT (optional)
	app.Get("/profile", middleware.JWTProtected(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{
			"message": "authenticated successfully",
			"user_id": userID,
		})
	})

	// 8️⃣ Jalankan server di port 3000
	log.Println("Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
