package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	database "github.com/luongquochai/promotional-campaign-system/database/mysql"
	"github.com/luongquochai/promotional-campaign-system/models"
)

func CreatePurchase(c *gin.Context, purchase models.Purchase) error {
	if err := database.DB.Create(&purchase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process purchase"})
		return err
	}
	return nil
}

func GetPurchaseHistory(c *gin.Context, userID uint) (*[]models.Purchase, error) {
	var purchases []models.Purchase
	if err := database.DB.Where("user_id = ?", userID).Find(&purchases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve purchases"})
		return nil, err
	}

	return &purchases, nil
}

func CheckValidVoucher(c *gin.Context, userID, campaignID uint) (*models.Voucher, error) {
	var voucher models.Voucher
	if err := database.DB.Where("campaign_id = ? AND user_id = ? AND used_at IS NULL", campaignID, userID).First(&voucher).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired voucher"})
		return nil, err
	}

	return &voucher, nil
}
