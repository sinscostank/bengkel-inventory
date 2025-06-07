package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user (admin or karyawan)
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Email     string         `json:"email" gorm:"size:255;not null;unique"`
	Password  string         `json:"password" gorm:"size:255;not null"`
	Role      string         `json:"role" gorm:"type:enum('admin','karyawan');not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Category represents a product category (oli, ban, etc.)
type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"size:255;not null;unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Products  []Product      `json:"products" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

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

// Activity represents a sales transaction header
type Activity struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Date      time.Time      `json:"date" gorm:"not null"`
	Status    string         `json:"status" gorm:"type:enum('success','failed');not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Items     []ActivityItem `json:"items" gorm:"foreignKey:ActivityID"`
}

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
