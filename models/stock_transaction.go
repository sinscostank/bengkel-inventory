package models

import (
	"time"

	"gorm.io/gorm"
)

// StockTransaction logs every stock change (inbound/outbound)
type StockTransaction struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID      uint           `json:"product_id" gorm:"not null;index"`
	Product        Product        `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ActivityItemID *uint          `json:"activity_item_id" gorm:"index"`
	ActivityItem   *ActivityItem  `json:"activity_item" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ChangeQuantity int            `json:"change_quantity" gorm:"not null;check:change_quantity<>0"`
	Date           time.Time      `json:"date" gorm:"not null"`
	Note           string         `json:"note" gorm:"size:255"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
