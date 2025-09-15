package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/pdf"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/generate", func(c *gin.Context) {
		var payload map[string]string
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pdfPath, err := pdf.FillForm("templates/delegacny_list_ligy.pdf", payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.FileAttachment(pdfPath, "delegacny.pdf")
	})

	r.GET("/list-fields", func(c *gin.Context) {
		err := pdf.ListFillableFields("templates/delegacny_list_ligy.pdf")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Fields listed to console. Check server logs for details."})
	})

	// Test endpoint to try our HTTP client
	r.GET("/test-api", func(c *gin.Context) {
		// Let's test with a free public API
		data, err := GetDataFromApi("https://jsonplaceholder.typicode.com/posts/1")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": data})
	})
}

// GetDataFromApi makes a simple HTTP GET request to an external API
func GetDataFromApi(url string) (map[string]interface{}, error) {
	// Create an HTTP client with a timeout
	// Why timeout? Without it, if the API is slow or down, your app could hang forever
	client := &http.Client{
		Timeout: 30 * time.Second, // 30 seconds max wait time
	}

	// Make the HTTP GET request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close() // Always close the response body to free resources

	// Check if the response was successful (status code 200-299)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse the JSON response into a Go map
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return result, nil
}
