package models

import "time"

// Voucher represents a discount voucher for a user within a campaign
type Voucher struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Code       string     `gorm:"size:255;uniqueIndex" json:"code"` // Use VARCHAR with a maximum size of 255
	UserID     uint       `gorm:"not null" json:"user_id"`          // User who the voucher belongs to
	CampaignID uint       `gorm:"not null" json:"campaign_id"`      // Campaign for which the voucher was generated
	Discount   float64    `gorm:"not null" json:"discount"`         // Discount percentage (e.g., 0.30 for 30%)
	ValidFrom  time.Time  `gorm:"not null" json:"valid_from"`       // When the voucher becomes valid
	ValidTo    time.Time  `gorm:"not null" json:"valid_to"`         // When the voucher expires
	UsedAt     *time.Time `json:"used_at,omitempty"`                // When the voucher was used (null if unused)
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type VoucherCampaign struct {
	Voucher  *Voucher
	Campaign *Campaign
}

// VoucherResponse represents the response format for a voucher along with its campaign
// swagger:response VoucherResponse
type VoucherResponse struct {
	// in: body
	Body Voucher
}

// VoucherCampaign represents the response format for a voucher along with its campaign
// swagger:response VoucherCampaignResponse
type VoucherCampaignResponse struct {
	// in: body
	Body VoucherCampaign
}

// VoucherValidationResponse represents the response structure when validating a voucher
// swagger:response VoucherValidationResponse
type VoucherValidationResponse struct {
	IsUsed       bool    `json:"is_used"`       // Indicates if the voucher has been used
	CampaignID   uint    `json:"campaign_id"`   // The campaign ID associated with the voucher
	CampaignName string  `json:"campaign_name"` // The name of the campaign
	DiscountRate float64 `json:"discount_rate"` // The discount rate for the campaign
}
