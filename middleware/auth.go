package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userIDValue := session.Get("user_id")
		roleValue := session.Get("role")

		userID, ok := toUint(userIDValue)
		if !ok || roleValue == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "session tidak valid atau belum login",
			})
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok || role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "session role tidak valid",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}

func toUint(value interface{}) (uint, bool) {
	switch typed := value.(type) {
	case uint:
		return typed, true
	case uint8:
		return uint(typed), true
	case uint16:
		return uint(typed), true
	case uint32:
		return uint(typed), true
	case uint64:
		return uint(typed), true
	case int:
		if typed < 0 {
			return 0, false
		}
		return uint(typed), true
	case int64:
		if typed < 0 {
			return 0, false
		}
		return uint(typed), true
	case float64:
		if typed < 0 {
			return 0, false
		}
		return uint(typed), true
	case string:
		parsed, err := strconv.ParseUint(typed, 10, 64)
		if err != nil {
			return 0, false
		}
		return uint(parsed), true
	default:
		return 0, false
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "akses hanya untuk admin"})
			c.Abort()
			return
		}
		c.Next()
	}
}
