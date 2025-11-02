package delivery

import (
	"fmt"
	"time"
	"math/rand"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/aulrizki/evermos-intern/internal/domain"
)

type TransaksiHandler struct {
	DB *gorm.DB
}

func NewTransaksiHandler(db *gorm.DB) *TransaksiHandler {
	return &TransaksiHandler{DB: db}
}

// 🟢 POST /transaksi — Buat transaksi baru
func (h *TransaksiHandler) CreateTransaksi(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req struct {
		AlamatID     uint   `json:"alamat_id"`
		MethodBayar  string `json:"method_bayar"`
		Items        []struct {
			ProdukID  uint `json:"produk_id"`
			Kuantitas int  `json:"kuantitas"`
		} `json:"items"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "invalid request"})
	}

	if len(req.Items) == 0 {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "items cannot be empty"})
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	totalHarga := 0
	kodeInvoice := fmt.Sprintf("INV-%d-%d", userID, rand.Intn(100000))

	trx := domain.Trx{
		UserID:      userID,
		AlamatID:    req.AlamatID,
		MethodBayar: req.MethodBayar,
		KodeInvoice: kodeInvoice,
		CreatedAt:   time.Now(),
	}

	if err := tx.Create(&trx).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"status": false, "message": err.Error()})
	}

	for _, item := range req.Items {
		var produk domain.Produk
		if err := tx.First(&produk, item.ProdukID).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"status": false, "message": "produk not found"})
		}

		if produk.Stok < item.Kuantitas {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"status": false, "message": fmt.Sprintf("stok produk %s tidak cukup", produk.NamaProduk)})
		}

		// Kurangi stok
		produk.Stok -= item.Kuantitas
		tx.Save(&produk)

		// Simpan log produk (snapshot)
		log := domain.LogProduk{
			ProdukID:      produk.ID,
			NamaProduk:    produk.NamaProduk,
			Slug:          produk.Slug,
			HargaReseller: produk.HargaReseller,
			HargaKonsumen: produk.HargaKonsumen,
			Deskripsi:     produk.Deskripsi,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		tx.Create(&log)

		// Detail transaksi
		detail := domain.DetailTrx{
			TrxID:      trx.ID,
			LogProdukID: log.ID,
			TokoID:     produk.TokoID,
			Kuantitas:  item.Kuantitas,
			HargaTotal: produk.HargaKonsumen * item.Kuantitas,
			CreatedAt:  time.Now(),
		}
		tx.Create(&detail)

		totalHarga += produk.HargaKonsumen * item.Kuantitas
	}

	// Update total harga transaksi
	trx.HargaTotal = totalHarga
	tx.Save(&trx)

	tx.Commit()

	return c.Status(201).JSON(fiber.Map{
		"status":  true,
		"message": "transaksi berhasil dibuat",
		"data": fiber.Map{
			"id_transaksi": trx.ID,
			"kode_invoice": trx.KodeInvoice,
			"harga_total":  trx.HargaTotal,
		},
	})
}

// 🟢 GET /transaksi — daftar transaksi milik user
func (h *TransaksiHandler) GetAllTransaksi(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "5"))
	offset := (page - 1) * limit

	var trxs []domain.Trx
	query := h.DB.Where("user_id = ?", userID).Preload("DetailTrx")

	var total int64
	query.Model(&domain.Trx{}).Count(&total)
	query.Offset(offset).Limit(limit).Find(&trxs)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "success get transaksi",
		"data": fiber.Map{
			"transaksi": trxs,
			"page":      page,
			"limit":     limit,
			"total":     total,
		},
	})
}

// 🟢 GET /transaksi/:id — detail transaksi (milik user saja)
func (h *TransaksiHandler) GetTransaksiDetail(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var trx domain.Trx
	if err := h.DB.Preload("DetailTrx").First(&trx, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "transaksi not found"})
	}

	if trx.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"status": false, "message": "forbidden: not your transaction"})
	}

	return c.JSON(fiber.Map{"status": true, "message": "success get detail transaksi", "data": trx})
}
