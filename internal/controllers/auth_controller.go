package controllers

import (
	"bad_boyes/internal/models"
	"bad_boyes/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid registration request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	if err := c.authService.Register(req); err != nil {
		switch err {
		case services.ErrUserAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{
				"code":    409,
				"status":  "error",
				"message": "User already exists",
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
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"status":  "error",
				"message": "Registration failed",
				"error":   err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"status":  "success",
		"message": "User registered successfully",
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid login request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	response, err := c.authService.Login(req)
	if err != nil {
		switch err {
		case services.ErrUserNotFound, services.ErrInvalidPassword:
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"status":  "error",
				"message": "Invalid credentials",
				"error":   err.Error(),
			})
		case services.ErrJWTSecretMissing:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"status":  "error",
				"message": "Server configuration error",
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
				"message": "Login failed",
				"error":   err.Error(),
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Login successful",
		"data": gin.H{
			"token":      response.Token,
			"expires_at": response.ExpiresAt,
			"user_id":    response.UserID,
			"username":   response.Username,
			"email":      response.Email,
			"name":       response.Name,
			"role":       response.Role,
		},
	})
}

func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"status":  "error",
			"message": "Unauthorized",
			"error":   "user ID not found in context",
		})
		return
	}

	profile, err := c.authService.GetUserProfile(userID)
	if err != nil {
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

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"status":  "success",
		"message": "Profile retrieved successfully",
		"data":    profile,
	})
}
