package handlers

import (
	"fmt"
	"gin-gcs-backend/gcsclient"
	"gin-gcs-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"strconv"
	"cloud.google.com/go/storage"
)

const (
	uploadFolder  = "product-images/"
	maxUploadSize = 20 << 20 // 20 MB
)

// UploadProduct handles image upload and product creation
func UploadProduct(c *gin.Context) {
	// ✅ Limit the request body size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	// ✅ Validate required fields
	sellerName := strings.TrimSpace(c.PostForm("sellerName"))
	name := strings.TrimSpace(c.PostForm("name"))
	description := strings.TrimSpace(c.PostForm("description"))
	price := strings.TrimSpace(c.PostForm("price"))
	category := strings.TrimSpace(c.PostForm("category"))
	quantityStr := strings.TrimSpace(c.PostForm("quantity"))

	if sellerName == "" || name == "" || price == "" || category == "" || quantityStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields are missing"})
		return
	}

	// ✅ Convert quantity to int
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be a positive integer"})
		return
	}

	// ✅ Handle image file
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Println("Image file read error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	defer file.Close()

	// ✅ Upload to Google Cloud Storage
	imageURL, err := uploadToGCS(file, header)
	if err != nil {
		log.Println("GCS upload error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
		return
	}

	// ✅ Save product to database
	_, err = models.DB.Exec(`
		INSERT INTO products 
		(seller_name, name, description, price, category, quantity, image_path)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		sellerName, name, description, price, category, quantity, imageURL,
	)
	if err != nil {
		log.Println("Database insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product"})
		return
	}

	// ✅ Success response
	c.JSON(http.StatusOK, gin.H{
		"message":   "Product uploaded successfully!",
		"image_url": imageURL,
	})
}


// uploadToGCS uploads the file to the GCS bucket and returns the public URL
func uploadToGCS(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(header.Filename)
	uniqueID := uuid.New().String()
	objectName := uploadFolder + uniqueID + ext

	wc := gcsclient.Client.Bucket(gcsclient.BucketName).Object(objectName).NewWriter(gcsclient.Ctx)
	wc.ContentType = header.Header.Get("Content-Type")
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	// Copy file content to GCS
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write to bucket: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close bucket writer: %v", err)
	}

	// Return public URL
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", gcsclient.BucketName, objectName)
	return imageURL, nil
}


// GetAllProducts returns all products from the database
func GetAllProducts(c *gin.Context) {
	rows, err := models.DB.Query("SELECT id, seller_name, name, description, price, category, quantity, image_path FROM products")
	if err != nil {
		log.Println("Database fetch error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	var products []map[string]interface{}

	for rows.Next() {
		var (
			id          int
			sellerName  string
			name        string
			description string
			price       string
			category    string
			quantity    int
			imagePath   string
		)

		err = rows.Scan(&id, &sellerName, &name, &description, &price, &category, &quantity, &imagePath)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		product := map[string]interface{}{
			"id":           id,
			"seller_name":  sellerName,
			"name":         name,
			"description":  description,
			"price":        price,
			"category":     category,
			"quantity":     quantity,
			"image_url":    imagePath,
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}



// GetProductByID returns a single product by its ID
func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	var (
		productID    int
		sellerName   string
		name         string
		description  string
		price        string
		category     string
		quantity     int
		imagePath    string
	)

	err := models.DB.QueryRow(`
		SELECT id, seller_name, name, description, price, category, quantity, image_path
		FROM products
		WHERE id = ?`, id,
	).Scan(&productID, &sellerName, &name, &description, &price, &category, &quantity, &imagePath)

	if err != nil {
		log.Println("Product fetch error:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product := map[string]interface{}{
		"id":           productID,
		"seller_name":  sellerName,
		"name":         name,
		"description":  description,
		"price":        price,
		"category":     category,
		"quantity":     quantity,
		"image_url":    imagePath,
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}



func GetProductsBySeller(c *gin.Context) {
	seller := c.Query("seller")

	rows, err := models.DB.Query(`
		SELECT id, name, price, quantity, category, image_path
		FROM products
		WHERE seller_name = ?`, seller)
	if err != nil {
		log.Println("DB query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	products := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, price, category, image string
		var quantity int

		err := rows.Scan(&id, &name, &price, &quantity, &category, &image)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}

		products = append(products, map[string]interface{}{
			"id":        id,
			"name":      name,
			"price":     price,
			"quantity":  quantity,
			"category":  category,
			"image_url": image,
		})
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}
// UpdateProductQuantity reduces the quantity of a product after a successful purchase
// UpdateProductQuantity reduces the quantity of a product after a purchase
func UpdateProductQuantity(c *gin.Context) {
	id := c.Param("id")

	// Parse request body for quantityPurchased
	var reqBody struct {
		QuantityPurchased int `json:"quantityPurchased"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil || reqBody.QuantityPurchased <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantityPurchased"})
		return
	}

	// Get current quantity
	var currentQuantity int
	err := models.DB.QueryRow("SELECT quantity FROM products WHERE id = ?", id).Scan(&currentQuantity)
	if err != nil {
		log.Println("Failed to get product quantity:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if currentQuantity < reqBody.QuantityPurchased {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
		return
	}

	// Update quantity
	newQuantity := currentQuantity - reqBody.QuantityPurchased
	_, err = models.DB.Exec("UPDATE products SET quantity = ? WHERE id = ?", newQuantity, id)
	if err != nil {
		log.Println("Failed to update quantity:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quantity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "✅ Product quantity updated",
		"new_quantity": newQuantity,
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var reqBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Category    string `json:"category"`
		Quantity    int    `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := models.DB.Exec(`
		UPDATE products
		SET name = ?, description = ?, price = ?, category = ?, quantity = ?
		WHERE id = ?`,
		reqBody.Name, reqBody.Description, reqBody.Price, reqBody.Category, reqBody.Quantity, id,
	)

	if err != nil {
		log.Println("Update error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Product updated successfully"})
}


func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	_, err := models.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		log.Println("Delete error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Product deleted successfully"})
}



