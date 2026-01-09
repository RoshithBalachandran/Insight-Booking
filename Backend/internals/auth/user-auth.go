package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
	"github.com/insight/internals/token"
	"github.com/insight/internals/utils"
	"gorm.io/gorm"
)

// Registeration
func Register(c *gin.Context) {
	var req models.RegisterRequest

	//Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Check if email already exists
	var existingUser models.User
	err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email already registered",
		})
		return
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
		return
	}

	//Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to secure password",
		})
		return
	}

	//Create user
	if req.Role == "" {
		req.Role = constant.User
	}
	user := models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.Phone,
		Password:    hashedPassword,
		Role:        req.Role,
	}

	//Save user
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful",
	})
}

// Login
func Login(c *gin.Context) {
	var req models.LoginRequest

	//Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Find user
	var user models.User
	err := database.DB.Where("email = ?", req.Email).First(&user).Error
	
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
		return
	}

	//Check password
	if err := utils.CheckPassword(user.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Check active status
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "account disabled",
		})
		return
	}

	//Generate JWT AccesToken
	acces, err := token.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate acess token"})
		return
	}

	refresh, err := token.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate Refresh token"})
		return
	}

	//Respond
	c.JSON(http.StatusOK, gin.H{"sucess": "Login sucessfull",
		"userName ":    user.FirstName,
		"AcessToken":   acces,
		"Refreshtoken": refresh,
	})

}
