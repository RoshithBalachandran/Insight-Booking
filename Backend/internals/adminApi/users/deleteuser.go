package adminusers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
)

func AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.
		Where("id = ? AND role = ?", id, constant.User).
		First(&user).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
