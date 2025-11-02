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
	// 🔹 1. Inisialisasi koneksi database
	db := config.InitDB()

	// 🔹 2. Migrasi model ke database
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

	// 🔹 3. Inisialisasi Fiber
	app := fiber.New()

	// 🔹 4. Inisialisasi semua handler
	authHandler := delivery.NewAuthHandler(db)
	userHandler := delivery.NewUserHandler(db)
	tokoHandler := delivery.NewTokoHandler(db)
	alamatHandler := delivery.NewAlamatHandler(db)
	kategoriHandler := delivery.NewKategoriHandler(db)
	produkHandler := delivery.NewProdukHandler(db)
	transaksiHandler := delivery.NewTransaksiHandler(db)

	// =======================
	// 🔹 5. ROUTES
	// =======================

	// --- AUTH ---
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// --- USER ---
	user := app.Group("/user", middleware.JWTProtected())
	user.Get("/profile", userHandler.GetProfile)
	user.Put("/profile", userHandler.UpdateProfile)

	// --- TOKO ---
	toko := app.Group("/toko", middleware.JWTProtected())
	toko.Get("/", tokoHandler.GetMyToko)
	toko.Put("/", tokoHandler.UpdateToko)

	// --- ALAMAT ---
	alamat := app.Group("/alamat", middleware.JWTProtected())
	alamat.Post("/", alamatHandler.CreateAlamat)
	alamat.Get("/", alamatHandler.GetAllAlamat)
	alamat.Get("/:id", alamatHandler.GetAlamatByID)
	alamat.Put("/:id", alamatHandler.UpdateAlamat)
	alamat.Delete("/:id", alamatHandler.DeleteAlamat)

	// --- kategori ---
	// Semua user bisa lihat kategori
	kategori := app.Group("/kategori", middleware.JWTProtected())
	kategori.Get("/", kategoriHandler.GetAllKategori)

	// Hanya admin yang boleh create/update/delete
	kategoriAdmin := app.Group("/kategori", middleware.JWTProtected(), middleware.IsAdmin())
	kategoriAdmin.Post("/", kategoriHandler.CreateKategori)
	kategoriAdmin.Put("/:id", kategoriHandler.UpdateKategori)
	kategoriAdmin.Delete("/:id", kategoriHandler.DeleteKategori)

	// --- PRODUK ---
	produk := app.Group("/produk", middleware.JWTProtected())
	produk.Post("/", produkHandler.CreateProduk)
	produk.Get("/", produkHandler.GetAllProduk)
	produk.Put("/:id", produkHandler.UpdateProduk)
	produk.Delete("/:id", produkHandler.DeleteProduk)

	// --- TRANSAKSI ---
	transaksi := app.Group("/transaksi", middleware.JWTProtected())
	transaksi.Post("/", transaksiHandler.CreateTransaksi)
	transaksi.Get("/", transaksiHandler.GetAllTransaksi)
	transaksi.Get("/:id", transaksiHandler.GetTransaksiDetail)


	// --- TEST ROUTE (opsional) ---
	app.Get("/profile", middleware.JWTProtected(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{
			"message": "authenticated successfully",
			"user_id": userID,
		})
	})

	// =======================
	// 🔹 6. Jalankan server
	// =======================
	log.Println("🚀 Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}