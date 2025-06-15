// main.go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/go-playground/validator/v10"
	"github.com/sinscostank/bengkel-inventory/controller"
	"github.com/sinscostank/bengkel-inventory/db"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/repository"
	"github.com/sinscostank/bengkel-inventory/route"
)

func main() {

	// Create a new validator instance
	validate := validator.New()

	// Register the custom validator function
	validate.RegisterValidation("fullname", forms.ValidateFullName)

	// 1. Load .env (DB creds, PORT, JWT_SECRET, dsb.)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// 2. Inisialisasi koneksi DB & AutoMigrate via GORM
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
		return
	}

	// Create repository
	userRepo := repository.NewUserRepository(dbConn)
	productRepo := repository.NewProductRepository(dbConn)
	categoryRepo := repository.NewCategoryRepository(dbConn)
	activityRepo := repository.NewActivityRepository(dbConn)
	actiityItemRepo := repository.NewActivityItemRepository(dbConn)
	stockTransactionRepo := repository.NewStockTransactionRepository(dbConn)
	priceHistoryRepo := repository.NewPriceHistoryRepository(dbConn)

	// Create controllers
	userController := controller.NewUserController(userRepo)
	productController := controller.NewProductController(productRepo, categoryRepo, priceHistoryRepo)
	categoryController := controller.NewCategoryController(categoryRepo)
	activityController := controller.NewActivityController(activityRepo, actiityItemRepo, stockTransactionRepo, productRepo)
	

	// 3. Buat Gin router
	router := route.SetupRoutes(userController, productController, categoryController, activityController)

	// 4. (Opsional) Pasang middleware, misalnya CORS atau JWT auth
	// router.Use(middleware.CORSMiddleware())
	// router.Use(middleware.AuthMiddleware())

	// 6. Jalankan server pada port dari .env (default 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s\n", port)
	router.Run(":" + port)
}
