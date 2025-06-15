// config/db.go
package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sinscostank/bengkel-inventory/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file")
	}

	// Bentuk DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	// Debug print
	fmt.Println("▶️ Connecting to database with DSN:", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// AutoMigrate all models structs
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Activity{},
		&models.ActivityItem{},
		&models.StockTransaction{},
		&models.PriceHistory{},
	); err != nil {
		panic("AutoMigrate failed: " + err.Error())
	}

	DB = db

	// Return the initialized db connection
	return db, nil	
}
