package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/models"
)

func CreatePurchase(c *gin.Context, purchase models.Purchase) error {
	if err := config.DB.Create(&purchase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process purchase"})
		return err
	}
	return nil
}

func GetPurchaseHistory(c *gin.Context, userID uint) (*[]models.Purchase, error) {
	var purchases []models.Purchase
	if err := config.DB.Where("user_id = ?", userID).Find(&purchases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve purchases"})
		return nil, err
	}

	return &purchases, nil
}
