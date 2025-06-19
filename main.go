package main

import (
	"gin-gcs-backend/gcsclient"
	"gin-gcs-backend/handlers"
	"gin-gcs-backend/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	// Initialize MySQL
	models.ConnectDatabase()

	// Initialize Google Cloud Storage
	err := gcsclient.ConnectGCS()
	if err != nil {
		log.Fatal("Failed to connect to Google Cloud Storage:", err)
	}
	defer gcsclient.CloseGCS()

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Product service is running",
		})
	})

	// Route
	r.POST("/products", handlers.UploadProduct)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
