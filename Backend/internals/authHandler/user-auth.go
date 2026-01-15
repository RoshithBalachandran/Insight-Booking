package authHandler

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
	"github.com/insight/internals/redis"
	"github.com/insight/internals/token"
	"github.com/insight/internals/utils"
)

/* ================= REGISTER ================= */

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println("Password hashing failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if req.Role == "" {
		req.Role = constant.User
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hash,
		Role:      req.Role,
	}

	// Set PhoneNumber only if provided
	if req.Phone != "" {
		user.PhoneNumber = &req.Phone
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Println("User creation failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registered successfully"})
}

/* ================= LOGIN ================= */

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := token.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		log.Println("Access token generation failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	jti := uuid.New().String()
	refreshToken, err := token.GenerateRefreshToken(user.ID, user.Role, jti)
	if err != nil {
		log.Println("Refresh token generation failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	userKey := "user_sessions:" + strconv.Itoa(int(user.ID))
	count, _ := redis.RDB.SCard(redis.Ctx, userKey).Result()
	if count >= 3 {
		oldJTI, _ := redis.RDB.SPop(redis.Ctx, userKey).Result()
		redis.RDB.Del(redis.Ctx, "refresh:"+oldJTI)
	}

	redis.RDB.Set(redis.Ctx, "refresh:"+jti, user.ID, 7*24*time.Hour)
	redis.RDB.SAdd(redis.Ctx, userKey, jti)
	redis.RDB.Expire(redis.Ctx, userKey, 7*24*time.Hour)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, 900, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/auth/refresh", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    user.FirstName,
		"email":   user.Email,
	})
}

/* ================= REFRESH ================= */

func Refresh(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	claims := &token.RefreshClaims{}
	tkn, err := jwt.ParseWithClaims(rt, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	if err != nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	key := "refresh:" + claims.JTI
	if _, err := redis.RDB.Get(redis.Ctx, key).Result(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
		return
	}

	redis.RDB.Del(redis.Ctx, key)
	redis.RDB.SRem(redis.Ctx, "user_sessions:"+strconv.Itoa(int(claims.UserID)), claims.JTI)

	newAccess, _ := token.GenerateAccessToken(claims.UserID, "", "")
	newJTI := uuid.New().String()
	newRefresh, _ := token.GenerateRefreshToken(claims.UserID, claims.Role, newJTI)

	redis.RDB.Set(redis.Ctx, "refresh:"+newJTI, claims.UserID, 7*24*time.Hour)
	redis.RDB.SAdd(redis.Ctx, "user_sessions:"+strconv.Itoa(int(claims.UserID)), newJTI)

	c.SetCookie("access_token", newAccess, 900, "/", "", true, true)
	c.SetCookie("refresh_token", newRefresh, 604800, "/auth/refresh", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}

/* ================= LOGOUT ================= */

func Logout(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err == nil {
		claims := &token.RefreshClaims{}
		jwt.ParseWithClaims(rt, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
		})
		redis.RDB.Del(redis.Ctx, "refresh:"+claims.JTI)
		redis.RDB.SRem(redis.Ctx, "user_sessions:"+strconv.Itoa(int(claims.UserID)), claims.JTI)
	}

	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/auth/refresh", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
