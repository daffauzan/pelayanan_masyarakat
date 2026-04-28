package services

import (
	"time"

	"pelayanan_publik/config"
	"pelayanan_publik/models"
)

func CreateSurat(userID uint, jenisSurat, keperluan, filePendukung string) (*models.Surat, error) {
	surat := models.Surat{
		UserID:        userID,
		JenisSurat:    jenisSurat,
		Keperluan:     keperluan,
		FilePendukung: filePendukung,
		Status:        "pending",
		SubmittedAt:   time.Now(),
	}

	if err := config.DB.Create(&surat).Error; err != nil {
		return nil, err
	}
	return &surat, nil
}

func GetSuratByUser(userID uint) ([]models.Surat, error) {
	var list []models.Surat
	err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error
	return list, err
}

func GetAllSurat() ([]models.Surat, error) {
	var list []models.Surat
	err := config.DB.Preload("User").Order("created_at desc").Find(&list).Error
	return list, err
}

func GetSuratByID(id uint) (*models.Surat, error) {
	var surat models.Surat
	err := config.DB.Preload("User").First(&surat, id).Error
	return &surat, err
}

func UpdateStatusSurat(id uint, status, catatanAdmin string) (*models.Surat, error) {
	var surat models.Surat
	if err := config.DB.First(&surat, id).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	err := config.DB.Model(&surat).Updates(map[string]interface{}{
		"status":        status,
		"catatan_admin": catatanAdmin,
		"processed_at":  &now,
	}).Error

	return &surat, err
}
