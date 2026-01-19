package adminusers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
	"github.com/insight/internals/utils"
)

func AdminUpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.
		Where("id = ? AND role = ?", userID, constant.User).
		First(&user).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	type UpdateUserRequest struct {
		FirstName   *string `json:"first_name"`
		LastName    *string `json:"last_name"`
		Email       *string `json:"email"`
		PhoneNumber *string `json:"phone_number"`
		Password    *string `json:"password"`
		IsActive    *bool   `json:"is_active"`
		Block       *bool   `json:"block"`
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	updates := map[string]interface{}{}

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = req.PhoneNumber
	}
	if req.Password != nil {
		hash, _ := utils.HashPassword(*req.Password)
		updates["password"] = hash
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.Block != nil {
		updates["block"] = *req.Block
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
