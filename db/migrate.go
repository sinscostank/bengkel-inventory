package db

import (
    "log"

	"github.com/sinscostank/bengkel-inventory/models"
    "gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
    err := db.AutoMigrate(
        &models.User{},
        &models.Category{},
        &models.Product{},
        &models.Activity{},
        &models.ActivityItem{},
        &models.StockTransaction{},
        &models.PriceHistory{},
    )

    if err != nil {
        log.Fatalf("AutoMigrate failed: %v", err)
    }

    log.Println("✅ Database migrated successfully.")
}