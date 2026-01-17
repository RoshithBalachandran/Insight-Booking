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
	// 1. Request DTO
	type AdminLoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login payload"})
		return
	}

	// 2. Fetch admin
	var admin models.User
	if err := database.DB.
		Where("email = ? AND role = ?", req.Email, constant.Admin).
		First(&admin).Error; err != nil {

		// intentionally vague
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// 3. Account state checks
	if !admin.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin account is inactive"})
		return
	}
	if admin.Block {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin account is blocked"})
		return
	}

	// 4. Password verification
	if err := utils.CheckPassword(admin.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// 5. Generate tokens
	jti := uuid.NewString()

	accessTokenTTL, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL"))
	if accessTokenTTL == 0 {
		accessTokenTTL = 900 // 15 minutes default
	}

	refreshTokenTTL := 7 * 24 * time.Hour

	accessToken, err := token.GenerateAccessToken(admin.ID, admin.Email, admin.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	refreshToken, err := token.GenerateRefreshToken(admin.ID, admin.Role, jti)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	// 6. Store refresh session (Redis)
	if err := redis.RDB.Set(
		redis.Ctx,
		"refresh:"+jti,
		admin.ID,
		refreshTokenTTL,
	).Err(); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "session creation failed"})
		return
	}

	if err := redis.RDB.
		SAdd(redis.Ctx, "admin_sessions:"+strconv.Itoa(int(admin.ID)), jti).
		Err(); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "session tracking failed"})
		return
	}

	// 7. Secure cookies
	isProd := os.Getenv("ENV") == "production"

	c.SetCookie(
		"access_token",
		accessToken,
		accessTokenTTL,
		"/admin",
		"",
		isProd,
		true,
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(refreshTokenTTL.Seconds()),
		"/auth/refresh",
		"",
		isProd,
		true,
	)

	// 8. Optional: audit login timestamp
	_ = database.DB.
		Model(&models.User{}).
		Where("id = ?", admin.ID).
		Update("last_login_at", time.Now()).Error

	// 9. Response
	c.JSON(http.StatusOK, gin.H{
		"message": "admin login successful",
	})
}
