package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sinscostank/bengkel-inventory/utils"
)

// AuthMiddleware stub so route.go compiles.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the "Authorization" header value
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Split the string into "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Extract the token
		token := parts[1]

		// Validate the token (this should return the user info if the token is valid)
		userClaims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Attach the user claims to the context (for later use in controllers)
		c.Set("userClaims", userClaims)

		// Continue to the next handler
		c.Next()
	}
}

// AdminMiddleware is a middleware to ensure that only admin users can access a route.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve user claims from the context
		userClaims, exists := c.Get("userClaims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User claims not found"})
			c.Abort()
			return
		}

		// Type assert the userClaims to your user model (make sure it's the correct type)
		claims, ok := userClaims.(*utils.UserClaims) // assuming UserClaims is the structure that contains role information
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user claims"})
			c.Abort()
			return
		}

		// Check if the user's role is "admin"
		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required admin role"})
			c.Abort()
			return
		}

		// Continue to the next handler
		c.Next()
	}
}