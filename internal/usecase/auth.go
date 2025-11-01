package usecase

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "github.com/aulrizki/evermos-intern/internal/domain"
)

func Register(db *gorm.DB, nama, noTelp, password, email string) (*domain.User, error) {
    // cek duplicate noTelp
    var exists domain.User
    if err := db.Where("no_telp = ?", noTelp).First(&exists).Error; err == nil {
        return nil, errors.New("no_telp sudah terdaftar")
    }

    // hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := domain.User{
        Nama: nama,
        NoTelp: noTelp,
        KataSandi: string(hashed),
        Email: email,
    }

    // transaksional: create user + create toko
    err = db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&user).Error; err != nil {
            return err
        }
        toko := domain.Toko{
            UserID: user.ID,
            NamaToko: nama + "'s store",
        }
        if err := tx.Create(&toko).Error; err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        return nil, err
    }
    return &user, nil
}
