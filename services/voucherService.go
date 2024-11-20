package services

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luongquochai/promotional-campaign-system/config"
	"github.com/luongquochai/promotional-campaign-system/models"
)

// GenerateVoucher generates a new voucher for a user in the given campaign.
func GenerateVoucher(userID uint, campaignID uint) (*models.Voucher, error) {
	var campaign models.Campaign

	// Fetch the campaign details
	if err := config.DB.First(&campaign, campaignID).Error; err != nil {
		return nil, errors.New("campaign not found")
	}

	// Campaign validity check
	if time.Now().Before(campaign.StartDate) {
		// If the current time is before the start date, the campaign is not active yet
		return nil, errors.New("campaign is not active")
	}

	if time.Now().After(campaign.EndDate) {
		// If the current time is after the end date, the campaign has expired
		return nil, errors.New("campaign has expired")
	}

	// Check if the campaign voucher is generated for the user
	var voucher models.Voucher
	if err := config.DB.Where("user_id = ? AND campaign_id = ?", userID, campaignID).Find(&voucher).Error; err != nil {
		return nil, errors.New("failed to check voucher campaign")
	}
	fmt.Println(voucher)
	if voucher.Code != "" {
		return nil, fmt.Errorf("voucher for campaign id %d of user id %d is generated", campaignID, userID)
	}

	// Check if the maximum number of vouchers has been reached
	var voucherCount int64
	if err := config.DB.Model(&voucher).Where("campaign_id = ?", campaignID).Count(&voucherCount).Error; err != nil {
		return nil, errors.New("failed to check voucher count")
	}

	if voucherCount >= int64(campaign.MaxUsers) {
		return nil, fmt.Errorf("voucher limit reached for campaign %d", campaignID)
	}

	// Generate a unique voucher code
	code := generateVoucherCode()

	// Voucher validity period (30 days in this case)
	validFrom := time.Now()
	validTo := validFrom.Add(30 * 24 * time.Hour) // Valid for 30 days

	// Create the voucher instance
	voucherInfo := models.Voucher{
		Code:       code,
		UserID:     userID,
		CampaignID: campaignID,
		Discount:   campaign.Discount, // Use the campaign's discount
		ValidFrom:  validFrom,
		ValidTo:    validTo,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save voucher to DB
	if err := config.DB.Create(&voucherInfo).Error; err != nil {
		return nil, err
	}

	return &voucherInfo, nil
}

// generateVoucherCode generates a random 8-character voucher code.
func generateVoucherCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var code []byte
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 8; i++ {
		code = append(code, charset[rand.Intn(len(charset))])
	}

	return string(code)
}

func ValidateVoucher(userID uint, code string) (*models.VoucherCampaign, error) {
	var voucher models.Voucher
	if err := config.DB.Where("user_id = ? AND code = ? AND valid_to >= ?", userID, code, time.Now()).First(&voucher).Error; err != nil {
		return nil, fmt.Errorf("code %s is invalid or expired voucher", code)
	}

	var campaign models.Campaign

	// Fetch the campaign details
	if err := config.DB.First(&campaign, voucher.CampaignID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch campaign %d", voucher.CampaignID)
	}

	return &models.VoucherCampaign{
		Voucher:  &voucher,
		Campaign: &campaign,
	}, nil
}

func UpdateVoucher(c *gin.Context, voucher models.Voucher) error {
	if err := config.DB.Model(&voucher).Update("used_at", time.Now()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark voucher as used"})
		return err
	}

	return nil
}
