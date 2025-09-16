package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

// DownloadChessResultsExcel downloads an Excel file from chess-results.com
func DownloadChessResultsExcel(tournamentID string) (string, error) {
	// Construct the URL for the Excel download
	url := fmt.Sprintf("https://chess-results.com/tnr%s.aspx?lan=1&zeilen=0&art=2&prt=4&excel=2010", tournamentID)
 	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Second, // Longer timeout for file download
	}

	// Make the HTTP GET request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download Excel file: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("chess-results returned status code: %d", resp.StatusCode)
	}

	// Create a dedicated directory for Excel files
	excelDir := "assets/tempfiles/"
	if err := os.MkdirAll(excelDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create downloads directory: %v", err)
	}

	// Create a permanent file to store the Excel data
	fileName := fmt.Sprintf("chess_results_%s_%d.xlsx", tournamentID, time.Now().Unix())
	filePath := filepath.Join(excelDir, fileName)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write Excel file: %v", err)
	}

	return filePath, nil
}

// CleanupTempFile removes the temporary Excel file
func CleanupTempFile(filePath string) error {
	return os.Remove(filePath)
}

// ExtractTournamentIDFromLeague extracts the tournament ID from a league's ChessResultsLink
func ExtractTournamentIDFromLeague(league *data.League) (string, error) {
	if league.ChessResultsLink == "" {
		return "", fmt.Errorf("league has no ChessResultsLink")
	}

	// Regular expression to match tournament ID in chess-results URLs
	// Pattern: tnr followed by digits, e.g., tnr123456.aspx
	re := regexp.MustCompile(`tnr(\d+)\.aspx`)
	matches := re.FindStringSubmatch(league.ChessResultsLink)

	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract tournament ID from link: %s", league.ChessResultsLink)
	}

	return matches[1], nil
}

// DownloadExcelForLeague downloads Excel file for a given league
func DownloadExcelForLeague(league *data.League) (string, error) {
	// Extract tournament ID from league's ChessResultsLink
	tournamentID, err := ExtractTournamentIDFromLeague(league)
	if err != nil {
		return "", fmt.Errorf("failed to extract tournament ID: %v", err)
	}

	// Download the Excel file
	filePath, err := DownloadChessResultsExcel(tournamentID)
	if err != nil {
		return "", fmt.Errorf("failed to download Excel file: %v", err)
	}

	return filePath, nil
}

