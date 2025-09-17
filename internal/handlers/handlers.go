package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/pdf"
	"github.com/gin-gonic/gin"
)

// Global session data storage
var sessionData *data.SessionData

// buildURLWithParams constructs a URL with query parameters
func buildURLWithParams(baseURL string, params map[string]string) string {
	if len(params) == 0 {
		return baseURL
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// filterActiveArbiters filters arbiters to only include active ones
// TEMPORARY: This function should be removed when chess.sk API properly supports status=active parameter
func filterActiveArbiters(rawData interface{}) (map[string]interface{}, error) {
	// Extract the actual data array from our wrapped structure
	dataMap, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("raw data is not a map")
	}

	dataArray, ok := dataMap["data"]
	if !ok {
		return nil, fmt.Errorf("no 'data' field in raw data")
	}

	// Convert to slice of interfaces
	arbitersSlice, ok := dataArray.([]interface{})
	if !ok {
		return nil, fmt.Errorf("data is not an array")
	}

	// Filter for active arbiters
	var activeArbiters []interface{}
	for _, arbiterInterface := range arbitersSlice {
		arbiterMap, ok := arbiterInterface.(map[string]interface{})
		if !ok {
			continue // Skip invalid entries
		}

		// Check if IsActive is true
		if isActive, exists := arbiterMap["IsActive"]; exists {
			if isActiveBool, ok := isActive.(bool); ok && isActiveBool {
				activeArbiters = append(activeArbiters, arbiterInterface)
			}
		}
	}

	// Create new data structure with filtered arbiters
	resultMap := make(map[string]interface{})
	resultMap["data"] = activeArbiters

	return resultMap, nil
}

// InitializeSessionData initializes the global session data storage
func InitializeSessionData() {
	fmt.Println("DEBUG: Initializing session data...")
	sessionData = data.NewSessionData()
	fmt.Println("DEBUG: Session data initialized successfully")
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
		fieldNames, err := pdf.ListFillableFields("templates/delegacny_list_ligy.pdf")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Fields listed successfully",
			"fields":  fieldNames,
		})
	})

	// Load external data button endpoint
	r.POST("/load-external-data", func(c *gin.Context) {
		fmt.Println("DEBUG: Load external data endpoint called")
		if sessionData == nil {
			fmt.Println("DEBUG: Session data is nil!")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}
		fmt.Println("DEBUG: Session data is initialized, proceeding with API calls")

		// Parse request body to get season year
		var requestBody struct {
			SeasonStartYear string `json:"seasonStartYear"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Load arbiters data from your real API with hardcoded active status parameter
		// TODO: Remove client-side filtering when chess.sk API properly supports status=active parameter
		arbitersURL := buildURLWithParams("https://chess.sk/api/matrika.php/v1/arbiters", map[string]string{
			"status": "active", // Currently ignored by API, but kept for when it gets fixed
		})

		err := sessionData.LoadData("arbiters", arbitersURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load arbiters: " + err.Error()})
			return
		}

		// TEMPORARY: Client-side filtering for active arbiters until chess.sk API supports status=active
		arbitersData, exists := sessionData.Get("arbiters")
		if exists {
			filteredArbiters, err := filterActiveArbiters(arbitersData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter arbiters: " + err.Error()})
				return
			}
			sessionData.Set("arbiters", filteredArbiters)
		}

		// Load leagues data from your real API with season parameter
		leaguesURL := fmt.Sprintf("https://chess.sk/api/leagues.php/v1/leagues?saisonStartYear=%s", requestBody.SeasonStartYear)
		err = sessionData.LoadData("leagues", leaguesURL)
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

	// Get specific league by ID
	r.GET("/leagues/:id", func(c *gin.Context) {
		leagueID := c.Param("id")

		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		// Convert string ID to int
		var id int
		if _, err := fmt.Sscanf(leagueID, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
			return
		}

		league, err := sessionData.GetLeagueByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"league": league})
	})

	// Prepare PDF data with arbiter and league from frontend
	r.POST("/prepare-pdf-data", func(c *gin.Context) {
		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		// Parse request body to get arbiter and league IDs
		var requestBody struct {
			ArbiterID int `json:"arbiterId"`
			LeagueID  int `json:"leagueId"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Get arbiter by ID
		arbiter, err := sessionData.GetArbiterByID(requestBody.ArbiterID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Arbiter not found: " + err.Error()})
			return
		}

		// Get league by ID
		league, err := sessionData.GetLeagueByID(requestBody.LeagueID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "League not found: " + err.Error()})
			return
		}

		// Prepare PDF data
		pdfData := PreparePDFData(arbiter, league)

		// Print the data to console
		PrintPDFData(pdfData)

		// Return the prepared data to frontend
		c.JSON(http.StatusOK, gin.H{
			"message": "PDF data prepared and printed to console",
			"data":    pdfData,
		})
	})

	// Download Excel file for a specific league
	r.POST("/download-excel", func(c *gin.Context) {
		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		// Parse request body to get league ID
		var requestBody struct {
			LeagueID int `json:"leagueId"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Get league by ID
		league, err := sessionData.GetLeagueByID(requestBody.LeagueID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "League not found: " + err.Error()})
			return
		}

		// Download Excel file for the league
		filePath, err := ParseExcelForLeagueToRounds(league)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download Excel file: " + err.Error()})
			return
		}

		// Return success response with file path
		c.JSON(http.StatusOK, gin.H{
			"message":  "Excel file downloaded successfully",
			"filePath": filePath,
			"league":   league.LeagueName,
		})
	})

	// Get rounds data for a specific league
	r.POST("/get-rounds", func(c *gin.Context) {
		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		// Parse request body to get league ID
		var requestBody struct {
			LeagueID int `json:"leagueId"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Get league by ID
		league, err := sessionData.GetLeagueByID(requestBody.LeagueID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "League not found: " + err.Error()})
			return
		}

		// Parse Excel file to get rounds
		rounds, err := ParseExcelForLeagueToRounds(league)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse rounds: " + err.Error()})
			return
		}

		// Store rounds in session data for later editing
		sessionData.Set("current_rounds", rounds)

		// Return rounds data
		c.JSON(http.StatusOK, gin.H{
			"message": "Rounds data loaded successfully",
			"rounds":  rounds,
			"league":  league,
		})
	})

	r.POST("/delegate-arbiters", func(c *gin.Context) {
		if sessionData == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not initialized"})
			return
		}

		var requestBody []data.PDFData
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Prepare data for function call
		printPDFDataArray(requestBody)

		// Convert PDFData array to interface{} array for PDF generation
		var pdfDataArray []interface{}
		for _, pdfData := range requestBody {
			// Convert PDFData to map[string]interface{}
			pdfMap := make(map[string]interface{})

			// Convert arbiter
			arbiterMap := make(map[string]interface{})
			arbiterMap["firstName"] = pdfData.Arbiter.FirstName
			arbiterMap["lastName"] = pdfData.Arbiter.LastName
			arbiterMap["playerId"] = pdfData.Arbiter.PlayerID
			pdfMap["arbiter"] = arbiterMap

			// Convert league
			leagueMap := make(map[string]interface{})
			leagueMap["name"] = pdfData.League.Name
			leagueMap["year"] = pdfData.League.Year
			pdfMap["league"] = leagueMap

			// Convert match
			matchMap := make(map[string]interface{})
			matchMap["homeTeam"] = pdfData.Match.HomeTeam
			matchMap["guestTeam"] = pdfData.Match.GuestTeam
			matchMap["dateTime"] = pdfData.Match.DateTime
			matchMap["address"] = pdfData.Match.Address
			pdfMap["match"] = matchMap

			// Convert director
			directorMap := make(map[string]interface{})
			directorMap["contact"] = pdfData.Director.Contact
			pdfMap["director"] = directorMap

			// Add contact person
			pdfMap["contactPerson"] = pdfData.ContactPerson

			pdfDataArray = append(pdfDataArray, pdfMap)
		}

		// Generate PDFs
		generatedFiles, err := pdf.GeneratePDFsFromDelegateArbiters(pdfDataArray, "templates/delegacny_list_ligy.pdf")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDFs: " + err.Error()})
			return
		}

		// Create zip file with all generated PDFs
		zipName := fmt.Sprintf("delegacne_listy_%d.zip", time.Now().Unix())
		zipPath, err := pdf.CreateZipFromFiles(generatedFiles, zipName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create zip file: " + err.Error()})
			return
		}

		// Clean up individual PDF files after creating zip
		for _, file := range generatedFiles {
			if err := os.Remove(file); err != nil {
				fmt.Printf("Warning: failed to remove temporary PDF file %s: %v\n", file, err)
			}
		}

		// Return the zip file for download
		c.FileAttachment(zipPath, zipName)
	})
}

// printPDFDataArray is a temporary function to print PDFData array for debugging
func printPDFDataArray(pdfDataArray []data.PDFData) {
	fmt.Printf("\n=== PDFData Array Debug Output ===\n")
	fmt.Printf("Total items: %d\n", len(pdfDataArray))
	fmt.Printf("=====================================\n")

	for i, pdfData := range pdfDataArray {
		fmt.Printf("\n--- Item %d ---\n", i+1)
		fmt.Printf("League: %s (%s)\n", pdfData.League.Name, pdfData.League.Year)
		fmt.Printf("Director: %s\n", pdfData.Director.Contact)
		fmt.Printf("Arbiter: %s %s (ID: %s)\n", pdfData.Arbiter.FirstName, pdfData.Arbiter.LastName, pdfData.Arbiter.PlayerID)
		fmt.Printf("Match: %s vs %s\n", pdfData.Match.HomeTeam, pdfData.Match.GuestTeam)
		fmt.Printf("DateTime: %s\n", pdfData.Match.DateTime)
		fmt.Printf("Address: %s\n", pdfData.Match.Address)
		fmt.Printf("Contact Person: %s\n", pdfData.ContactPerson)
	}

	fmt.Printf("\n=== End Debug Output ===\n")
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
