package app

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/excel"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/pdf"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all HTTP routes for the application.
// It sets up both GET and POST endpoints for data retrieval, PDF generation, and Excel processing.
// The routes include endpoints for arbiters, leagues, external data loading, and delegation management.
func (app *App) RegisterRoutes(r *gin.Engine) {
	r.GET("/external-data/:type", app.getExternalData)
	r.GET("/arbiters", app.getArbiters)
	r.GET("/leagues", app.getLeagues)
	r.GET("/arbiters/:id", app.getArbiterByID)
	r.GET("/leagues/:id", app.getLeagueByID)

	r.POST("/prepare-pdf-data", app.preparePDFData)
	r.POST("/download-excel", app.downloadExcel)
	r.POST("/get-rounds", app.getRounds)
	r.POST("/delegate-arbiters", app.delegateArbiters)
	r.POST("/load-external-data", app.loadExternalData)
}

// loadExternalData loads arbiters and leagues data from external APIs.
// It expects a JSON request body with a "seasonStartYear" field.
// The function loads data from chess.sk API for both arbiters and leagues.
// Returns a JSON response indicating success and whether data was loaded.
func (app *App) loadExternalData(c *gin.Context) {
	fmt.Println("========== START loadExternalData handler ==========")
	fmt.Printf("[HANDLER] Request method: %s\n", c.Request.Method)
	fmt.Printf("[HANDLER] Request URL: %s\n", c.Request.URL.String())
	fmt.Printf("[HANDLER] Request headers: %v\n", c.Request.Header)

	// Parse request body to get season year
	var requestBody struct {
		SeasonStartYear string `json:"seasonStartYear"`
	}

	fmt.Println("[HANDLER] Attempting to bind JSON request body")
	if err := c.BindJSON(&requestBody); err != nil {
		fmt.Printf("[HANDLER] ✗ Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Println("========== END loadExternalData handler (bind error) ==========")
		return
	}

	fmt.Printf("[HANDLER] ✓ Request body parsed successfully\n")
	fmt.Printf("[HANDLER] Season start year: %s\n", requestBody.SeasonStartYear)

	// Load arbiters data
	fmt.Println("[HANDLER] Loading arbiters data...")
	arbitersStartTime := time.Now()
	err := app.LoadArbiters()
	arbitersDuration := time.Since(arbitersStartTime)

	if err != nil {
		fmt.Printf("[HANDLER] ✗ Failed to load arbiters (took %v): %v\n", arbitersDuration, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load arbiters: " + err.Error()})
		fmt.Println("========== END loadExternalData handler (arbiters error) ==========")
		return
	}
	fmt.Printf("[HANDLER] ✓ Arbiters loaded successfully in %v\n", arbitersDuration)

	// Load leagues data
	fmt.Printf("[HANDLER] Loading leagues data for season %s...\n", requestBody.SeasonStartYear)
	leaguesStartTime := time.Now()
	err = app.LoadLeagues(requestBody.SeasonStartYear)
	leaguesDuration := time.Since(leaguesStartTime)

	if err != nil {
		fmt.Printf("[HANDLER] ✗ Failed to load leagues (took %v): %v\n", leaguesDuration, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load leagues: " + err.Error()})
		fmt.Println("========== END loadExternalData handler (leagues error) ==========")
		return
	}
	fmt.Printf("[HANDLER] ✓ Leagues loaded successfully in %v\n", leaguesDuration)

	arbitersLoaded := app.storage.HasData("arbiters")
	leaguesLoaded := app.storage.HasData("leagues")

	fmt.Printf("[HANDLER] Storage check - arbiters loaded: %v, leagues loaded: %v\n", arbitersLoaded, leaguesLoaded)

	c.JSON(http.StatusOK, gin.H{
		"message":         "External data loaded successfully",
		"arbiters_loaded": arbitersLoaded,
		"leagues_loaded":  leaguesLoaded,
	})

	fmt.Println("[HANDLER] ✓ Response sent successfully")
	fmt.Println("========== END loadExternalData handler (success) ==========")
}

// getExternalData returns raw external data by type from session storage.
// It expects a URL parameter "type" (arbiters or leagues) and returns the raw data.
// This endpoint is useful for debugging and inspecting the loaded data structure.
func (app *App) getExternalData(c *gin.Context) {
	dataType := c.Param("type")

	data, exists := app.storage.Get(dataType)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data loaded for type: " + dataType})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// getArbiters returns all loaded arbiters from session storage.
// It retrieves and returns all arbiters that have been loaded from the chess.sk API.
// Returns a JSON response with the arbiters array or an error if no data is loaded.
func (app *App) getArbiters(c *gin.Context) {
	arbiters, err := app.storage.GetAllArbiters()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"arbiters": arbiters})
}

// getLeagues returns all leagues
func (app *App) getLeagues(c *gin.Context) {
	leagues, err := app.storage.GetAllLeagues()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leagues": leagues})
}

// getArbiterByID returns a specific arbiter by ID
func (app *App) getArbiterByID(c *gin.Context) {
	arbiterID := c.Param("id")

	// Convert string ID to int
	var id int
	if _, err := fmt.Sscanf(arbiterID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid arbiter ID"})
		return
	}

	arbiter, err := app.storage.GetArbiterByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"arbiter": arbiter})
}

// getLeagueByID returns a specific league by ID
func (app *App) getLeagueByID(c *gin.Context) {
	leagueID := c.Param("id")

	// Convert string ID to int
	var id int
	if _, err := fmt.Sscanf(leagueID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league ID"})
		return
	}

	league, err := app.storage.GetLeagueByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"league": league})
}

// preparePDFData prepares PDF data with arbiter and league from frontend
func (app *App) preparePDFData(c *gin.Context) {
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
	arbiter, err := app.storage.GetArbiterByID(requestBody.ArbiterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arbiter not found: " + err.Error()})
		return
	}

	// Get league by ID
	league, err := app.storage.GetLeagueByID(requestBody.LeagueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "League not found: " + err.Error()})
		return
	}

	// Prepare PDF data
	pdfData := pdf.PreparePDFDataFromArbiterAndLeague(arbiter, league)

	// Return the prepared data to frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "PDF data prepared and printed to console",
		"data":    pdfData,
	})
}

// downloadExcel downloads Excel file for a specific league
func (app *App) downloadExcel(c *gin.Context) {
	// Parse request body to get league ID
	var requestBody struct {
		LeagueID int `json:"leagueId"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get league by ID
	league, err := app.storage.GetLeagueByID(requestBody.LeagueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "League not found: " + err.Error()})
		return
	}

	// Download Excel file for the league
	rounds, err := excel.ParseExcelForLeagueToRounds(league)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download Excel file: " + err.Error()})
		return
	}

	// Return success response with rounds data
	c.JSON(http.StatusOK, gin.H{
		"message": "Excel file downloaded successfully",
		"rounds":  rounds,
		"league":  league.LeagueName,
	})
}

// getRounds gets rounds data for a specific league
func (app *App) getRounds(c *gin.Context) {
	fmt.Println("========== START getRounds handler ==========")
	fmt.Printf("[GET-ROUNDS] Request method: %s\n", c.Request.Method)
	fmt.Printf("[GET-ROUNDS] Request URL: %s\n", c.Request.URL.String())

	// Parse request body to get league ID
	var requestBody struct {
		LeagueID int `json:"leagueId"`
	}

	fmt.Println("[GET-ROUNDS] Attempting to bind JSON request body")
	if err := c.BindJSON(&requestBody); err != nil {
		fmt.Printf("[GET-ROUNDS] ✗ Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		fmt.Println("========== END getRounds handler (bind error) ==========")
		return
	}

	fmt.Printf("[GET-ROUNDS] ✓ Request body parsed successfully\n")
	fmt.Printf("[GET-ROUNDS] League ID: %d\n", requestBody.LeagueID)

	// Get league by ID
	fmt.Printf("[GET-ROUNDS] Fetching league with ID %d from storage\n", requestBody.LeagueID)
	leagueStartTime := time.Now()
	league, err := app.storage.GetLeagueByID(requestBody.LeagueID)
	leagueDuration := time.Since(leagueStartTime)

	if err != nil {
		fmt.Printf("[GET-ROUNDS] ✗ League not found (took %v): %v\n", leagueDuration, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "League not found: " + err.Error()})
		fmt.Println("========== END getRounds handler (league not found) ==========")
		return
	}
	fmt.Printf("[GET-ROUNDS] ✓ League found in %v\n", leagueDuration)
	fmt.Printf("[GET-ROUNDS] League details: Name='%s', ChessResultsLink='%s'\n", league.LeagueName, league.ChessResultsLink)

	// Parse Excel file to get rounds
	fmt.Println("[GET-ROUNDS] Calling excel.ParseExcelForLeagueToRounds()")
	parseStartTime := time.Now()
	rounds, err := excel.ParseExcelForLeagueToRounds(league)
	parseDuration := time.Since(parseStartTime)

	if err != nil {
		fmt.Printf("[GET-ROUNDS] ✗ Failed to parse rounds (took %v): %v\n", parseDuration, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse rounds: " + err.Error()})
		fmt.Println("========== END getRounds handler (parse error) ==========")
		return
	}
	fmt.Printf("[GET-ROUNDS] ✓ Rounds parsed successfully in %v\n", parseDuration)
	fmt.Printf("[GET-ROUNDS] Rounds count: %d\n", len(rounds))

	// Log detailed round information
	for i, round := range rounds {
		fmt.Printf("[GET-ROUNDS] Round %d: Number=%d, DateTime='%s', Matches=%d\n",
			i+1, round.Number, round.DateTime, len(round.Matches))
		for j, match := range round.Matches {
			fmt.Printf("[GET-ROUNDS]   Match %d: %s vs %s at %s\n",
				j+1, match.HomeTeam, match.GuestTeam, match.DateTime)
		}
	}

	// Store rounds in session data for later editing
	fmt.Println("[GET-ROUNDS] Storing rounds in session data")
	app.storage.Set("current_rounds", rounds)
	fmt.Println("[GET-ROUNDS] ✓ Rounds stored in session")

	// Return rounds data
	c.JSON(http.StatusOK, gin.H{
		"message": "Rounds data loaded successfully",
		"rounds":  rounds,
		"league":  league,
	})

	fmt.Println("[GET-ROUNDS] ✓ Response sent successfully")
	fmt.Println("========== END getRounds handler (success) ==========")
}

// delegateArbiters handles the main PDF generation for delegated arbiters
func (app *App) delegateArbiters(c *gin.Context) {
	var requestBody []data.PDFData
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// DEBUG: Log received data
	fmt.Printf("DEBUG: Received %d PDF data items\n", len(requestBody))
	for i, pdfData := range requestBody {
		fmt.Printf("DEBUG: PDF Data %d:\n", i)
		fmt.Printf("  Arbiter: FirstName='%s', LastName='%s', PlayerID='%s'\n",
			pdfData.Arbiter.FirstName, pdfData.Arbiter.LastName, pdfData.Arbiter.PlayerID)
		fmt.Printf("  League: Name='%s', Year='%s'\n",
			pdfData.League.Name, pdfData.League.Year)
		fmt.Printf("  Match: HomeTeam='%s', GuestTeam='%s', DateTime='%s', Address='%s'\n",
			pdfData.Match.HomeTeam, pdfData.Match.GuestTeam, pdfData.Match.DateTime, pdfData.Match.Address)
		fmt.Printf("  Director: Contact='%s'\n", pdfData.Director.Contact)
		fmt.Printf("  ContactPerson: '%s'\n", pdfData.ContactPerson)
	}

	// Generate PDFs
	generatedFiles, err := pdf.GeneratePDFsFromDelegateArbiters(requestBody, "templates/delegacny_list_ligy.pdf")
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
	c.Header("Content-Type", "application/zip")
	c.FileAttachment(zipPath, zipName)
}

// buildURLWithParams constructs a URL with query parameters from a base URL and parameter map.
// It safely parses the base URL and adds the provided parameters as query strings.
// Returns the constructed URL or the original base URL if parsing fails.
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

// LoadArbiters loads arbiters data from the chess.sk API and stores it in session storage.
// It downloads active arbiters data and applies client-side filtering until the API supports status filtering.
// The function handles URL construction, API calls, and data processing.
// Returns an error if the API call fails or data processing encounters issues.
func (app *App) LoadArbiters() error {
	fmt.Println("========== START LoadArbiters ==========")

	// Load arbiters data from your real API with hardcoded active status parameter
	// TODO: Remove client-side filtering when chess.sk API properly supports status=active parameter
	arbitersURL := buildURLWithParams("https://chess.sk/api/matrika.php/v1/arbiters", map[string]string{
		"status": "active", // Currently ignored by API, but kept for when it gets fixed
	})

	fmt.Printf("[LOAD-ARBITERS] Target URL: %s\n", arbitersURL)
	fmt.Println("[LOAD-ARBITERS] Calling storage.LoadData('arbiters', url)")

	loadStartTime := time.Now()
	err := app.storage.LoadData("arbiters", arbitersURL)
	loadDuration := time.Since(loadStartTime)

	if err != nil {
		fmt.Printf("[LOAD-ARBITERS] ✗ Failed to load arbiters data (took %v): %v\n", loadDuration, err)
		return fmt.Errorf("failed to load arbiters: %v", err)
	}
	fmt.Printf("[LOAD-ARBITERS] ✓ Raw arbiters data loaded successfully in %v\n", loadDuration)

	// TEMPORARY: Client-side filtering for active arbiters until chess.sk API supports status=active
	fmt.Println("[LOAD-ARBITERS] Checking if arbiters data exists in storage")
	arbitersData, exists := app.storage.Get("arbiters")
	if exists {
		fmt.Println("[LOAD-ARBITERS] Arbiters data found, applying active filter")
		fmt.Printf("[LOAD-ARBITERS] Raw data type: %T\n", arbitersData)

		filterStartTime := time.Now()
		filteredArbiters, err := filterActiveArbiters(arbitersData)
		filterDuration := time.Since(filterStartTime)

		if err != nil {
			fmt.Printf("[LOAD-ARBITERS] ✗ Failed to filter arbiters (took %v): %v\n", filterDuration, err)
			return fmt.Errorf("failed to filter arbiters: %v", err)
		}
		fmt.Printf("[LOAD-ARBITERS] ✓ Arbiters filtered successfully in %v\n", filterDuration)

		// Log some statistics about filtering
		if dataMap, ok := arbitersData.(map[string]interface{}); ok {
			if dataArray, ok := dataMap["data"].([]interface{}); ok {
				originalCount := len(dataArray)
				if filteredMap, ok := filteredArbiters["data"].([]interface{}); ok {
					filteredCount := len(filteredMap)
					fmt.Printf("[LOAD-ARBITERS] Filtering stats: %d total arbiters -> %d active arbiters\n", originalCount, filteredCount)
				}
			}
		}

		app.storage.Set("arbiters", filteredArbiters)
		fmt.Println("[LOAD-ARBITERS] ✓ Filtered data stored back to session")
	} else {
		fmt.Println("[LOAD-ARBITERS] ⚠ Warning: Arbiters data not found in storage after loading")
	}

	fmt.Println("========== END LoadArbiters (success) ==========")
	return nil
}

// LoadLeagues loads leagues data from the chess.sk API for the specified season and stores it in session storage.
// It constructs the API URL with the season parameter and downloads the leagues data.
// The function handles URL construction, API calls, and data storage.
// Returns an error if the API call fails or data processing encounters issues.
func (app *App) LoadLeagues(seasonStartYear string) error {
	fmt.Println("========== START LoadLeagues ==========")
	fmt.Printf("[LOAD-LEAGUES] Season start year: %s\n", seasonStartYear)

	// Load leagues data from your real API with season parameter
	leaguesURL := fmt.Sprintf("https://chess.sk/api/leagues.php/v1/leagues?saisonStartYear=%s", seasonStartYear)
	fmt.Printf("[LOAD-LEAGUES] Target URL: %s\n", leaguesURL)
	fmt.Println("[LOAD-LEAGUES] Calling storage.LoadData('leagues', url)")

	loadStartTime := time.Now()
	err := app.storage.LoadData("leagues", leaguesURL)
	loadDuration := time.Since(loadStartTime)

	if err != nil {
		fmt.Printf("[LOAD-LEAGUES] ✗ Failed to load leagues data (took %v): %v\n", loadDuration, err)
		return fmt.Errorf("failed to load leagues: %v", err)
	}
	fmt.Printf("[LOAD-LEAGUES] ✓ Leagues data loaded successfully in %v\n", loadDuration)

	// Log some stats about the loaded data
	leaguesData, exists := app.storage.Get("leagues")
	if exists {
		if dataMap, ok := leaguesData.(map[string]interface{}); ok {
			if dataArray, ok := dataMap["data"].([]interface{}); ok {
				fmt.Printf("[LOAD-LEAGUES] Loaded %d leagues\n", len(dataArray))
			}
		}
	} else {
		fmt.Println("[LOAD-LEAGUES] ⚠ Warning: Leagues data not found in storage after loading")
	}

	fmt.Println("========== END LoadLeagues (success) ==========")
	return nil
}
