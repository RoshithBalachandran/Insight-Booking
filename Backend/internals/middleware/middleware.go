package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/insight/internals/token"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Cookie first
		if t, err := c.Cookie("access_token"); err == nil {
			tokenString = t
		}

		// Header fallback
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		claims := &token.AccessClaims{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
		})

		if err != nil || !tkn.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
