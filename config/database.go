package config

import (
	"fmt"
	"log"
	"os"

	"pelayanan_publik/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	err = database.AutoMigrate(
		&models.User{},
		&models.Surat{},
		&models.Pengaduan{},
		&models.Tracking{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	DB = database
	log.Println("MySQL connected successfully")
}
