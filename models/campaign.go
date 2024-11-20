package models

import "time"

type Campaign struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Discount  float64   `gorm:"not null" json:"discount_percentage"`
	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`
	MaxUsers  int       `gorm:"not null" json:"max_users"`
	Status    string    `gorm:"type:enum('active', 'paused', 'expired');default:'active';not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Foreign key relationship with User model (if applicable)
	CreatorID uint `gorm:"not null" json:"creator_id"`          // ID of the user who created the campaign
	Creator   User `gorm:"foreignKey:CreatorID" json:"creator"` // The associated User who created the campaign
}
