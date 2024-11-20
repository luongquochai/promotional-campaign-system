package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/services"
)

//TODO: this file to basic logic that is have account
// We should base on role and authorization for account creation

// Create a new campaign
func CreateCampaign(c *gin.Context) {
	// Retrieve user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	// Type assertion to uint (or whatever type user_id is in your DB)
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id type"})
		return
	}

	campaign, err := services.CreateCampaign(c, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err, "data": campaign})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": campaign})
}

// List all campaigns
func ListCampaigns(c *gin.Context) {
	// Retrieve user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	// Type assertion to uint (or whatever type user_id is in your DB)
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id type"})
		return
	}

	campaigns, err := services.ListCampaigns(c, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve campaigns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": campaigns})
}

// Get campaign details by ID
func GetCampaign(c *gin.Context) {
	// Retrieve user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	// Type assertion to uint (or whatever type user_id is in your DB)
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id type"})
		return
	}

	id := c.Param("id")
	campaign, err := services.GetCampaignID(c, userIDInt, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": campaign})
}

// Update a campaign
func UpdateCampaign(c *gin.Context) {
	// Retrieve user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	// Type assertion to uint (or whatever type user_id is in your DB)
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id type"})
		return
	}

	id := c.Param("id")

	// Retrieve the existing campaign by ID
	campaign, err := services.GetCampaignID(c, userIDInt, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	updateInfo, err := services.BuildUpdateInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated values to the database
	if err := config.DB.Model(&campaign).Updates(updateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign"})
		return
	}

	// Return the updated campaign in the response
	c.JSON(http.StatusOK, gin.H{"data": campaign})
}

func DeleteCampaign(c *gin.Context) {
	// Retrieve user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	// Type assertion to uint (or whatever type user_id is in your DB)
	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id type"})
		return
	}

	id := c.Param("id")
	// Check if the campaign exists
	campaign, err := services.GetCampaignID(c, userIDInt, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	// Delete the campaign from the database
	if err := services.DeleteCampaignID(c, userIDInt, campaign); err != nil {
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}
