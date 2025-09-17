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

// loadExternalData loads arbiters and leagues data from external APIs
func (app *App) loadExternalData(c *gin.Context) {
	fmt.Println("DEBUG: Load external data endpoint called")

	// Parse request body to get season year
	var requestBody struct {
		SeasonStartYear string `json:"seasonStartYear"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Load arbiters data
	err := app.LoadArbiters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load arbiters: " + err.Error()})
		return
	}

	// Load leagues data
	err = app.LoadLeagues(requestBody.SeasonStartYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load leagues: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "External data loaded successfully",
		"arbiters_loaded": app.storage.HasData("arbiters"),
		"leagues_loaded":  app.storage.HasData("leagues"),
	})
}

// getExternalData returns raw external data by type
func (app *App) getExternalData(c *gin.Context) {
	dataType := c.Param("type")

	data, exists := app.storage.Get(dataType)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data loaded for type: " + dataType})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// getArbiters returns all arbiters
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get league by ID
	league, err := app.storage.GetLeagueByID(requestBody.LeagueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "League not found: " + err.Error()})
		return
	}

	// Parse Excel file to get rounds
	rounds, err := excel.ParseExcelForLeagueToRounds(league)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse rounds: " + err.Error()})
		return
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
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

// LoadArbiters loads arbiters data from the API and stores it in storage
func (app *App) LoadArbiters() error {
	// Load arbiters data from your real API with hardcoded active status parameter
	// TODO: Remove client-side filtering when chess.sk API properly supports status=active parameter
	arbitersURL := buildURLWithParams("https://chess.sk/api/matrika.php/v1/arbiters", map[string]string{
		"status": "active", // Currently ignored by API, but kept for when it gets fixed
	})

	err := app.storage.LoadData("arbiters", arbitersURL)
	if err != nil {
		return fmt.Errorf("failed to load arbiters: %v", err)
	}

	// TEMPORARY: Client-side filtering for active arbiters until chess.sk API supports status=active
	arbitersData, exists := app.storage.Get("arbiters")
	if exists {
		filteredArbiters, err := filterActiveArbiters(arbitersData)
		if err != nil {
			return fmt.Errorf("failed to filter arbiters: %v", err)
		}
		app.storage.Set("arbiters", filteredArbiters)
	}

	return nil
}

// LoadLeagues loads leagues data from the API and stores it in storage
func (app *App) LoadLeagues(seasonStartYear string) error {
	// Load leagues data from your real API with season parameter
	leaguesURL := fmt.Sprintf("https://chess.sk/api/leagues.php/v1/leagues?saisonStartYear=%s", seasonStartYear)
	err := app.storage.LoadData("leagues", leaguesURL)
	if err != nil {
		return fmt.Errorf("failed to load leagues: %v", err)
	}

	return nil
}
