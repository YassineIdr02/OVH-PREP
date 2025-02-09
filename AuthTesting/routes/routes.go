package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/YassineIdr02/ovh-prep/E2E-Tests/controllers"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/register", controllers.RegisterUser)
		api.POST("/login", controllers.LoginUser)
	}
}
