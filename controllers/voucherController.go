package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/models"
	"github.com/luongquochai/promotional-campaign-system/services"
	"github.com/luongquochai/promotional-campaign-system/utils"
)

type CodeVoucher struct {
	Code string `json:"code" binding:"required"`
}

// GenerateVoucher generates a new voucher for a user in the given campaign.
// @Summary Generate a new voucher
// @Description Generate a voucher for the user in the specified campaign
// @Tags Voucher
// @Accept  json
// @Produce  json
// @Param campaign_id body controllers.CampaignID true "Campaign ID"
// @Success 200 {object} models.VoucherResponse "Voucher generated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /voucher/generate [post]
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

// ValidateVoucher validates a voucher and checks if it has been used.
// @Summary Validate a voucher
// @Description Validate a voucher code and return voucher usage information
// @Tags Voucher
// @Accept  json
// @Produce  json
// @Param code body controllers.CodeVoucher true "Voucher Code"
// @Success 200 {object} models.VoucherValidationResponse "Voucher validated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /voucher/validate [post]
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

	response := models.VoucherValidationResponse{
		IsUsed:       voucherUsed,
		CampaignID:   voucherCampaign.Voucher.CampaignID,
		CampaignName: voucherCampaign.Campaign.Name,
		DiscountRate: voucherCampaign.Campaign.Discount,
	}

	c.JSON(http.StatusOK, response)

}
