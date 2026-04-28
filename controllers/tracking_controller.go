package controllers

import (
	"net/http"
	"strconv"

	"pelayanan_publik/config"
	"pelayanan_publik/models"

	"github.com/gin-gonic/gin"
)

// GET /api/tracking?reference_id=1&service_type=surat
// User atau admin melihat riwayat tracking suatu layanan
func GetTracking(c *gin.Context) {
	referenceIDStr := c.Query("reference_id")
	serviceType := c.Query("service_type")

	if referenceIDStr == "" || serviceType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "reference_id dan service_type wajib diisi"})
		return
	}

	referenceID, err := strconv.Atoi(referenceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "reference_id tidak valid"})
		return
	}

	var trackingList []models.Tracking
	if err := config.DB.
		Preload("Admin").
		Where("reference_id = ? AND service_type = ?", referenceID, serviceType).
		Order("created_at asc").
		Find(&trackingList).Error; err != nil {
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
