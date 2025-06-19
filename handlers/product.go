package handlers

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"gin-gcs-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

const (
	bucketName     = "ecommerce-details"
	uploadFolder   = "product-images/"
	publicBaseURL  = "https://storage.googleapis.com/" + bucketName + "/"
)

func UploadProduct(c *gin.Context) {
	sellerName := c.PostForm("sellerName")
	name := c.PostForm("name")
	description := c.PostForm("description")
	price := c.PostForm("price")
	category := c.PostForm("category")

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	defer file.Close()

	imageURL, err := uploadToGCS(c, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"message":   "Product uploaded successfully!",
		"image_url": imageURL,
	})
}

func uploadToGCS(c *gin.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	ext := filepath.Ext(header.Filename)
	uniqueID := uuid.New().String()
	objectName := uploadFolder + uniqueID + ext

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	wc.ContentType = header.Header.Get("Content-Type")
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}} // make public

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write to bucket: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close bucket writer: %v", err)
	}

	// Public URL
	imageURL := publicBaseURL + objectName
	return imageURL, nil
}
