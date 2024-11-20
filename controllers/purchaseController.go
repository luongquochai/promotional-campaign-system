package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/models"
	"github.com/luongquochai/promotional-campaign-system/services"
	"golang.org/x/exp/rand"
)

type CampaignID struct {
	CampaignID uint `json:"campaign_id"`
}

// ProcessPurchase handles processing a discounted subscription purchase
func CreatePurchase(c *gin.Context) {
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

	// Step 1: Retrieve the voucher (e.g., using the campaign_id in the request)
	var campaignID CampaignID
	if err := c.ShouldBindJSON(&campaignID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var voucher models.Voucher
	if err := config.DB.Where("campaign_id = ? AND user_id = ? AND used_at IS NULL", campaignID.CampaignID, userIDInt).First(&voucher).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired voucher"})
		return
	}

	// Step 2: Apply discount
	// For simplicity, let's assume the subscription price is fixed, e.g., $100
	subscriptionPrice := 100.0
	discountAmount := subscriptionPrice * (voucher.Discount / 100)
	finalPrice := subscriptionPrice - discountAmount

	// Step 3: Create the purchase record
	purchase := models.Purchase{
		UserID:          userIDInt,
		TransactionID:   fmt.Sprintf("TXN%06d", rand.Intn(1e6)), // Unique Transaction ID
		SubscriptionID:  voucher.CampaignID,                     // Get the subscription from the voucher
		DiscountApplied: discountAmount,
		FinalPrice:      finalPrice,
		Status:          "completed", // assuming the purchase is successful
	}

	if err := services.CreatePurchase(c, purchase); err != nil {
		return
	}

	// Step 4: Mark voucher as used
	if err := services.UpdateVoucher(c, voucher); err != nil {
		return
	}

	// Step 5: Respond with the purchase details
	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase successful",
		"data":    purchase,
	})
}

func GetPurchaseHistory(c *gin.Context) {
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

	purchases, err := services.GetPurchaseHistory(c, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving purchase history"})
		return
	}

	c.JSON(http.StatusOK, purchases)
}
