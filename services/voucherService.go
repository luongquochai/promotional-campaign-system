package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	database "github.com/luongquochai/promotional-campaign-system/database/mysql"
	"github.com/luongquochai/promotional-campaign-system/database/redis" // Import Redis utility
	"github.com/luongquochai/promotional-campaign-system/models"
)

var ctx = context.Background()

// GenerateVoucher generates a new voucher for a user in the given campaign.
func GenerateVoucher(userID uint, campaignID uint) (*models.Voucher, error) {
	// Define cache keys
	voucherCacheKey := fmt.Sprintf("voucher:user:%d:campaign:%d", userID, campaignID)
	campaignCacheKey := fmt.Sprintf("campaign:%d", campaignID)

	// Check Redis cache for campaign
	var campaign models.Campaign
	campaignCached, err := redis.Get(campaignCacheKey)
	if err == nil {
		// Parse cached campaign data
		if err := json.Unmarshal([]byte(campaignCached), &campaign); err != nil {
			return nil, errors.New("failed to parse cached campaign data")
		}
	} else {
		// Fetch campaign details from DB
		if err := database.DB.First(&campaign, campaignID).Error; err != nil {
			return nil, errors.New("campaign not found")
		}

		// Cache campaign data
		campaignData, _ := json.Marshal(campaign)
		redis.Set(campaignCacheKey, string(campaignData), 5*time.Minute)
	}

	// Check campaign status
	if campaign.Status != "active" {
		return nil, errors.New("campaign is not active")
	}
	if time.Now().Before(campaign.StartDate) {
		return nil, errors.New("campaign has not started yet")
	}
	if time.Now().After(campaign.EndDate) {
		return nil, errors.New("campaign has expired")
	}

	// Check Redis cache for voucher
	voucherCached, err := redis.Get(voucherCacheKey)
	if err == nil {
		// If voucher exists in cache, return an error
		return nil, fmt.Errorf("voucher code %s for campaign id %d of user id %d already generated", voucherCached, campaignID, userID)
	}

	// Check if the maximum number of vouchers has been reached
	var voucherCount int64
	if err := database.DB.Model(&models.Voucher{}).Where("campaign_id = ?", campaignID).Count(&voucherCount).Error; err != nil {
		return nil, errors.New("failed to check voucher count")
	}
	if voucherCount >= int64(campaign.MaxUsers) {
		return nil, fmt.Errorf("voucher limit reached for campaign %d", campaignID)
	}

	// Generate a unique voucher code
	code := generateVoucherCode()

	// Voucher validity period (30 days)
	validFrom := time.Now()
	validTo := validFrom.Add(30 * 24 * time.Hour)

	// Create the voucher instance
	voucherInfo := models.Voucher{
		Code:       code,
		UserID:     userID,
		CampaignID: campaignID,
		Discount:   campaign.Discount,
		ValidFrom:  validFrom,
		ValidTo:    validTo,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save voucher to DB
	if err := database.DB.Create(&voucherInfo).Error; err != nil {
		return nil, err
	}

	// Cache the voucher in Redis
	voucherData, _ := json.Marshal(voucherInfo)
	redis.Set(voucherCacheKey, string(voucherData), 30*time.Minute)

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
	// Define cache key
	voucherCacheKey := fmt.Sprintf("voucher:user:%d:code:%s", userID, code)

	// Check Redis cache for voucher
	var voucher models.Voucher
	voucherCached, err := redis.Get(voucherCacheKey)
	if err == nil {
		// Parse cached voucher data
		if err := json.Unmarshal([]byte(voucherCached), &voucher); err != nil {
			return nil, errors.New("failed to parse cached voucher data")
		}
	} else {
		// Fetch voucher from DB
		if err := database.DB.Where("user_id = ? AND code = ? AND valid_to >= ?", userID, code, time.Now()).First(&voucher).Error; err != nil {
			return nil, fmt.Errorf("code %s is invalid or expired voucher", code)
		}

		// Cache the voucher
		voucherData, _ := json.Marshal(voucher)
		redis.Set(voucherCacheKey, string(voucherData), 30*time.Minute)
	}

	var campaign models.Campaign

	// Fetch the campaign details
	if err := database.DB.First(&campaign, voucher.CampaignID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch campaign %d", voucher.CampaignID)
	}

	return &models.VoucherCampaign{
		Voucher:  &voucher,
		Campaign: &campaign,
	}, nil
}

func UpdateVoucher(c *gin.Context, voucher *models.Voucher) error {
	if err := database.DB.Model(&voucher).Update("used_at", time.Now()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark voucher as used"})
		return err
	}

	// Invalidate Redis cache for the voucher
	voucherCacheKey := fmt.Sprintf("voucher:user:%d:campaign:%d", voucher.UserID, voucher.CampaignID)
	redis.Delete(voucherCacheKey)

	return nil
}
