package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("Processing request: %s %s", ctx.Request.Method, ctx.Request.URL.Path)

		// Get the Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("No Authorization header found")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"status":  "error",
				"message": "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Invalid Authorization header format")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"status":  "error",
				"message": "Invalid Authorization header format",
			})
			ctx.Abort()
			return
		}

		// Get the token
		tokenString := parts[1]

		// Get JWT secret from environment
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Printf("JWT_SECRET environment variable is not set")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"status":  "error",
				"message": "Server configuration error",
			})
			ctx.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			log.Printf("Token validation failed: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"status":  "error",
				"message": "Invalid token",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			log.Printf("Token is invalid")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"status":  "error",
				"message": "Invalid token",
			})
			ctx.Abort()
			return
		}

		// Get claims from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Failed to get claims from token")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"status":  "error",
				"message": "Invalid token claims",
			})
			ctx.Abort()
			return
		}

		// Get user ID from claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Printf("Failed to get user_id from claims")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"status":  "error",
				"message": "Invalid token claims",
			})
			ctx.Abort()
			return
		}

		// Set user ID in context
		ctx.Set("user_id", uint(userID))
		ctx.Set("claims", claims)
		log.Printf("User authenticated: ID=%d", uint(userID))

		ctx.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"status":  "error",
				"message": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		role, ok := claims.(jwt.MapClaims)["role"].(string)
		if !ok || role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"status":  "error",
				"message": "Admin access required",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
