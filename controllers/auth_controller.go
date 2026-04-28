package controllers

import (
	"net/http"

	"pelayanan_publik/config"
	"pelayanan_publik/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// POST /api/auth/register
func Register(c *gin.Context) {
	var input struct {
		Nama     string `json:"nama"     binding:"required"`
		Email    string `json:"email"    binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		NIK      string `json:"nik"      binding:"required,len=16"`
		NoTelp   string `json:"no_telp"`
		Alamat   string `json:"alamat"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Cek apakah email sudah terdaftar
	var existing models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "email sudah terdaftar"})
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal memproses password"})
		return
	}

	user := models.User{
		Nama:     input.Nama,
		Email:    input.Email,
		Password: string(hashed),
		NIK:      input.NIK,
		NoTelp:   input.NoTelp,
		Alamat:   input.Alamat,
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mendaftarkan user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registrasi berhasil",
		"data": gin.H{
			"id":    user.ID,
			"nama":  user.Nama,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Login godoc
// POST /api/auth/login
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"    binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "email atau password salah"})
		return
	}

	token, err := config.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login berhasil",
		"token":   token,
		"data": gin.H{
			"id":    user.ID,
			"nama":  user.Nama,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
