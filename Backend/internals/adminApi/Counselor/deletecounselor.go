package counselor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/models"
)

func AdminDeleteCounselor(c *gin.Context) {
	id := c.Param("id")

	var counselor models.Counselor
	if err := database.DB.First(&counselor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "counselor not found"})
		return
	}

	if err := database.DB.Delete(&counselor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "counselor deleted successfully"})
}
