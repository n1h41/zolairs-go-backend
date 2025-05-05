package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// GinAuthMiddleware checks for user authentication
// This is a simplified version - in a real app, use JWT or OAuth2
func GinAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from header
		userID := c.GetHeader("X-User-ID")

		// For demo purposes only - in production, never use a default user
		if userID == "" {
			c.JSON(401, gin.H{"status": false, "message": "Unauthorized: Missing user ID"})
			c.Abort()
			return
		}

		// Log authentication
		log.Printf("Authenticated request for user: %s", userID)

		// Add user ID to request context
		c.Set(string(UserIDKey), userID)

		c.Next()
	}
}

// GinLoggerMiddleware logs request details
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Log request details
		log.Printf("%s | %s | %d | %s | %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
		)
	}
}

// GetUserIDFromGin extracts the user ID from the Gin context
func GetUserIDFromGin(c *gin.Context) string {
	userID, exists := c.Get(string(UserIDKey))
	if !exists {
		return ""
	}
	return userID.(string)
}
