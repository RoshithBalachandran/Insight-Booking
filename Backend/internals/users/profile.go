package users

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
	"github.com/insight/internals/utils"
)

func GetProfile(c *gin.Context) {
	//Get authenticated user ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//Response DTO 
	type UserProfileResponse struct {
		ID        uint      `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		IsActive  bool      `json:"is_active"`
		JoinedAt  time.Time `json:"joined_at"`
	}

	var profile UserProfileResponse

	// Fetch only USER role and allowed fields
	err := database.DB.
		Model(&models.User{}).
		Select("id, first_name, last_name, email, role, is_active, created_at").
		Where("id = ? AND role = ?", userID, constant.User).
		Scan(&profile).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	//Account state checks
	if !profile.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "user account is inactive"})
		return
	}

	//Success response
	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}

// UserProfileUpdate updates profile info for logged-in users
func UserProfileUpdate(c *gin.Context) {
	//Get authenticated user ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//Input DTO
	type UpdateUserProfileInput struct {
		FirstName   *string `json:"first_name" binding:"omitempty,min=2,max=100"`
		LastName    *string `json:"last_name"  binding:"omitempty,min=2,max=100"`
		Email       *string `json:"email"      binding:"omitempty,email"`
		PhoneNumber *string `json:"phone_number"`
		Password    *string `json:"password"   binding:"omitempty,min=8"`
	}

	var input UpdateUserProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	//Fetch user
	var user models.User
	if err := database.DB.
		Where("id = ? AND role = ?", userID, constant.User).
		First(&user).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	//Account state checks
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "user account is inactive"})
		return
	}
	if user.Block {
		c.JSON(http.StatusForbidden, gin.H{"error": "user account is blocked"})
		return
	}

	//Email uniqueness check
	if input.Email != nil && *input.Email != user.Email {
		var count int64
		if err := database.DB.
			Model(&models.User{}).
			Where("email = ?", *input.Email).
			Count(&count).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "email validation failed"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
	}

	//Build updates map
	updates := map[string]interface{}{}

	if input.FirstName != nil {
		updates["first_name"] = *input.FirstName
	}
	if input.LastName != nil {
		updates["last_name"] = *input.LastName
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.PhoneNumber != nil {
		updates["phone_number"] = input.PhoneNumber
	}

	// Password update (hashed)
	if input.Password != nil {
		hashedPassword, err := utils.HashPassword(*input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		updates["password"] = hashedPassword
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided for update"})
		return
	}

	//Update DB
	if err := database.DB.
		Model(&user).
		Updates(updates).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	// 8. Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"data": gin.H{
			"id":           user.ID,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"role":         user.Role,
		},
	})
}
