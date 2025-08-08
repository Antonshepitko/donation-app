// middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"donation-backend/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token"})
			c.Abort()
			return
		}
		raw := strings.TrimPrefix(auth, "Bearer ")
		username, err := utils.ParseJWT(raw) // ← было ParseToken
		if err != nil || username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
