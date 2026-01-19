package adminhandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
)

// GetAdminProfile returns profile info for admins only
func GetAdminProfile(c *gin.Context) {
	//Get authenticated admin ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//Response DTO
	type AdminProfileResponse struct {
		ID        uint      `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		IsActive  bool      `json:"is_active"`
		JoinedAt  time.Time `json:"joined_at"`
	}

	var profile AdminProfileResponse

	//Fetch admin safely
	err := database.DB.
		Model(&models.User{}).
		Select("id, first_name, last_name, email, role, is_active, created_at").
		Where("id = ? AND role = ?", userID, constant.Admin).
		Scan(&profile).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	//Account state checks
	if !profile.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin account is inactive"})
		return
	}

	//Success response
	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}

func UpdateAdminProfile(c *gin.Context) {
	// 1. Get logged-in user ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	//Request DTO
	type UpdateAdminProfileInput struct {
		FirstName *string `json:"first_name" binding:"omitempty,min=2,max=100"`
		LastName  *string `json:"last_name"  binding:"omitempty,min=2,max=100"`
		Email     *string `json:"email"      binding:"omitempty,email"`
	}

	var input UpdateAdminProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	//Fetch admin from DB
	var admin models.User
	if err := database.DB.
		Where("id = ? AND role = ?", userID, constant.Admin).
		First(&admin).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	//Status checks
	if !admin.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin account is inactive"})
		return
	}
	if admin.Block {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin account is blocked"})
		return
	}

	//Email uniqueness validation (only if changed)
	if input.Email != nil && *input.Email != admin.Email {
		var count int64
		if err := database.DB.
			Model(&models.User{}).
			Where("email = ?", *input.Email).
			Count(&count).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
	}

	//Apply only provided fields
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

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided for update"})
		return
	}

	//Update DB
	if err := database.DB.
		Model(&admin).
		Updates(updates).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	//Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"data": gin.H{
			"id":         admin.ID,
			"first_name": admin.FirstName,
			"last_name":  admin.LastName,
			"email":      admin.Email,
			"role":       admin.Role,
		},
	})
}
