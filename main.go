package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strconv"

    "go-backend/db"
    "go-backend/models"
)

func loadCategories(filePath string, dbConn *db.DB) {
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
            Name: record[0],
        }
        dbConn.DB.Create(&category)
    }

    fmt.Println("✅ Categories loaded successfully")
}

func loadProducts(filePath string, dbConn *db.DB) {
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
        categoryID, _ := strconv.Atoi(record[2])
        price, _ := strconv.ParseFloat(record[1], 64)

        product := models.Product{
            Name:       record[0],
            Price:      price,
            CategoryID: uint(categoryID),
        }
        dbConn.DB.Create(&product)
    }

    fmt.Println("✅ Products loaded successfully")
}

func main() {
    dbConn := db.Connect()
    dbConn.DB.AutoMigrate(&models.Category{}, &models.Product{})

    loadCategories("data/category.csv", dbConn)
    loadProducts("data/product.csv", dbConn)
}
