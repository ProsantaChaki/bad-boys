package routes

import (
	"bad_boyes/internal/controllers"
	"bad_boyes/internal/middleware"
	"bad_boyes/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	authService := services.NewAuthService()
	authController := controllers.NewAuthController(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), authController.GetProfile)
	}
}
