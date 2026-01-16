package adminapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/models"
	"github.com/insight/internals/utils"
	"gorm.io/gorm"
)

// Get all counselor
func AdminGetAllCounselors(c *gin.Context) {
	var counselors []models.Counselor

	if err := database.DB.Find(&counselors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch counselors: " + err.Error()})
		return
	}

	// Hide passwords
	for i := range counselors {
		counselors[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Counselors fetched successfully",
		"counselors": counselors,
	})
}

// creates a new counselor (admin-only)
func AdminCreateCounselor(c *gin.Context) {
	var input models.Counselor

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate required fields
	if input.Name == "" || input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, Email, and Password are required"})
		return
	}

	// Force default values
	input.Block = false
	input.IsActive = true

	// Check if email already exists
	var existing models.Counselor
	if err := database.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "A counselor with this email already exists"})
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = hashedPassword

	// Create counselor
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create counselor: " + err.Error()})
		return
	}

	// Do not return password in response
	input.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Counselor created successfully",
		"counselor": input,
	})
}

// Update Counselor
func AdminUpdateCounselor(c *gin.Context) {
	id := c.Param("id")
	var counselor models.Counselor

	// Find counselor
	if err := database.DB.First(&counselor, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Counselor not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	var input struct {
		Name           *string `json:"name"`
		Qualification  *string `json:"qualification"`
		Specialization *string `json:"specialization"`
		Email          *string `json:"email"`
		Password       *string `json:"password"`
		Block          *bool   `json:"block"`
		IsActive       *bool   `json:"is_active"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Check email uniqueness if updated
	if input.Email != nil && *input.Email != counselor.Email {
		var existing models.Counselor
		if err := database.DB.Where("email = ?", *input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		} else if err != nil && err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}
	}

	// Update fields safely
	if input.Name != nil {
		counselor.Name = *input.Name
	}
	if input.Qualification != nil {
		counselor.Qualification = *input.Qualification
	}
	if input.Specialization != nil {
		counselor.Specialization = *input.Specialization
	}
	if input.Email != nil {
		counselor.Email = *input.Email
	}
	if input.Password != nil {
		hashed, err := utils.HashPassword(*input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		counselor.Password = hashed
	}
	if input.Block != nil {
		counselor.Block = *input.Block
	}
	if input.IsActive != nil {
		counselor.IsActive = *input.IsActive
	}

	// Save changes
	if err := database.DB.Save(&counselor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update counselor: " + err.Error()})
		return
	}

	counselor.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message":        "Counselor updated successfully",
		"counselor":      counselor.Name,
		"Qualification":  counselor.Qualification,
		"Specialization": counselor.Specialization,
		"Email":          counselor.Email,
	})
}

//delete counselor

func AdminDeleteCounselor(c *gin.Context) {
	id := c.Param("id")
	var counselor models.Counselor

	// Find counselor
	if err := database.DB.First(&counselor, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Counselor not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Soft delete counselor
	if err := database.DB.Delete(&counselor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete counselor: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Counselor deleted successfully"})
}
