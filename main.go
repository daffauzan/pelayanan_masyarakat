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
	// Load .env hanya untuk local development
	// Di ECS, environment variable langsung disediakan oleh container
	if os.Getenv("AWS_EXECUTION_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not loaded: %v", err)
		}
	} else {
		log.Println("Running on AWS ECS, skipping .env loading")
	}

	config.ConnectDatabase()
	config.InitAWS()
	config.SeedAdmin()

	r := gin.Default()

	// Trusted Proxies
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

	// Session Middleware
	r.Use(config.NewSessionMiddleware())

	// Routes
	routes.SetupRoutes(r)

	// Port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
