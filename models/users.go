package models

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