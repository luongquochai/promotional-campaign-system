package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/models"
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
	var campaign models.Campaign
	var creator models.User

	campaign.CreatorID = userIDInt
	campaign.Creator = creator

	if err := c.ShouldBindJSON(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if campaign.StartDate.After(campaign.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date"})
		return
	}

	if err := config.DB.Create(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
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

	var campaigns []models.Campaign
	if err := config.DB.Where("creator_id = ?", userIDInt).Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaigns"})
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
	var campaign models.Campaign
	if err := config.DB.Where("creator_id = ?", userIDInt).First(&campaign, id).Error; err != nil {
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
	var campaign models.Campaign

	// Retrieve the existing campaign by ID
	if err := config.DB.Where("creator_id = ?", userIDInt).First(&campaign, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	var input models.Campaign
	// Bind the incoming JSON data to the input variable
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepare a map to store the updated fields
	updates := make(map[string]interface{})
	// Update only the fields that are passed in the request
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Discount != 0 {
		updates["discount"] = input.Discount
	}
	if !input.StartDate.IsZero() {
		updates["start_date"] = input.StartDate
	}
	if !input.EndDate.IsZero() {
		updates["end_date"] = input.EndDate
	}
	if input.MaxUsers != 0 {
		updates["max_users"] = input.MaxUsers
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	updates["updated_at"] = time.Now()

	// Save the updated values to the database
	if err := config.DB.Model(&campaign).Updates(updates).Error; err != nil {
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
	var campaign models.Campaign

	// Check if the campaign exists
	if err := config.DB.Where("creator_id = ?", userIDInt).First(&campaign, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	// Delete the campaign from the database
	if err := config.DB.Where("creator_id = ?", userIDInt).Delete(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign"})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}
