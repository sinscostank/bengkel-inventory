package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents an item in inventory
type Product struct {
	ID         uint               `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string             `json:"name" gorm:"size:255;not null"`
	Stock      int                `json:"stock" gorm:"not null;check:stock>=0"`
	Price      float64            `json:"price" gorm:"not null;check:price>=0"`
	Location   string             `json:"location" gorm:"size:255;not null"`
	CategoryID uint               `json:"category_id" gorm:"not null;index"`
	Category   Category           `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	DeletedAt  gorm.DeletedAt     `json:"deleted_at" gorm:"index"`
	Items      []ActivityItem     `json:"items" gorm:"foreignKey:ProductID"`
	StockTx    []StockTransaction `json:"stock_transactions" gorm:"foreignKey:ProductID"`
	PriceHist  []PriceHistory     `json:"price_history" gorm:"foreignKey:ProductID"`
}
