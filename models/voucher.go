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
