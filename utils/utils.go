package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserID extracts the user_id from the Gin context
func GetUserID(c *gin.Context) (uint, error) {
	// Retrieve user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user not authenticated")
	}

	// Type assertion to uint
	userIDInt, ok := userID.(uint)
	if !ok {
		return 0, errors.New("invalid user_id type")
	}

	return userIDInt, nil
}

// RespondUnauthorized sends a standard unauthorized response
func RespondUnauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
}

// ErrorResponse represents the structure of an error response
// swagger:response ErrorResponse
type ErrorResponse struct {
	Message string `json:"message"` // Message describes the error
	Code    int    `json:"code"`    // Code represents the error code
}

// SuccessResponse represents the structure of a successful response
// swagger:response SuccessResponse
type SuccessResponse struct {
	Message string `json:"message"` // A message indicating the success
	Data    any    `json:"data"`    // The data returned with the response
}
