// main.go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/go-playground/validator/v10"
	"github.com/sinscostank/bengkel-inventory/db"
	"github.com/sinscostank/bengkel-inventory/forms"
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

	// 3. Buat Gin router
	router := route.SetupRoutes(dbConn)

	// 4. Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s\n", port)
	router.Run(":" + port)
}
