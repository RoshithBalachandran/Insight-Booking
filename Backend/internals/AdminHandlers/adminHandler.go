package adminhandlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
	"github.com/insight/internals/redis"
	"github.com/insight/internals/token"
	"github.com/insight/internals/utils"
)

func AdminLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password required"})
		return
	}

	var admin models.User
	if err := database.DB.
		Where("email = ? AND role = ?", req.Email, constant.Admin).
		First(&admin).Error; err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := utils.CheckPassword(admin.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := token.GenerateAccessToken(admin.ID, admin.Email, admin.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	jti := uuid.NewString()
	refreshToken, err := token.GenerateRefreshToken(admin.ID, admin.Role,jti)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	if err := redis.RDB.Set(redis.Ctx, "refresh:"+jti, admin.ID, 7*24*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session store failed"})
		return
	}

	redis.RDB.SAdd(redis.Ctx, "admin_sessions:"+strconv.Itoa(int(admin.ID)), jti)

	secure := os.Getenv("ENV") == "production"

	c.SetCookie("access_token", accessToken, 900, "/admin", "", secure, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/auth/refresh", "", secure, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "admin login successful",
	})
}
