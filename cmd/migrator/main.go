// Migrator created to execute migrations on database as creating table for next manipulation
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    // reding drovided configuration 
    pUser := os.Getenv("POSTGRES_USER")
    pPassword := os.Getenv("POSTGRES_PASSWORD")
    migrationsPath := os.Getenv("MIGRATIONS_PATH")

    if pUser == "" || pPassword == "" {
        panic("database user, password are required")
    }
    if migrationsPath == "" {
        panic("migrations-path is required")
    }

    // setting up connection to database
    postgresURL := fmt.Sprintf("postgres://%s:%s@postgres:5432/messages?sslmode=disable", pUser, pPassword)
    m, err := migrate.New("file://" + migrationsPath, postgresURL)
    if err != nil {
        panic(err)
    }
    // executing migrations
    if err := m.Up(); err != nil {
        if errors.Is(err, migrate.ErrNoChange) {
            fmt.Println("no migrations to apply")
            return
        }
        panic(err)
    }

    fmt.Println("migrations applied")
}