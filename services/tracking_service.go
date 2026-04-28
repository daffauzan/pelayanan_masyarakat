package services

import (
	"pelayanan_publik/config"
	"pelayanan_publik/models"
)

func AddTracking(referenceID uint, serviceType, status, keterangan string, adminID uint) (*models.Tracking, error) {
	tracking := models.Tracking{
		ReferenceID: referenceID,
		ServiceType: serviceType,
		Status:      status,
		Keterangan:  keterangan,
		UpdatedBy:   adminID,
	}

	if err := config.DB.Create(&tracking).Error; err != nil {
		return nil, err
	}
	return &tracking, nil
}

func GetTracking(referenceID uint, serviceType string) ([]models.Tracking, error) {
	var list []models.Tracking
	err := config.DB.
		Preload("Admin").
		Where("reference_id = ? AND service_type = ?", referenceID, serviceType).
		Order("created_at asc").
		Find(&list).Error
	return list, err
}
