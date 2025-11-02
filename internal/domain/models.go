package domain

import (
    "time"
)

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Nama      string    `gorm:"size:255" json:"nama"`
    KataSandi string    `gorm:"size:255" json:"-"`
    NoTelp    string    `gorm:"size:255;uniqueIndex" json:"no_telp"`
    TanggalLahir *time.Time `json:"tanggal_lahir,omitempty"`
    JenisKelamin string  `gorm:"size:255" json:"jenis_kelamin,omitempty"`
    Tentang   string    `gorm:"type:text" json:"tentang,omitempty"`
    Pekerjaan string    `gorm:"size:255" json:"pekerjaan,omitempty"`
    Email     string    `gorm:"size:255;uniqueIndex" json:"email"`
    IdProvinsi string   `gorm:"size:255" json:"id_provinsi,omitempty"`
    IdKota    string    `gorm:"size:255" json:"id_kota,omitempty"`
    IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Toko      Toko      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"toko"`
    Alamat    []Alamat  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"alamat"`
}

type Toko struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"uniqueIndex" json:"user_id"`
    NamaToko  string    `gorm:"size:255" json:"nama_toko"`
    UrlFoto   string    `gorm:"size:255" json:"url_foto,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Produk    []Produk  `json:"produk"`
}

type Alamat struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    UserID       uint      `json:"user_id"`
    JudulAlamat  string    `gorm:"size:255" json:"judul_alamat"`
    NamaPenerima string    `gorm:"size:255" json:"nama_penerima"`
    NoTelp       string    `gorm:"size:255" json:"no_telp"`
    DetailAlamat string    `gorm:"type:text" json:"detail_alamat"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type Category struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Nama      string    `gorm:"size:255" json:"nama"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Produk    []Produk
}

type Produk struct {
    ID               uint      `gorm:"primaryKey" json:"id"`
    NamaProduk       string    `gorm:"size:255" json:"nama_produk"`
    Slug             string    `gorm:"size:255;index" json:"slug"`
    HargaReseller    int       `json:"harga_reseller"`
    HargaKonsumen    int       `json:"harga_konsumen"`
    Stok             int       `json:"stok"`
    Deskripsi        string    `gorm:"type:text" json:"deskripsi"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
    TokoID           uint      `json:"id_toko"`
    CategoryID       uint      `json:"id_category"`
    FotoProduk       []FotoProduk
}

type FotoProduk struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    ProdukID  uint      `json:"id_produk"`
    URL       string    `gorm:"size:255" json:"url"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Trx struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    UserID          uint      `json:"id_user"`
    AlamatID        uint      `json:"alamat_pengiriman"`
    HargaTotal      int       `json:"harga_total"`
    KodeInvoice     string    `gorm:"size:255" json:"kode_invoice"`
    MethodBayar     string    `gorm:"size:255" json:"method_bayar"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    DetailTrx       []DetailTrx
}

type DetailTrx struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    TrxID      uint      `json:"id_trx"`
    LogProdukID uint     `json:"id_log_produk"`
    TokoID     uint      `json:"id_toko"`
    Kuantitas  int       `json:"kuantitas"`
    HargaTotal int       `json:"harga_total"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type LogProduk struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    ProdukID        uint      `json:"id_produk"`
    NamaProduk      string    `gorm:"size:255" json:"nama_produk"`
    Slug            string    `gorm:"size:255" json:"slug"`
    HargaReseller   int       `json:"harga_reseller"`
    HargaKonsumen   int       `json:"harga_konsumen"`
    Deskripsi       string    `gorm:"type:text" json:"deskripsi"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    // relasi ke detail_trx via ID nanti
}
