package routes

import (
	"github.com/gin-gonic/gin"
	adminhandlers "github.com/insight/internals/AdminHandlers"
	adminapi "github.com/insight/internals/adminApi"
	"github.com/insight/internals/authHandler"
	"github.com/insight/internals/middleware"
)

func SetUpAdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")

	// Public route
	admin.POST("/login", adminhandlers.AdminLogin)

	// Protected admin routes
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.GET("/profile", adminhandlers.GetAdminProfile)
		admin.PUT("/profile", adminhandlers.UpdateAdminProfile)

		// User endpoints
		admin.GET("/users", adminapi.AdminGetAllUsers)
		admin.POST("/users", adminapi.AdminCreateUser)
		admin.PUT("/users/:id", adminapi.AdminUpdateUser)

		// Counselor endpoints
		admin.GET("/counselors", adminapi.AdminGetAllCounselors)
		admin.POST("/counselors", adminapi.AdminCreateCounselor)
		admin.PUT("/counselors/:id", adminapi.AdminUpdateCounselor)

		admin.POST("/logout", authHandler.Logout)
	}
}
