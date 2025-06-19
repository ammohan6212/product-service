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

	if sellerName == "" || name == "" || price == "" || category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required fields are missing"})
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
		(seller_name, name, description, price, category, image_path)
		VALUES (?, ?, ?, ?, ?, ?)`,
		sellerName, name, description, price, category, imageURL,
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
