package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Global variable for database connection
var DB *sql.DB

// DBConfig holds PostgreSQL config values
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Connect initializes and connects to PostgreSQL with retries
func Connect() {
	config := getDBConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)

	var err error
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		DB, err = sql.Open("postgres", dsn)
		if err == nil && DB.Ping() == nil {
			log.Println("✅ Connected to PostgreSQL database successfully")
			return
		}
		log.Printf("⏳ Waiting for database connection (%d/%d)...", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("❌ Failed to connect to PostgreSQL after %d retries: %v", maxRetries, err)
}

// getDBConfig fetches DB config from environment or uses fallback
func getDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "yourpassword"),
		DBName:   getEnv("DB_NAME", "mydatabase"),
	}
}

// getEnv returns environment variable or fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
