package app

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/excel"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/logger"
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
	// Parse request body to get season year
	var requestBody struct {
		SeasonStartYear string `json:"seasonStartYear"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		logger.Error("Failed to parse loadExternalData request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	logger.Info("Loading external data for season %s", requestBody.SeasonStartYear)

	// Load arbiters data
	if err := app.LoadArbiters(); err != nil {
		logger.Error("Failed to load arbiters: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load arbiters: " + err.Error()})
		return
	}

	// Load leagues data
	if err := app.LoadLeagues(requestBody.SeasonStartYear); err != nil {
		logger.Error("Failed to load leagues for season %s: %v", requestBody.SeasonStartYear, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load leagues: " + err.Error()})
		return
	}

	logger.Info("External data loaded successfully - Arbiters: %v, Leagues: %v",
		app.storage.HasData("arbiters"), app.storage.HasData("leagues"))

	c.JSON(http.StatusOK, gin.H{
		"message":         "External data loaded successfully",
		"arbiters_loaded": app.storage.HasData("arbiters"),
		"leagues_loaded":  app.storage.HasData("leagues"),
	})
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
	// Parse request body to get league ID
	var requestBody struct {
		LeagueID int `json:"leagueId"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		logger.Error("Failed to parse getRounds request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	logger.Info("Fetching rounds for league ID: %d", requestBody.LeagueID)

	// Get league by ID
	league, err := app.storage.GetLeagueByID(requestBody.LeagueID)
	if err != nil {
		logger.Error("League not found (ID: %d): %v", requestBody.LeagueID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "League not found: " + err.Error()})
		return
	}

	// Parse Excel file to get rounds
	rounds, err := excel.ParseExcelForLeagueToRounds(league)
	if err != nil {
		logger.Error("Failed to parse rounds for league '%s': %v", league.LeagueName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse rounds: " + err.Error()})
		return
	}

	logger.Info("Successfully loaded %d rounds for league '%s'", len(rounds), league.LeagueName)

	// Store rounds in session data for later editing
	app.storage.Set("current_rounds", rounds)

	// Return rounds data
	c.JSON(http.StatusOK, gin.H{
		"message": "Rounds data loaded successfully",
		"rounds":  rounds,
		"league":  league,
	})
}

// delegateArbiters handles the main PDF generation for delegated arbiters
func (app *App) delegateArbiters(c *gin.Context) {
	var requestBody []data.PDFData
	if err := c.BindJSON(&requestBody); err != nil {
		logger.Error("Failed to parse delegateArbiters request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	logger.Info("Generating PDFs for %d arbiters", len(requestBody))
	logger.Debug("PDF generation data: %+v", requestBody)

	// Generate PDFs
	generatedFiles, err := pdf.GeneratePDFsFromDelegateArbiters(requestBody, "templates/delegacny_list_ligy.pdf")
	if err != nil {
		logger.Error("Failed to generate PDFs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDFs: " + err.Error()})
		return
	}

	// Create zip file with all generated PDFs
	zipName := fmt.Sprintf("delegacne_listy_%d.zip", time.Now().Unix())
	zipPath, err := pdf.CreateZipFromFiles(generatedFiles, zipName)
	if err != nil {
		logger.Error("Failed to create zip file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create zip file: " + err.Error()})
		return
	}

	// Clean up individual PDF files after creating zip
	for _, file := range generatedFiles {
		if err := os.Remove(file); err != nil {
			logger.Error("Failed to remove temporary PDF file %s: %v", file, err)
		}
	}

	logger.Info("Successfully generated delegation package: %s", zipName)

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
	// Load arbiters data from your real API with hardcoded active status parameter
	// TODO: Remove client-side filtering when chess.sk API properly supports status=active parameter
	arbitersURL := buildURLWithParams("https://chess.sk/api/matrika.php/v1/arbiters", map[string]string{
		"status": "active", // Currently ignored by API, but kept for when it gets fixed
	})

	logger.Debug("Loading arbiters from: %s", arbitersURL)

	if err := app.storage.LoadData("arbiters", arbitersURL); err != nil {
		return fmt.Errorf("failed to load arbiters: %v", err)
	}

	// TEMPORARY: Client-side filtering for active arbiters until chess.sk API supports status=active
	arbitersData, exists := app.storage.Get("arbiters")
	if !exists {
		return fmt.Errorf("arbiters data not found in storage after loading")
	}

	filteredArbiters, err := filterActiveArbiters(arbitersData)
	if err != nil {
		return fmt.Errorf("failed to filter arbiters: %v", err)
	}

	// Log filtering statistics
	if dataMap, ok := arbitersData.(map[string]interface{}); ok {
		if dataArray, ok := dataMap["data"].([]interface{}); ok {
			originalCount := len(dataArray)
			if filteredMap, ok := filteredArbiters["data"].([]interface{}); ok {
				filteredCount := len(filteredMap)
				logger.Info("Loaded arbiters: %d total -> %d active", originalCount, filteredCount)
			}
		}
	}

	app.storage.Set("arbiters", filteredArbiters)
	return nil
}

// LoadLeagues loads leagues data from the chess.sk API for the specified season and stores it in session storage.
// It constructs the API URL with the season parameter and downloads the leagues data.
// The function handles URL construction, API calls, and data storage.
// Returns an error if the API call fails or data processing encounters issues.
func (app *App) LoadLeagues(seasonStartYear string) error {
	// Load leagues data from your real API with season parameter
	leaguesURL := fmt.Sprintf("https://chess.sk/api/leagues.php/v1/leagues?saisonStartYear=%s", seasonStartYear)

	logger.Debug("Loading leagues from: %s", leaguesURL)

	if err := app.storage.LoadData("leagues", leaguesURL); err != nil {
		return fmt.Errorf("failed to load leagues: %v", err)
	}

	// Log statistics about loaded leagues
	leaguesData, exists := app.storage.Get("leagues")
	if exists {
		if dataMap, ok := leaguesData.(map[string]interface{}); ok {
			if dataArray, ok := dataMap["data"].([]interface{}); ok {
				logger.Info("Loaded %d leagues for season %s", len(dataArray), seasonStartYear)
			}
		}
	}

	return nil
}
