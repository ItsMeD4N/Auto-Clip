package main

import (
	"auto-clip/api"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	_ = godotenv.Load()

	r := gin.Default()

	// CORS Setup (Simple)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	r.POST("/generate", api.HandleGenerate)
	r.GET("/status/:id", api.HandleStatus)
	r.GET("/download/:id", api.HandleDownload)

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
