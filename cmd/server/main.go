// Package main provides the entry point for the chess arbiter delegation generator server.
// It sets up the HTTP server, serves static assets, and registers API routes.
package main

import (
	"log"
	"os"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/app"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/logger"
	"github.com/gin-gonic/gin"
)

// main is the entry point of the application.
// It initializes the application, sets up the Gin router, serves static files,
// registers API routes, and starts the HTTP server on port 8080.
func main() {
	// Initialize logger
	enableDebug := os.Getenv("DEBUG") == "true"
	if err := logger.Init("logs", enableDebug); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Clean up old logs (keep last 30 days)
	if err := logger.CleanOldLogs("logs", 30); err != nil {
		logger.Error("Failed to clean old logs: %v", err)
	}

	logger.Info("Starting Chess Arbiter Delegation Generator")

	// Create new App instance with all dependencies
	application := app.New()

	// Set Gin mode based on environment
	if !enableDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Serve static assets
	r.Static("/assets", "./web/assets")

	// Register API routes
	application.RegisterRoutes(r)

	// Serve frontend for root path
	r.GET("/", func(c *gin.Context) {
		c.File("web/index.html")
	})

	logger.Info("Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		logger.Error("Server failed to start: %v", err)
		log.Fatal(err)
	}
}
