package main

import (
	"github.com/gin-gonic/gin"
	"gin-gcs-backend/models"
	"gin-gcs-backend/handlers"
)

func main() {
	r := gin.Default()

	// CORS
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

	models.ConnectDatabase()

	r.POST("/products", handlers.UploadProduct)

	r.Run(":8080")
}
