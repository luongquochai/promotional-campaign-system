package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	database "github.com/luongquochai/promotional-campaign-system/database/mysql"
	"github.com/luongquochai/promotional-campaign-system/models"
)

func CreateCampaign(c *gin.Context, userID uint) (*models.Campaign, error) {
	var campaign models.Campaign

	campaign.CreatorID = userID

	if err := c.ShouldBindJSON(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	if campaign.Discount <= 0 || campaign.Discount > 100 {
		return nil, errors.New("discount_percentage must be between 1 and 100")
	}

	if campaign.StartDate.After(campaign.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date"})
		return nil, errors.New("start date must be before end date")
	}

	// check duplicate campaign
	dupCampaign, err := checkDuplicateCampaign(&campaign)
	if err != nil {
		return nil, err
	}

	fmt.Println(dupCampaign)

	if dupCampaign != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "duplicate campaign"})
		return dupCampaign, errors.New("duplicate campaign")
	}

	if err := database.DB.Create(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
		return nil, err
	}

	return &campaign, nil
}

func ListCampaigns(c *gin.Context, userID uint) ([]*models.Campaign, error) {
	var campaigns []*models.Campaign
	if err := database.DB.Where("creator_id = ?", userID).Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaigns"})
		return nil, err
	}

	return campaigns, nil
}

func GetCampaignID(c *gin.Context, userID uint, campaignID string) (*models.Campaign, error) {
	var campaign models.Campaign
	if err := database.DB.Where("creator_id = ?", userID).First(&campaign, campaignID).Error; err != nil {
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
	if err := database.DB.Where("creator_id = ?", userID).Delete(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign"})
		return err
	}
	return nil
}

func checkDuplicateCampaign(req *models.Campaign) (*models.Campaign, error) {
	var campaignModel []*models.Campaign
	err := database.DB.Where("discount = ? AND start_date <= ? AND end_date >= ?", req.Discount, req.EndDate, req.StartDate).
		Find(&campaignModel).Error
	if err != nil {
		return nil, errors.New("database query error")
	}
	for _, d := range campaignModel {
		if d.Status == "active" {
			return d, nil
		}
	}
	return nil, nil
}

func UpdateCampaignID(c *gin.Context, updateInfo map[string]interface{}) error {
	// Save the updated values to the database
	var campaign models.Campaign
	if err := database.DB.Model(&campaign).Updates(updateInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign"})
		return err
	}
	return nil
}
