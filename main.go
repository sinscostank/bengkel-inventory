// main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/sinscostank/bengkel-inventory/config"
	"github.com/sinscostank/bengkel-inventory/route"
)

func main() {
	// 1. Load .env (DB creds, PORT, JWT_SECRET, dsb.)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// 2. Inisialisasi koneksi DB & AutoMigrate via GORM
	config.InitDB()

	// 3. Buat Gin router
	router := gin.Default()

	// 4. (Opsional) Pasang middleware, misalnya CORS atau JWT auth
	// router.Use(middleware.CORSMiddleware())
	// router.Use(middleware.AuthMiddleware())

	// 5. Daftarkan semua route dari route.SetupRoutes
	route.SetupRoutes(router)

	// 6. Jalankan server pada port dari .env (default 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s\n", port)
	router.Run(":" + port)
}
