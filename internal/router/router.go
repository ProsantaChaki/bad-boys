package router

import (
	"bad_boyes/internal/handler"
	"bad_boyes/internal/middleware"
	"bad_boyes/internal/repository"
	"bad_boyes/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Initialize dependencies
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
	}
}
