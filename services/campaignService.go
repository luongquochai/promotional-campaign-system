package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/models"
)

func CreateCampaign(c *gin.Context, userID uint) (*models.Campaign, error) {
	var campaign models.Campaign

	campaign.CreatorID = userID

	if err := c.ShouldBindJSON(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	if campaign.StartDate.After(campaign.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date"})
		return nil, fmt.Errorf("start date must be before end date")
	}

	if err := config.DB.Create(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
		return nil, err
	}

	return &campaign, nil
}

func ListCampaigns(c *gin.Context, userID uint) ([]*models.Campaign, error) {
	var campaigns []*models.Campaign
	if err := config.DB.Where("creator_id = ?", userID).Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaigns"})
		return nil, err
	}

	return campaigns, nil
}

func GetCampaignID(c *gin.Context, userID uint, campaignID string) (*models.Campaign, error) {
	var campaign models.Campaign
	if err := config.DB.Where("creator_id = ?", userID).First(&campaign, campaignID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return nil, err
	}

	return &campaign, nil
}

func BuildUpdateInfo(c *gin.Context) (map[string]interface{}, error) {
	var input models.Campaign
	// Bind the incoming JSON data to the input variable
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
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

	return updates, nil
}

func DeleteCampaignID(c *gin.Context, userID uint, campaign *models.Campaign) error {
	if err := config.DB.Where("creator_id = ?", userID).Delete(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign"})
		return err
	}
	return nil
}
