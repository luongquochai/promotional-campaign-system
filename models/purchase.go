package models

import "time"

// Purchase represents the purchase details of a user
type Purchase struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null" json:"user_id"`                             // User ID is required (foreign key)
	TransactionID   string    `gorm:"size:50;not null;unique;index" json:"transaction_id"` // Transaction ID is unique and required
	SubscriptionID  uint      `gorm:"not null" json:"subscription_id"`                     // Subscription ID (Campaign ID) is required
	DiscountApplied float64   `gorm:"not null" json:"discount_applied"`                    // Discount applied is required
	FinalPrice      float64   `gorm:"not null" json:"final_price"`                         // Final price after discount is required
	Status          string    `gorm:"not null" json:"status"`                              // Status is required (e.g. "completed", "failed")
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`                    // Automatically set the creation time
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`                    // Automatically update the time when modified
}
