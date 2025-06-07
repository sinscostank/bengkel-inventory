package middleware

import "github.com/gin-gonic/gin"

// AuthMiddleware stub so route.go compiles.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
