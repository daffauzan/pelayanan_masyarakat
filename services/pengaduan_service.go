package services

import (
	"time"

	"pelayanan_publik/config"
	"pelayanan_publik/models"
)

func CreatePengaduan(userID uint, judul, deskripsi, kategori, lampiran string) (*models.Pengaduan, error) {
	var lampiranPtr *string
	if lampiran != "" {
		lampiranPtr = &lampiran
	}

	pengaduan := models.Pengaduan{
		UserID:    userID,
		Judul:     judul,
		Deskripsi: deskripsi,
		Kategori:  kategori,
		Lampiran:  lampiranPtr,
		Status:    "open",
	}

	if err := config.DB.Create(&pengaduan).Error; err != nil {
		return nil, err
	}
	return &pengaduan, nil
}

func GetPengaduanByUser(userID uint) ([]models.Pengaduan, error) {
	var list []models.Pengaduan
	err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error
	return list, err
}

func GetAllPengaduan() ([]models.Pengaduan, error) {
	var list []models.Pengaduan
	err := config.DB.Preload("User").Order("created_at desc").Find(&list).Error
	return list, err
}

func GetPengaduanByID(id uint) (*models.Pengaduan, error) {
	var pengaduan models.Pengaduan
	err := config.DB.Preload("User").First(&pengaduan, id).Error
	return &pengaduan, err
}

func UpdateStatusPengaduan(id uint, status, tanggapan string) (*models.Pengaduan, error) {
	var pengaduan models.Pengaduan
	if err := config.DB.First(&pengaduan, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"status":          status,
		"tanggapan_admin": tanggapan,
	}
	if status == "resolved" {
		now := time.Now()
		updates["resolved_at"] = &now
	}

	err := config.DB.Model(&pengaduan).Updates(updates).Error
	return &pengaduan, err
}
