package adminusers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
	"github.com/insight/internals/utils"
)

// Admin CreateUser
func AdminCreateUser(c *gin.Context) {
	type CreateUserRequest struct {
		FirstName   string  `json:"first_name" binding:"required,min=2"`
		LastName    string  `json:"last_name"`
		Email       string  `json:"email" binding:"required,email"`
		PhoneNumber *string `json:"phone_number"`
		Password    string  `json:"password" binding:"required,min=8"`
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var count int64
	database.DB.Model(&models.User{}).
		Where("email = ? OR phone_number = ?", req.Email, req.PhoneNumber).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "email or phone already exists"})
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		return
	}

	user := models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hash,
		Role:        constant.User,
		IsActive:    true,
		Block:       false,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}
