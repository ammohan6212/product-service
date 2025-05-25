package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

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

func getDBConfig() DB {
	return DB{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "yourpassword"),
		DBName:   getEnv("DB_NAME", "mydatabase"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
