package counselor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/models"
	"github.com/insight/internals/utils"
)

func AdminCreateCounselor(c *gin.Context) {
	type CreateCounselorRequest struct {
		Name           string `json:"name" binding:"required,min=2"`
		Qualification  string `json:"qualification" binding:"required"`
		Specialization string `json:"specialization" binding:"required"`
		Email          string `json:"email" binding:"required,email"`
		Password       string `json:"password" binding:"required,min=8"`
	}

	var req CreateCounselorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var count int64
	database.DB.Model(&models.Counselor{}).
		Where("email = ?", req.Email).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		return
	}

	counselor := models.Counselor{
		Name:           req.Name,
		Qualification:  req.Qualification,
		Specialization: req.Specialization,
		Email:          req.Email,
		Password:       hash,
		IsActive:       true,
		Block:          false,
	}

	if err := database.DB.Create(&counselor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "counselor creation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "counselor created successfully"})
}
