package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/models"
)

func HomePage(c *gin.Context) {
	var services []models.Service
	var counselors []models.Counselor

	database.DB.
		Where("is_active = true").
		Find(&services)

	database.DB.
		Where("is_active = true").
		Find(&counselors)

	c.JSON(http.StatusOK, gin.H{
		"center": gin.H{
			"name":    "Insight Counseling Center",
			"tagline": "Your mind matters",
			"phone":   "8304852769",
			"address": "13th Mile, Neerveli, Kannur 670701",
		},
		"services":   services,
		"counselors": counselors,
	})
}
