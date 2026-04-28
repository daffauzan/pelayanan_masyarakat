package controllers

import (
	"net/http"
	"strconv"
	"time"

	"pelayanan_publik/config"
	"pelayanan_publik/models"

	"github.com/gin-gonic/gin"
)

// POST /api/surat — user mengajukan surat
func CreateSurat(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		JenisSurat    string `json:"jenis_surat"    binding:"required"`
		Keperluan     string `json:"keperluan"`
		FilePendukung string `json:"file_pendukung"` // URL file (dari S3 jika sudah diintegrasikan)
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	surat := models.Surat{
		UserID:        userID,
		JenisSurat:    input.JenisSurat,
		Keperluan:     input.Keperluan,
		FilePendukung: input.FilePendukung,
		Status:        "pending",
		SubmittedAt:   time.Now(),
	}

	if err := config.DB.Create(&surat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengajukan surat"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "surat berhasil diajukan", "data": surat})
}

// GET /api/surat — user melihat surat miliknya sendiri
func GetMySurat(c *gin.Context) {
	userID := c.GetUint("user_id")

	var suratList []models.Surat
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&suratList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data surat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": suratList})
}

// GET /api/admin/surat — admin melihat semua surat
func GetAllSurat(c *gin.Context) {
	var suratList []models.Surat
	if err := config.DB.Preload("User").Order("created_at desc").Find(&suratList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data surat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": suratList})
}

// GET /api/surat/:id — detail surat
func GetSuratByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}

	var surat models.Surat
	if err := config.DB.Preload("User").First(&surat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "surat tidak ditemukan"})
		return
	}

	// user hanya bisa lihat suratnya sendiri
	if role != "admin" && surat.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "akses ditolak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": surat})
}

// PUT /api/admin/surat/:id — admin update status & catatan
func UpdateStatusSurat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}

	var input struct {
		Status       string `json:"status"        binding:"required"`
		CatatanAdmin string `json:"catatan_admin"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var surat models.Surat
	if err := config.DB.First(&surat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "surat tidak ditemukan"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":        input.Status,
		"catatan_admin": input.CatatanAdmin,
		"processed_at":  &now,
	}

	if err := config.DB.Model(&surat).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal update status surat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status surat diperbarui", "data": surat})
}
