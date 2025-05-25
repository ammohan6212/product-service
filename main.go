package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go-backend/db"
	"go-backend/models"
	"gorm.io/gorm"
	"encoding/json"
	
)

func loadCategories(filePath string, dbConn *gorm.DB) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open category.csv: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("failed to read category.csv: %v", err)
	}

	for _, record := range records[1:] {
		category := models.Category{
			Name:     record[0],
			ImageURL: record[1],
		}
		dbConn.Create(&category)
	}

	fmt.Println("✅ Categories loaded successfully")
}

func loadProducts(filePath string, dbConn *gorm.DB) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open product.csv: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("failed to read product.csv: %v", err)
	}

	for _, record := range records[1:] {
		price, _ := strconv.ParseFloat(record[1], 64)
		categoryID, _ := strconv.Atoi(record[2])
		stock, _ := strconv.Atoi(record[3])

		product := models.Product{
			Name:       record[0],
			Price:      price,
			CategoryID: uint(categoryID),
			Stock:      stock,
			ImageURL:   record[4],
		}
		dbConn.Create(&product)
	}

	fmt.Println("✅ Products loaded successfully")
}
func getCategories(w http.ResponseWriter, r *http.Request) {
    var categories []models.Category
    dbConn.Find(&categories)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
    var products []models.Product
    dbConn.Find(&products)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

func main() {
	dbConn := db.Connect()

	// Auto migrate schema
	dbConn.AutoMigrate(&models.Category{}, &models.Product{})

	loadCategories("data/category.csv", dbConn)
	loadProducts("data/product.csv", dbConn)
	http.HandleFunc("/api/categories", getCategories)
    http.HandleFunc("/api/products", getProducts)

    fmt.Println("🚀 Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
