package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Nama      string `gorm:"column:nama;size:100;not null"`
	Email     string `gorm:"column:email;size:100;not null;unique"`
	Password  string `gorm:"column:password;size:255;not null"`
	NIK       string `gorm:"column:nik;unique;size:16"`
	NoTelp    string `gorm:"column:no_telp;size:20"`
	Alamat    string `gorm:"column:alamat;type:text"`
	Role      string `gorm:"type:enum('admin','user');default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
