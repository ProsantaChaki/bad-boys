package handler

import (
	"bad_boyes/internal/models"
	"bad_boyes/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	log.Printf("Received registration request from IP: %s", ctx.ClientIP())

	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid registration request from IP %s: %v", ctx.ClientIP(), err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Processing registration for email: %s", req.Email)
	if err := h.authService.Register(req); err != nil {
		log.Printf("Registration failed for email %s: %v", req.Email, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Registration failed",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Registration successful for email: %s", req.Email)
	ctx.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"status":  "success",
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	log.Printf("Received login request from IP: %s", ctx.ClientIP())

	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid login request from IP %s: %v", ctx.ClientIP(), err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Processing login for email: %s", req.Email)
	response, err := h.authService.Login(req)
	if err != nil {
		log.Printf("Login failed for email %s: %v", req.Email, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"status":  "error",
			"message": "Login failed",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Login successful for email: %s", req.Email)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Login successful",
		"data":    response,
	})
}

func (h *AuthHandler) GetProfile(ctx *gin.Context) {
	log.Printf("Received profile request from IP: %s", ctx.ClientIP())

	userID := ctx.GetUint("user_id")
	log.Printf("Processing profile request for user ID: %d", userID)

	if userID == 0 {
		log.Printf("Unauthorized profile access attempt from IP: %s", ctx.ClientIP())
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"status":  "error",
			"message": "Unauthorized",
			"error":   "user ID not found in context",
		})
		return
	}

	profile, err := h.authService.GetUserProfile(userID)
	if err != nil {
		log.Printf("Failed to fetch profile for user ID %d: %v", userID, err)
		switch err {
		case services.ErrUserNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"status":  "error",
				"message": "User not found",
				"error":   err.Error(),
			})
		case services.ErrDatabaseError:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"status":  "error",
				"message": "Database error occurred",
				"error":   err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"status":  "error",
				"message": "Failed to fetch profile",
				"error":   err.Error(),
			})
		}
		return
	}

	log.Printf("Profile fetched successfully for user ID: %d", userID)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Profile retrieved successfully",
		"data":    profile,
	})
}
