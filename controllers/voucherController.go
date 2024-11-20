package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/services"
	"github.com/luongquochai/promotional-campaign-system/utils"
)

type CodeVoucher struct {
	Code string `json:"code" binding:"required"`
}

// GenerateVoucher generates a new voucher for a user in the given campaign.
func GenerateVoucher(c *gin.Context) {
	// Retrieve user_id from context
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondUnauthorized(c, err)
		return
	}

	// Get the campaign ID from the URL parameter
	var campaignID CampaignID
	if err := c.ShouldBindJSON(&campaignID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate the voucher
	voucher, err := services.GenerateVoucher(userID, campaignID.CampaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the voucher data
	c.JSON(http.StatusOK, gin.H{"voucher": voucher})
}

func ValidateVoucher(c *gin.Context) {
	// Retrieve user_id from context
	// Retrieve user_id from context
	userID, err := utils.GetUserID(c)
	if err != nil {
		utils.RespondUnauthorized(c, err)
		return
	}

	var codeVocher CodeVoucher
	if err := c.ShouldBindJSON(&codeVocher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Using services to validate Voucher and get VoucherCampaign
	voucherCampaign, err := services.ValidateVoucher(userID, codeVocher.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var voucherUsed bool
	if voucherCampaign.Voucher.UsedAt != nil {
		voucherUsed = true
	}

	c.JSON(http.StatusOK, gin.H{
		"is_used":       voucherUsed,
		"campaign_id":   voucherCampaign.Voucher.CampaignID,
		"campaign_name": voucherCampaign.Campaign.Name,
		"discount_rate": voucherCampaign.Campaign.Discount, // Assume static discount for simplicity
	})

}
