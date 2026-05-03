package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const SessionName = "pelayanan_publik_session"

var (
	sessionSecret     []byte
	sessionSecretOnce sync.Once
)

func getSessionSecret() []byte {
	sessionSecretOnce.Do(func() {
		secret := os.Getenv("SESSION_SECRET")

		if secret != "" {
			sessionSecret = []byte(secret)
			return
		}

		buf := make([]byte, 32)
		if _, err := rand.Read(buf); err != nil {
			sessionSecret = []byte("pelayanan-publik-session-fallback")
			log.Println("warning: SESSION_SECRET is not set, using fallback runtime session secret")
			return
		}

		sessionSecret = []byte(base64.StdEncoding.EncodeToString(buf))
		log.Println("warning: SESSION_SECRET is not set, using auto-generated runtime session secret")
	})

	return sessionSecret
}

func getSessionMaxAge() int {
	maxAgeHours := 24

	if value := os.Getenv("SESSION_EXPIRED_HOURS"); value != "" {
		if hours, err := strconv.Atoi(value); err == nil && hours > 0 {
			maxAgeHours = hours
		}
	}

	return int((time.Duration(maxAgeHours) * time.Hour).Seconds())
}

func NewSessionMiddleware() gin.HandlerFunc {
	store := cookie.NewStore(getSessionSecret())
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   getSessionMaxAge(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})

	return sessions.Sessions(SessionName, store)
}
