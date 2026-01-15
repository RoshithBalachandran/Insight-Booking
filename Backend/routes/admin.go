package routes

import (
	"github.com/gin-gonic/gin"
	adminhandlers "github.com/insight/internals/AdminHandlers"
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



		admin.POST("/logout", authHandler.Logout)
	}
}
