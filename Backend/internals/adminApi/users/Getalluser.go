package adminusers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
)

func AdminGetAllUsers(c *gin.Context) {
	type UserListResponse struct {
		ID          uint      `json:"id"`
		FirstName   string    `json:"first_name"`
		LastName    string    `json:"last_name"`
		Email       string    `json:"email"`
		PhoneNumber *string   `json:"phone_number"`
		IsActive    bool      `json:"is_active"`
		Block       bool      `json:"block"`
		JoinedAt    time.Time `json:"joined_at"`
	}

	var users []UserListResponse

	if err := database.DB.
		Model(&models.User{}).
		Select("id, first_name, last_name, email, phone_number, is_active, block, created_at").
		Where("role = ?", constant.User).
		Order("created_at DESC").
		Scan(&users).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}
