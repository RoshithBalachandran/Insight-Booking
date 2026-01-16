package adminapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/models"
	"github.com/insight/internals/utils"
	"gorm.io/gorm"
)

//get all user
func AdminGetAllUsers(c *gin.Context) {
	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users: " + err.Error()})
		return
	}

	// Hide passwords
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users fetched successfully",
		"users":   users,
	})
}

// Admin CreateUser
func AdminCreateUser(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate required fields
	if input.Email == "" || input.Password == "" || input.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, Password, and FirstName are required"})
		return
	}

	// Force role to USER
	input.Role = "USER"

	// Check for existing user by email or phone
	var existing models.User
	if err := database.DB.Where("email = ? OR phone_number = ?", input.Email, input.PhoneNumber).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email or phone number already exists"})
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Hash password
	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = hash

	// Create user
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": input})
}


//admin update users
func AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Find the user
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Bind input
	var input struct {
		FirstName   *string `json:"first_name"`
		LastName    *string `json:"last_name"`
		PhoneNumber *string `json:"phone_number"`
		Email       *string `json:"email"`
		Password    *string `json:"password"`
		IsActive    *bool   `json:"is_active"`
		Block       *bool   `json:"block"` // Admin can block/unblock
		// Role cannot be changed
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Check email uniqueness if updated
	if input.Email != nil && *input.Email != user.Email {
		var existing models.User
		if err := database.DB.Where("email = ?", *input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		} else if err != nil && err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}
	}

	// Check phone uniqueness if updated
	if input.PhoneNumber != nil && (user.PhoneNumber == nil || *input.PhoneNumber != *user.PhoneNumber) {
		var existing models.User
		if err := database.DB.Where("phone_number = ?", *input.PhoneNumber).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already in use"})
			return
		} else if err != nil && err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}
	}

	// Update fields safely
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.PhoneNumber != nil {
		user.PhoneNumber = input.PhoneNumber
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		hash, err := utils.HashPassword(*input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hash
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}
	if input.Block != nil {
		user.Block = *input.Block
	}

	// Save updated user
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	// Hide password in response
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

//Delete user
func AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Find user
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Soft delete user
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}


