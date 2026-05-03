package main

import (
	"log"
	"os"
	"strings"

	"pelayanan_publik/config"
	"pelayanan_publik/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}

	config.ConnectDatabase()
	config.InitAWS()
	config.SeedAdmin()

	r := gin.Default()

	trustedProxyEnv := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	trustedProxies := []string{"127.0.0.1", "::1"}
	if trustedProxyEnv != "" {
		trustedProxies = strings.Split(trustedProxyEnv, ",")
		for i := range trustedProxies {
			trustedProxies[i] = strings.TrimSpace(trustedProxies[i])
		}
	}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	r.Use(config.NewSessionMiddleware())
	routes.SetupRoutes(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
