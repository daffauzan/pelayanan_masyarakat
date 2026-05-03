package config

import (
	"log"

	"pelayanan_publik/models"

	"golang.org/x/crypto/bcrypt"
)

const (
	adminNama     = "admin"
	adminEmail    = "admin@email.com"
	adminPassword = "admin123"
)

func SeedAdmin() {
	nama := adminNama
	email := adminEmail
	password := adminPassword

	var existing models.User
	if err := DB.Where("email = ?", email).First(&existing).Error; err == nil {
		log.Printf("SeedAdmin: admin '%s' sudah ada, seeder dilewati", email)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("SeedAdmin: gagal hash password: %v", err)
	}

	admin := models.User{
		Nama:     nama,
		Email:    email,
		Password: string(hashed),
		Role:     "admin",
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Fatalf("SeedAdmin: gagal membuat admin: %v", err)
	}

	log.Printf("SeedAdmin: admin '%s' berhasil dibuat", email)
}
