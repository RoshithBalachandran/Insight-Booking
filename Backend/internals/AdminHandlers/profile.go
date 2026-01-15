package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
)

// GetAdminProfile returns profile info for admins only
func GetAdminProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var admin models.User
	err := database.DB.Select("id, first_name, last_name, email, role, is_active, created_at").
		Where("id = ? AND role = ?", userID, constant.Admin).
		First(&admin).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	if !admin.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin is inactive"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         admin.ID,
		"first_name": admin.FirstName,
		"last_name":  admin.LastName,
		"email":      admin.Email,
		"role":       admin.Role,
		"joined_at":  admin.CreatedAt,
	})
}

func UpdateAdminProfile(c *gin.Context) {
	//Get logged user id
	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised"})
		return
	}
	//bind json input
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// fetch data from db
	var admin models.User
	err := database.DB.Where("id=? AND Role=?", userId, constant.Admin).First(&admin).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin details not found"})
		return
	}

	//prevent admin is inactive
	if !admin.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin is inactive"})
		return
	}

	//update fields
	admin.FirstName = input.FirstName
	admin.LastName = input.LastName
	admin.Email = input.Email

	//save data into the database
	if err = database.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sucess": "Profile updated sucessfully",
		"id":         admin.ID,
		"first_name": admin.FirstName,
		"last_name":  admin.LastName,
		"email":      admin.Email,
		"role":       admin.Role,
	})
}
