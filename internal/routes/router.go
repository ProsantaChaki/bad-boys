package routes

import (
	"bad_boyes/internal/handler"
	"bad_boyes/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authHandler *handler.AuthHandler, postHandler *handler.PostHandler) {
	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// User profile route
		auth.GET("/profile", authHandler.GetProfile)

		// Post routes
		auth.POST("/posts", postHandler.CreatePost)
		auth.GET("/posts", postHandler.ListPosts)
		auth.GET("/posts/:id", postHandler.GetPost)
		auth.PUT("/posts/:id", postHandler.UpdatePost)
		auth.DELETE("/posts/:id", postHandler.DeletePost)
		auth.GET("/posts/:id/history", postHandler.GetPostHistory)

		// Report routes
		auth.POST("/posts/:id/report", postHandler.CreateReport)

		// Admin routes
		admin := auth.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.PUT("/reports/:id/status", postHandler.UpdateReportStatus)
			admin.GET("/reports", postHandler.ListReports)
		}
	}
}
