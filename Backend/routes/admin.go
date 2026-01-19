package routes

import (
	"github.com/gin-gonic/gin"
	adminhandlers "github.com/insight/internals/AdminHandlers"
	counselor "github.com/insight/internals/adminApi/Counselor"
	adminusers "github.com/insight/internals/adminApi/users"
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
		admin.PATCH("/profile", adminhandlers.UpdateAdminProfile)

		// User endpoints
		admin.GET("/users", adminusers.AdminGetAllUsers)
		admin.POST("/users", adminusers.AdminCreateUser)
		admin.PUT("/users/:id", adminusers.AdminUpdateUser)

		// Counselor endpoints
		admin.GET("/counselors", counselor.AdminGetAllCounselors)
		admin.POST("/counselors", counselor.AdminCreateCounselor)
		admin.PUT("/counselors/:id", counselor.AdminUpdateCounselor)

		admin.POST("/logout", authHandler.Logout)
	}
}
