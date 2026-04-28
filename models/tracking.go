package models

import "time"

type Tracking struct {
	ID          uint   `gorm:"primaryKey"`
	ReferenceID uint   `gorm:"not null"`
	ServiceType string `gorm:"size:50"` // surat, pengaduan
	Status      string `gorm:"size:50"`
	Keterangan  string `gorm:"type:text"`
	UpdatedBy   uint
	CreatedAt   time.Time

	Admin User `gorm:"foreignKey:UpdatedBy"`
}
