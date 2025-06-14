package models

import (
	"time"

	"gorm.io/gorm"
)

// PriceHistory logs price changes for products
type PriceHistory struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID   uint           `json:"product_id" gorm:"not null;index"`
	Product     Product        `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OldPrice    float64        `json:"old_price" gorm:"not null;check:old_price>=0"`
	NewPrice    float64        `json:"new_price" gorm:"not null;check:new_price>=0"`
	DateChanged time.Time      `json:"date_changed" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
