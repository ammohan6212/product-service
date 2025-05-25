package db

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type DB struct {
    DB *gorm.DB
}

func Connect() *DB {
    dsn := "host=postgres user=postgres password=postgres dbname=go_db port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("❌ Failed to connect to DB: %v", err)
    }

    fmt.Println("✅ Connected to Postgres")
    return &DB{DB: db}
}
