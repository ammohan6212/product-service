package handlers

import (
	"fmt"
	"gin-gcs-backend/gcsclient"
	"gin-gcs-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"cloud.google.com/go/storage"
)

const (
	uploadFolder   = "product-images/"
	maxUploadSize  = 20 << 20 // 20 MB
)

func UploadProduct(c *gin.Context) {
	// ✅ Enforce max upload size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	// Extract form fields
	sellerName := c.PostForm("sellerName")
	name := c.PostForm("name")
	description := c.PostForm("description")
	price := c.PostForm("price")
	category := c.PostForm("category")

	// Handle image file
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	defer file.Close()

	// Upload to GCS
	imageURL, err := uploadToGCS(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
		return
	}

	// Save product in DB
	_, err = models.DB.Exec(`
        INSERT INTO products 
        (seller_name, name, description, price, category, image_path)
        VALUES (?, ?, ?, ?, ?, ?)`,
		sellerName, name, description, price, category, imageURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product"})
		return
	}

	// Respond to client
	c.JSON(http.StatusOK, gin.H{
		"message":   "Product uploaded successfully!",
		"image_url": imageURL,
	})
}

func uploadToGCS(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(header.Filename)
	uniqueID := uuid.New().String()
	objectName := uploadFolder + uniqueID + ext

	wc := gcsclient.Client.Bucket(gcsclient.BucketName).Object(objectName).NewWriter(gcsclient.Ctx)
	wc.ContentType = header.Header.Get("Content-Type")
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}} // make public

	// Upload the file
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write to bucket: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close bucket writer: %v", err)
	}

	// Generate public URL
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", gcsclient.BucketName, objectName)
	return imageURL, nil
}
