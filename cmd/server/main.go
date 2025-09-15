package main

import (
	"log"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create Gin router
	r := gin.Default()

	// Serve static assets
	r.Static("/assets", "./web/assets")

	// Register API routes
	handlers.RegisterRoutes(r)

	// Serve frontend for root path
	r.GET("/", func(c *gin.Context) {
		c.File("web/index.html")
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(r.Run(":8080"))
}
