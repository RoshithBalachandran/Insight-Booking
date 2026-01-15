package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/insight/internals/authHandler"
	"github.com/insight/internals/middleware"
	"github.com/insight/internals/users"
)

func SetupRouter(r *gin.Engine) {

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.GET("/home", users.HomePage)
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", users.GetProfile)
		protected.PUT("/profile", users.UserProfileUpdate)
		protected.POST("/logout", authHandler.Logout)
	}
}
