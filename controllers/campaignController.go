package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/services"
	"github.com/luongquochai/promotional-campaign-system/utils"
)

//TODO: this file to basic logic that is have account
// We should base on role and authorization for account creation

// CreateCampaign creates a new campaign
// @Summary Create a new campaign
// @Description Create a campaign with user details and return the created campaign
// @Tags Campaign
// @Accept  json
// @Produce  json
// @Param user_id header string true "User ID"
// @Param campaign body models.Campaign true "Campaign data"
// @Success 201 {object} models.Campaign
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /campaigns [post]
func CreateCampaign(c *gin.Context) {
	// Retrieve user_id from context
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondUnauthorized(c, err)
		return
	}

	campaign, err := services.CreateCampaign(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err, "data": campaign})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": campaign})
}

// ListCampaigns returns a list of campaigns
// @Summary List all campaigns
// @Description Get a list of all campaigns created by the user
// @Tags Campaign
// @Accept  json
// @Produce  json
// @Param user_id header string true "User ID"
// @Success 200 {array} models.Campaign
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /campaigns [get]
func ListCampaigns(c *gin.Context) {
	// Retrieve user_id from context
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondUnauthorized(c, err)
		return
	}

	campaigns, err := services.ListCampaigns(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve campaigns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": campaigns})
}

// GetCampaign retrieves a single campaign by its ID
// @Summary Get a campaign by ID
// @Description Retrieve details of a specific campaign by its ID
// @Tags Campaign
// @Accept  json
// @Produce  json
// @Param id path string true "Campaign ID"
// @Success 200 {object} models.Campaign
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /campaigns/{id} [get]
func GetCampaign(c *gin.Context) {
	// Retrieve user_id from context
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondUnauthorized(c, err)
		return
	}

	id := c.Param("id")
	campaign, err := services.GetCampaignID(c, userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": campaign})
}

// UpdateCampaign updates a campaign
// @Summary Update a campaign
// @Description Update the details of a specific campaign
// @Tags Campaign
// @Accept  json
// @Produce  json
// @Param id path string true "Campaign ID"
// @Param campaign body models.Campaign true "Campaign update data"
// @Success 200 {object} models.Campaign
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /campaigns/{id} [put]
func UpdateCampaign(c *gin.Context) {
	// Retrieve user_id from context
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondUnauthorized(c, err)
		return
	}

	id := c.Param("id")

	// Retrieve the existing campaign by ID
	campaign, err := services.GetCampaignID(c, userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	updateInfo, err := services.BuildUpdateInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateCampaignID(c, updateInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the updated campaign in the response
	c.JSON(http.StatusOK, gin.H{"data": campaign})
}

// DeleteCampaign deletes a campaign
// @Summary Delete a campaign
// @Description Delete a specific campaign by its ID
// @Tags Campaign
// @Accept  json
// @Produce  json
// @Param id path string true "Campaign ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /campaigns/{id} [delete]
func DeleteCampaign(c *gin.Context) {
	// Retrieve user_id from context
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondUnauthorized(c, err)
		return
	}

	id := c.Param("id")
	// Check if the campaign exists
	campaign, err := services.GetCampaignID(c, userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	// Delete the campaign from the database
	if err := services.DeleteCampaignID(c, userID, campaign); err != nil {
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}
