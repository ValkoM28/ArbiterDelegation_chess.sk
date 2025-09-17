package app

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/chess"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/excel"
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/pdf"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all HTTP routes with the Gin engine
func (app *App) RegisterRoutes(r *gin.Engine) {
	// Simple PDF generation endpoint
	r.POST("/generate", app.generatePDF)

	// List PDF fields endpoint
	r.GET("/list-fields", app.listFields)

	// Load external data endpoint
	r.POST("/load-external-data", app.loadExternalData)

	// Get loaded data endpoint
	r.GET("/external-data/:type", app.getExternalData)

	// Get arbiters endpoint
	r.GET("/arbiters", app.getArbiters)

	// Get leagues endpoint
	r.GET("/leagues", app.getLeagues)

	// Get specific arbiter by ID
	r.GET("/arbiters/:id", app.getArbiterByID)

	// Get specific league by ID
	r.GET("/leagues/:id", app.getLeagueByID)

	// Prepare PDF data endpoint
	r.POST("/prepare-pdf-data", app.preparePDFData)

	// Download Excel file endpoint
	r.POST("/download-excel", app.downloadExcel)

	// Get rounds data endpoint
	r.POST("/get-rounds", app.getRounds)

	// Delegate arbiters endpoint
	r.POST("/delegate-arbiters", app.delegateArbiters)
}

// generatePDF handles simple PDF generation
func (app *App) generatePDF(c *gin.Context) {
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
}

// listFields lists all fillable fields in the PDF template
func (app *App) listFields(c *gin.Context) {
	// Note: ListFillableFields was removed from the refactored PDF package
	// This endpoint can be removed or implemented differently if needed
	c.JSON(http.StatusOK, gin.H{
		"message": "PDF field listing not available in refactored version",
		"fields":  []string{},
	})
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
	err := chess.LoadArbiters(app.storage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load arbiters: " + err.Error()})
		return
	}

	// Load leagues data
	err = chess.LoadLeagues(app.storage, requestBody.SeasonStartYear)
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
	arbiters, err := chess.GetArbiters(app.storage)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"arbiters": arbiters})
}

// getLeagues returns all leagues
func (app *App) getLeagues(c *gin.Context) {
	leagues, err := chess.GetLeagues(app.storage)
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

	arbiter, err := chess.GetArbiterByID(app.storage, id)
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

	league, err := chess.GetLeagueByID(app.storage, id)
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
	arbiter, err := chess.GetArbiterByID(app.storage, requestBody.ArbiterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Arbiter not found: " + err.Error()})
		return
	}

	// Get league by ID
	league, err := chess.GetLeagueByID(app.storage, requestBody.LeagueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "League not found: " + err.Error()})
		return
	}

	// Prepare PDF data
	pdfData := pdf.PreparePDFDataFromArbiterAndLeague(arbiter, league)

	// Print the data to console
	pdf.PrintPDFData(pdfData)

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
	league, err := chess.GetLeagueByID(app.storage, requestBody.LeagueID)
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
	league, err := chess.GetLeagueByID(app.storage, requestBody.LeagueID)
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

	// Prepare data for function call
	pdf.PrintPDFDataArray(requestBody)

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
