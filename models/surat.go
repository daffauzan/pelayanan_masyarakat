package models

import "time"

type Surat struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        uint   `gorm:"not null"`
	JenisSurat    string `gorm:"size:100;not null"`
	Keperluan     string `gorm:"type:text"`
	FilePendukung string
	Status        string `gorm:"default:'pending'"`
	CatatanAdmin  string `gorm:"type:text;nullable"`
	SubmittedAt   time.Time
	ProcessedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time

	User User `gorm:"foreignKey:UserID"`
}
