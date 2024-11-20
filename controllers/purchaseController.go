package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/models"
	"github.com/luongquochai/promotional-campaign-system/services"
	"github.com/luongquochai/promotional-campaign-system/utils"
	"golang.org/x/exp/rand"
)

type CampaignID struct {
	CampaignID uint `json:"campaign_id"`
}

// ProcessPurchase handles processing a discounted subscription purchase
// @Summary Process a discounted subscription purchase
// @Description Process a purchase with a valid voucher and apply discount
// @Tags Purchase
// @Accept  json
// @Produce  json
// @Param user_id header string true "User ID"
// @Param campaign_id body CampaignID true "Campaign ID for the purchase"
// @Success 200 {object} models.Purchase
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /purchase/create [post]
func CreatePurchase(c *gin.Context) {
	// Retrieve user_id from context
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondUnauthorized(c, err)
		return
	}

	// Step 1: Retrieve the voucher (e.g., using the campaign_id in the request)
	var campaignID CampaignID
	if err := c.ShouldBindJSON(&campaignID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	voucher, err := services.CheckValidVoucher(c, userID, campaignID.CampaignID)
	if err != nil {
		return
	}

	// Step 2: Apply discount
	// For simplicity, let's assume the subscription price is fixed, e.g., $100
	subscriptionPrice := 100.0
	discountAmount := subscriptionPrice * (voucher.Discount / 100)
	finalPrice := subscriptionPrice - discountAmount

	// Step 3: Create the purchase record
	purchase := models.Purchase{
		UserID:          userID,
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

// GetPurchaseHistory retrieves the purchase history of a user
// @Summary Get purchase history
// @Description Retrieve all purchase history for a specific user
// @Tags Purchase
// @Accept  json
// @Produce  json
// @Param user_id header string true "User ID"
// @Success 200 {array} models.Purchase
// @Failure 500 {object} utils.ErrorResponse
// @Router /purchase/history [get]
func GetPurchaseHistory(c *gin.Context) {
	// Retrieve user_id from context
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondUnauthorized(c, err)
		return
	}

	purchases, err := services.GetPurchaseHistory(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving purchase history"})
		return
	}

	c.JSON(http.StatusOK, purchases)
}
