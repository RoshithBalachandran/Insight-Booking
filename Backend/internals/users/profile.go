package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
	"github.com/insight/internals/utils"
)

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.Select("id, first_name, last_name, email, role, is_active, created_at").
		First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is inactive"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"role":       user.Role,
		"joined_at":  user.CreatedAt,
	})
}

// UserProfileUpdate updates profile info for logged-in users
func UserProfileUpdate(c *gin.Context) {
	// Get logged-in user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}

	// Bind input JSON
	var input struct {
		FirstName   string  `json:"first_name" binding:"required"`
		LastName    string  `json:"last_name" binding:"required"`
		PhoneNumber *string `json:"phone_number"`
		Email       string  `json:"email" binding:"required,email"`
		Password    string  `json:"password"` // optional, for password change
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// Fetch user from DB
	var user models.User
	if err := database.DB.Where("id = ? AND role = ?", userID, constant.User).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Prevent updates if user is inactive
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is inactive"})
		return
	}

	// Check email uniqueness if changed
	if user.Email != input.Email {
		var existing models.User
		if err := database.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
			return
		}
	}

	// Update fields
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Email = input.Email
	user.PhoneNumber = input.PhoneNumber

	// Optional password update
	if input.Password != "" {
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		user.Password = hashedPassword
	}

	// Save updates
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "profile updated successfully",
		"id":           user.ID,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"role":         user.Role,
	})
}
