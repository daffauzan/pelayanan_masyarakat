package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"pelayanan_publik/config"
	"pelayanan_publik/models"

	"github.com/gin-gonic/gin"
)

// GET /api/tracking?reference_id=1&service_type=surat
// User atau admin melihat riwayat tracking.
// Query parameter bersifat opsional; jika kosong, akan mengembalikan semua data yang diizinkan.
func GetTracking(c *gin.Context) {
	referenceIDStr := strings.TrimSpace(c.Query("reference_id"))
	serviceType := strings.TrimSpace(c.Query("service_type"))
	role := c.GetString("role")
	userID := c.GetUint("user_id")

	query := config.DB.Preload("Admin").Model(&models.Tracking{})

	if serviceType != "" {
		if serviceType != "surat" && serviceType != "pengaduan" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "service_type harus 'surat' atau 'pengaduan'"})
			return
		}
		query = query.Where("service_type = ?", serviceType)
	}

	if referenceIDStr != "" {
		referenceID, err := strconv.Atoi(referenceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "reference_id tidak valid"})
			return
		}
		query = query.Where("reference_id = ?", referenceID)
	}

	if role != "admin" {
		suratSubQuery := config.DB.Model(&models.Surat{}).Select("id").Where("user_id = ?", userID)
		pengaduanSubQuery := config.DB.Model(&models.Pengaduan{}).Select("id").Where("user_id = ?", userID)

		query = query.Where(
			"(service_type = ? AND reference_id IN (?)) OR (service_type = ? AND reference_id IN (?))",
			"surat", suratSubQuery,
			"pengaduan", pengaduanSubQuery,
		)
	}

	var trackingList []models.Tracking
	if err := query.Order("created_at asc").Find(&trackingList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data tracking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trackingList})
}

// POST /api/admin/tracking — admin menambahkan entri tracking
func AddTracking(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var input struct {
		ReferenceID uint   `json:"reference_id"  binding:"required"`
		ServiceType string `json:"service_type"  binding:"required"`
		Status      string `json:"status"        binding:"required"`
		Keterangan  string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tracking := models.Tracking{
		ReferenceID: input.ReferenceID,
		ServiceType: input.ServiceType,
		Status:      input.Status,
		Keterangan:  input.Keterangan,
		UpdatedBy:   adminID,
	}

	if err := config.DB.Create(&tracking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menambahkan tracking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "tracking berhasil ditambahkan", "data": tracking})
}
