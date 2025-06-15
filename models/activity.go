package models

import (
	"time"

	"gorm.io/gorm"
)

// Activity represents a sales transaction header
type Activity struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Date      time.Time      `json:"date" gorm:"not null"`
	Status    string         `json:"status" gorm:"type:enum('success','failed');not null"`
	Type      string         `json:"type" gorm:"type:enum('inbound','outbound');not null"`	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Items     []ActivityItem `json:"items" gorm:"foreignKey:ActivityID"`
}