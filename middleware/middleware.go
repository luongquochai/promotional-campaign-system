package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/database/redis"
)

// AuthMiddleware is a middleware function that validates the JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// The Authorization header should be in the form of "Bearer <token>"
		// Split the string into "Bearer" and the actual token
		tokenString := ""
		if len(strings.Split(authHeader, " ")) == 2 {
			tokenString = strings.Split(authHeader, " ")[1]
		}

		// If token string is empty, return unauthorized error
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			c.Abort()
			return
		}

		// Validate the JWT token
		_, err := config.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// // Extract user_id from the claims
		// userID, ok := claims["user_id"].(float64) // jwt.MapClaims uses float64 for numeric values
		// if !ok {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		// 	c.Abort()
		// 	return
		// }

		// Check Redis for user_id
		userID, err := redis.Get(tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}

		uintUserID, err := strconv.ParseUint(userID, 10, 64) // Base 10, 64-bit width
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Set the user_id in the Gin context for further use
		c.Set("user_id", uint(uintUserID))

		// Proceed to the next handler
		c.Next()
	}
}
