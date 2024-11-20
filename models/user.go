package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:255;uniqueIndex;not null" json:"username"` // Required and unique
	Email     string    `gorm:"size:255;uniqueIndex;not null" json:"email"`    // Required and unique
	Password  string    `gorm:"size:255;not null" json:"password"`             // Required
	Role      string    `gorm:"size:100" json:"role"`                          // Optional
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`              // Auto-create timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`              // Auto-update timestamp
}
