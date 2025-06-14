package models

import (
	"time"

	"gorm.io/gorm"
)


// ActivityItem represents each item in an activity (transaction)
type ActivityItem struct {
	ID             uint               `json:"id" gorm:"primaryKey;autoIncrement"`
	ActivityID     uint               `json:"activity_id" gorm:"not null;index"`
	Activity       Activity           `json:"activity" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID      uint               `json:"product_id" gorm:"not null;index"`
	Product        Product            `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Quantity       int                `json:"quantity" gorm:"not null;check:quantity>0"`
	PriceAtTime    float64            `json:"price_at_time" gorm:"not null;check:price_at_time>=0"`
	DiscountAmount float64            `json:"discount_amount" gorm:"not null;default:0;check:discount_amount>=0"`
	FinalPrice     float64            `json:"final_price" gorm:"not null;check:final_price>=0"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DeletedAt      gorm.DeletedAt     `json:"deleted_at" gorm:"index"`
	StockTx        []StockTransaction `json:"stock_transactions" gorm:"foreignKey:ActivityItemID"`
}