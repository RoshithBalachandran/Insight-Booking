package counselor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/models"
	"github.com/insight/internals/utils"
)

func AdminUpdateCounselor(c *gin.Context) {
	id := c.Param("id")

	var counselor models.Counselor
	if err := database.DB.First(&counselor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "counselor not found"})
		return
	}

	type UpdateCounselorRequest struct {
		Name           *string `json:"name"`
		Qualification  *string `json:"qualification"`
		Specialization *string `json:"specialization"`
		Email          *string `json:"email"`
		Password       *string `json:"password"`
		Block          *bool   `json:"block"`
		IsActive       *bool   `json:"is_active"`
	}

	var req UpdateCounselorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Qualification != nil {
		updates["qualification"] = *req.Qualification
	}
	if req.Specialization != nil {
		updates["specialization"] = *req.Specialization
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Password != nil {
		hash, _ := utils.HashPassword(*req.Password)
		updates["password"] = hash
	}
	if req.Block != nil {
		updates["block"] = *req.Block
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := database.DB.Model(&counselor).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "counselor updated successfully"})
}
