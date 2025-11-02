package delivery

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/aulrizki/evermos-intern/internal/domain"
)

type ProdukHandler struct {
	DB *gorm.DB
}

func NewProdukHandler(db *gorm.DB) *ProdukHandler {
	return &ProdukHandler{DB: db}
}

// 🟩 CREATE Produk (POST /produk)
func (h *ProdukHandler) CreateProduk(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Ambil toko milik user
	var toko domain.Toko
	if err := h.DB.Where("user_id = ?", userID).First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "toko not found"})
	}

	// Ambil form data
	nama := c.FormValue("nama_produk")
	slug := strings.ReplaceAll(strings.ToLower(nama), " ", "-")
	hargaReseller, _ := strconv.Atoi(c.FormValue("harga_reseller"))
	hargaKonsumen, _ := strconv.Atoi(c.FormValue("harga_konsumen"))
	stok, _ := strconv.Atoi(c.FormValue("stok"))
	deskripsi := c.FormValue("deskripsi")
	categoryID, _ := strconv.Atoi(c.FormValue("id_category"))

	if nama == "" || hargaKonsumen == 0 {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "nama dan harga wajib diisi"})
	}

	produk := domain.Produk{
		NamaProduk:    nama,
		Slug:          slug,
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok:          stok,
		Deskripsi:     deskripsi,
		TokoID:        toko.ID,
		CategoryID:    uint(categoryID),
	}

	// Simpan produk
	if err := h.DB.Create(&produk).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "failed to create product"})
	}

	// Upload multiple foto produk
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["foto_produk"]
		for _, file := range files {
			filename := fmt.Sprintf("uploads/products/%d_%d_%s", userID, time.Now().Unix(), file.Filename)
			if err := c.SaveFile(file, filename); err == nil {
				foto := domain.FotoProduk{
					ProdukID: produk.ID,
					URL:      filename,
				}
				h.DB.Create(&foto)
			}
		}
	}

	h.DB.Preload("FotoProduk").First(&produk)

	return c.Status(201).JSON(fiber.Map{
		"status":  true,
		"message": "produk berhasil dibuat",
		"data":    produk,
	})
}

// 🟩 READ Produk (GET /produk?page=1&limit=5&nama=kaos)
func (h *ProdukHandler) GetAllProduk(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "5"))
	offset := (page - 1) * limit

	var produk []domain.Produk
	query := h.DB.Preload("FotoProduk")

	// Filtering
	if nama := c.Query("nama"); nama != "" {
		query = query.Where("nama_produk LIKE ?", "%"+nama+"%")
	}
	if category := c.Query("id_category"); category != "" {
		query = query.Where("category_id = ?", category)
	}
	if minHarga := c.Query("min_harga"); minHarga != "" {
		query = query.Where("harga_konsumen >= ?", minHarga)
	}
	if maxHarga := c.Query("max_harga"); maxHarga != "" {
		query = query.Where("harga_konsumen <= ?", maxHarga)
	}

	var total int64
	query.Model(&domain.Produk{}).Count(&total)
	query.Offset(offset).Limit(limit).Find(&produk)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "success get produk",
		"data": fiber.Map{
			"produk": produk,
			"page":   page,
			"limit":  limit,
			"total":  total,
		},
	})
}

// 🟩 UPDATE Produk (PUT /produk/:id)
func (h *ProdukHandler) UpdateProduk(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var produk domain.Produk
	if err := h.DB.Preload("FotoProduk").First(&produk, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "produk not found"})
	}

	// Cek kepemilikan toko
	var toko domain.Toko
	h.DB.First(&toko, produk.TokoID)
	if toko.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"status": false, "message": "forbidden: not your product"})
	}

	// Update field yang dikirim
	if nama := c.FormValue("nama_produk"); nama != "" {
		produk.NamaProduk = nama
		produk.Slug = strings.ReplaceAll(strings.ToLower(nama), " ", "-")
	}
	if harga := c.FormValue("harga_konsumen"); harga != "" {
		hargaInt, _ := strconv.Atoi(harga)
		produk.HargaKonsumen = hargaInt
	}
	if stok := c.FormValue("stok"); stok != "" {
		stokInt, _ := strconv.Atoi(stok)
		produk.Stok = stokInt
	}
	if deskripsi := c.FormValue("deskripsi"); deskripsi != "" {
		produk.Deskripsi = deskripsi
	}

	// Upload foto baru (opsional)
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["foto_produk"]
		for _, file := range files {
			filename := fmt.Sprintf("uploads/products/%d_%d_%s", userID, time.Now().Unix(), file.Filename)
			if err := c.SaveFile(file, filename); err == nil {
				foto := domain.FotoProduk{
					ProdukID: produk.ID,
					URL:      filename,
				}
				h.DB.Create(&foto)
			}
		}
	}

	h.DB.Save(&produk)
	h.DB.Preload("FotoProduk").First(&produk)

	return c.JSON(fiber.Map{"status": true, "message": "produk updated successfully", "data": produk})
}

// 🟩 DELETE Produk (DELETE /produk/:id)
func (h *ProdukHandler) DeleteProduk(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var produk domain.Produk
	if err := h.DB.First(&produk, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "produk not found"})
	}

	var toko domain.Toko
	h.DB.First(&toko, produk.TokoID)
	if toko.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"status": false, "message": "forbidden: not your product"})
	}

	h.DB.Delete(&produk)
	return c.JSON(fiber.Map{"status": true, "message": "produk deleted successfully"})
}
