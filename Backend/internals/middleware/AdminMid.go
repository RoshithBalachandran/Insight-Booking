package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/internals/constant"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != constant.Admin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
