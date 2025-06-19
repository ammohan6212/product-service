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

const uploadFolder = "product-images/"

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

	imageURL, err := uploadToGCS(file, header)
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

func uploadToGCS(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(header.Filename)
	uniqueID := uuid.New().String()
	objectName := uploadFolder + uniqueID + ext

	wc := gcsclient.Client.Bucket(gcsclient.BucketName).Object(objectName).NewWriter(gcsclient.Ctx)
	wc.ContentType = header.Header.Get("Content-Type")
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write to bucket: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close bucket writer: %v", err)
	}

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", gcsclient.BucketName, objectName)
	return imageURL, nil
}
