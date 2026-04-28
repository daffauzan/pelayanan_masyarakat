package routes

import (
	"pelayanan_publik/controllers"
	"pelayanan_publik/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")

	// Auth — public
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// User routes — perlu login
	user := api.Group("/")
	user.Use(middleware.AuthMiddleware())
	{
		// Surat
		user.POST("/surat", controllers.CreateSurat)
		user.GET("/surat", controllers.GetMySurat)
		user.GET("/surat/:id", controllers.GetSuratByID)

		// Pengaduan
		user.POST("/pengaduan", controllers.CreatePengaduan)
		user.GET("/pengaduan", controllers.GetMyPengaduan)
		user.GET("/pengaduan/:id", controllers.GetPengaduanByID)

		// Tracking
		user.GET("/tracking", controllers.GetTracking)
	}

	// Admin routes — perlu login + role admin
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		// Surat
		admin.GET("/surat", controllers.GetAllSurat)
		admin.PUT("/surat/:id", controllers.UpdateStatusSurat)

		// Pengaduan
		admin.GET("/pengaduan", controllers.GetAllPengaduan)
		admin.PUT("/pengaduan/:id", controllers.UpdateStatusPengaduan)

		// Tracking
		admin.POST("/tracking", controllers.AddTracking)
	}
}
