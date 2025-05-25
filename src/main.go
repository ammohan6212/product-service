package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Serve React build static files
	r.Static("/static", "./frontend/build/static")
	r.LoadHTMLFiles("./frontend/build/index.html")

	// Serve frontend at root
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Simple API endpoint
	r.GET("/api/products", func(c *gin.Context) {
		products := []map[string]string{
			{"id": "1", "name": "Product A"},
			{"id": "2", "name": "Product B"},
		}
		c.JSON(http.StatusOK, gin.H{"products": products})
	})

	// Use port from environment or fallback to 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r.Run(":" + port)
}
