package counselor

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/models"
)

func AdminGetAllCounselors(c *gin.Context) {
	type CounselorListResponse struct {
		ID             uint      `json:"id"`
		Name           string    `json:"name"`
		Qualification  string    `json:"qualification"`
		Specialization string    `json:"specialization"`
		Email          string    `json:"email"`
		IsActive       bool      `json:"is_active"`
		Block          bool      `json:"block"`
		JoinedAt       time.Time `json:"joined_at"`
	}

	var counselors []CounselorListResponse

	if err := database.DB.
		Model(&models.Counselor{}).
		Select("id, name, qualification, specialization, email, is_active, block, created_at").
		Order("created_at DESC").
		Scan(&counselors).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch counselors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": counselors})
}
