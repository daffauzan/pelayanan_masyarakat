package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getAllowedOrigins() []string {
	origins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	}

	if frontendURL := strings.TrimSpace(os.Getenv("FRONTEND_URL")); frontendURL != "" {
		for _, origin := range strings.Split(frontendURL, ",") {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				origins = append(origins, origin)
			}
		}
	}

	return origins
}

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: getAllowedOrigins(),
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
