package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category (oli, ban, etc.)
type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"size:255;not null;unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Products  []Product      `json:"products" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}