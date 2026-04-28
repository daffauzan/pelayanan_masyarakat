package models

import "time"

type Pengaduan struct {
	ID             uint   `gorm:"primaryKey"`
	UserID         uint   `gorm:"not null"`
	Judul          string `gorm:"size:255;not null"`
	Deskripsi      string `gorm:"type:text"`
	Kategori       string `gorm:"size:100"`
	Lampiran       string
	Status         string `gorm:"default:'open'"`
	TanggapanAdmin string `gorm:"type:text"`
	ResolvedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time

	User User `gorm:"foreignKey:UserID"`
}
