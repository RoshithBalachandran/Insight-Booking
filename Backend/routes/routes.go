package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/insight/internals/auth"
)


func SetupRouter(r *gin.Engine){
	r.POST("/reg",auth.Register)
}