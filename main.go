package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go-backend/db"
	"go-backend/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

func loadCategories(filePath string, dbConn *gorm.DB) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("❌ Failed to open category.csv: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("❌ Failed to read category.csv: %v", err)
	}

	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		category := models.Category{
			Name:     record[1],
			ImageURL: record[2],
		}
		dbConn.Create(&category)
	}

	fmt.Println("✅ Categories loaded successfully")
}

func loadProducts(filePath string, dbConn *gorm.DB) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("❌ Failed to open product.csv: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("❌ Failed to read product.csv: %v", err)
	}

	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		price, _ := strconv.ParseFloat(record[2], 64)
		stock, _ := strconv.Atoi(record[3])
		categoryID, _ := strconv.Atoi(record[5])

		product := models.Product{
			Name:       record[1],
			Price:      price,
			Stock:      stock,
			ImageURL:   record[4],
			CategoryID: uint(categoryID),
		}
		dbConn.Create(&product)
	}

	fmt.Println("✅ Products loaded successfully")
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	DB.Find(&categories)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	DB.Preload("Category").Find(&products)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func main() {
	DB = db.Connect()

	DB.AutoMigrate(&models.Category{}, &models.Product{})

	// DB.Exec("DELETE FROM products")
	// DB.Exec("DELETE FROM categories")
	// DB.Exec("ALTER SEQUENCE categories_id_seq RESTART WITH 1") // optional: reset auto-increment
	// DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")   // optional: reset auto-increment
	

	loadCategories("data/category.csv", DB)
	loadProducts("data/product.csv", DB)

	http.HandleFunc("/categories", getCategories)
	http.HandleFunc("/products", getProducts)

	fmt.Println("🚀 Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
