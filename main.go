package main

import (
	"bad_boyes/internal/handler"
	"bad_boyes/internal/middleware"
	"bad_boyes/internal/repository"
	"bad_boyes/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/bad_boyes?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, auditRepo)
	postService := services.NewPostService(postRepo, auditRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)

	// Initialize router
	r := gin.Default()

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
