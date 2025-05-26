package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Connect() *gorm.DB {
	config := DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "yourpassword",
		DBName:   "mydatabase",
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)

	var err error
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Connected to PostgreSQL via GORM successfully")
			return DB
		}
		log.Printf("⏳ Waiting for DB connection via GORM (%d/%d)...", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("❌ Failed to connect to PostgreSQL with GORM: %v", err)
	return nil
}
