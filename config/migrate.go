// config/migrate.go
package config

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql" // blank import MySQL driver
	_ "github.com/golang-migrate/migrate/v4/source/file"    // blank import file source
)

// RunMigrations applies SQL migrations from db/migrations
func RunMigrations() {
	dsn := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to create migrate instance: %v", err))
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("migration up failed: %v", err))
	}
}
