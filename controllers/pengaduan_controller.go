package controllers

import (
	"net/http"
	"strconv"
	"time"

	"pelayanan_publik/config"
	"pelayanan_publik/models"

	"github.com/gin-gonic/gin"
)

// POST /api/pengaduan — user membuat pengaduan
func CreatePengaduan(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		Judul     string `json:"judul"     binding:"required"`
		Deskripsi string `json:"deskripsi"`
		Kategori  string `json:"kategori"`
		Lampiran  string `json:"lampiran"` // URL file (dari S3 jika sudah diintegrasikan)
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	pengaduan := models.Pengaduan{
		UserID:    userID,
		Judul:     input.Judul,
		Deskripsi: input.Deskripsi,
		Kategori:  input.Kategori,
		Lampiran:  input.Lampiran,
		Status:    "open",
	}

	if err := config.DB.Create(&pengaduan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal membuat pengaduan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "pengaduan berhasil dikirim", "data": pengaduan})
}

// GET /api/pengaduan — user melihat pengaduan miliknya
func GetMyPengaduan(c *gin.Context) {
	userID := c.GetUint("user_id")

	var list []models.Pengaduan
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data pengaduan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// GET /api/admin/pengaduan — admin melihat semua pengaduan
func GetAllPengaduan(c *gin.Context) {
	var list []models.Pengaduan
	if err := config.DB.Preload("User").Order("created_at desc").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mengambil data pengaduan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

// GET /api/pengaduan/:id — detail pengaduan
func GetPengaduanByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}

	var pengaduan models.Pengaduan
	if err := config.DB.Preload("User").First(&pengaduan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "pengaduan tidak ditemukan"})
		return
	}

	if role != "admin" && pengaduan.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "akses ditolak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pengaduan})
}

// PUT /api/admin/pengaduan/:id — admin menanggapi dan update status
func UpdateStatusPengaduan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id tidak valid"})
		return
	}

	var input struct {
		Status         string `json:"status"          binding:"required"`
		TanggapanAdmin string `json:"tanggapan_admin"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var pengaduan models.Pengaduan
	if err := config.DB.First(&pengaduan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "pengaduan tidak ditemukan"})
		return
	}

	updates := map[string]interface{}{
		"status":          input.Status,
		"tanggapan_admin": input.TanggapanAdmin,
	}

	if input.Status == "resolved" {
		now := time.Now()
		updates["resolved_at"] = &now
	}

	if err := config.DB.Model(&pengaduan).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal update status pengaduan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pengaduan diperbarui", "data": pengaduan})
}
