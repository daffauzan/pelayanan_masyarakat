package controllers

import (
	"log"
	"net/http"

	"pelayanan_publik/config"
	"pelayanan_publik/models"

	"github.com/gin-contrib/sessions"
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
		NoTelp   string `json:"no_telp"`
		Alamat   string `json:"alamat"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("register validation failed: %v", err)
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
		NoTelp:   input.NoTelp,
		Alamat:   input.Alamat,
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		log.Printf("register create user failed: %v", err)
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
		log.Printf("Login: email '%s' tidak ditemukan: %v", input.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("Login: password salah untuk email '%s': %v", input.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "email atau password salah"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("role", user.Role)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal membuat session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login berhasil",
		"data": gin.H{
			"id":    user.ID,
			"nama":  user.Nama,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Logout godoc
// POST /api/auth/logout
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal menghapus session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
}

// Profile godoc
// GET /api/auth/profile
func Profile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profil berhasil diambil",
		"data": gin.H{
			"id":      user.ID,
			"nama":    user.Nama,
			"email":   user.Email,
			"no_telp": user.NoTelp,
			"alamat":  user.Alamat,
			"role":    user.Role,
		},
	})
}
