package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/pdf"
	"github.com/gin-gonic/gin"
)

// Global session data storage
var sessionData *data.SessionData

// InitializeSessionData initializes the global session data storage
func InitializeSessionData() {
	sessionData = data.NewSessionData()
}

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

	// Load external data button endpoint
	r.POST("/load-external-data", func(c *gin.Context) {
		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		// Load arbiters data from your real API
		err := sessionData.LoadData("arbiters", "https://chess.sk/api/matrika.php/v1/arbiters")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load arbiters: " + err.Error()})
			return
		}

		// Load leagues data from your real API
		err = sessionData.LoadData("leagues", "https://chess.sk/api/leagues.php/v1/leagues")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load leagues: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":         "External data loaded successfully",
			"arbiters_loaded": sessionData.HasData("arbiters"),
			"leagues_loaded":  sessionData.HasData("leagues"),
		})
	})

	// Get loaded data endpoint
	r.GET("/external-data/:type", func(c *gin.Context) {
		dataType := c.Param("type")

		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		data, exists := sessionData.Get(dataType)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "No data loaded for type: " + dataType})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": data})
	})

	// Get arbiters as structured data
	r.GET("/arbiters", func(c *gin.Context) {
		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		arbiters, err := sessionData.GetAllArbiters()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"arbiters": arbiters})
	})

	// Get leagues as structured data
	r.GET("/leagues", func(c *gin.Context) {
		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		leagues, err := sessionData.GetAllLeagues()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"leagues": leagues})
	})

	// Get specific arbiter by ID
	r.GET("/arbiters/:id", func(c *gin.Context) {
		arbiterID := c.Param("id")

		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		// Convert string ID to int
		var id int
		if _, err := fmt.Sscanf(arbiterID, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid arbiter ID"})
			return
		}

		arbiter, err := sessionData.GetArbiterByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"arbiter": arbiter})
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
