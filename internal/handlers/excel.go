package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
	"github.com/xuri/excelize/v2"
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

func ParseChessResultsExcel(filepath string) ([]data.MatchInfo, error) {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	// Get all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet: %v", err)
	}

	var matches []data.MatchInfo
	var currentRoundDateTime string
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}

		// Check if this is a round header (e.g., "Round 1 on 2025/10/25 at 11:00")
		if strings.HasPrefix(row[0], "Round ") && strings.Contains(row[0], " on ") {
			// Extract date and time from round header
			// Format: "Round X on YYYY/MM/DD at HH:MM"
			re := regexp.MustCompile(`Round \d+ on (\d{4}/\d{2}/\d{2}) at (\d{2}:\d{2})`)
			matches := re.FindStringSubmatch(row[0])
			if len(matches) >= 3 {
				currentRoundDateTime = fmt.Sprintf("%s %s", matches[1], matches[2])
			}
			continue
		}

		// Check if this is a match row (starts with a number and has two team names)
		if len(row) >= 3 && isNumeric(row[0]) && row[1] != "" && row[2] != "" {
			// Skip the header row "No.,Team,Team,Res.,:,Res."
			if row[1] == "Team" && row[2] == "Team" {
				continue
			}

			// Create MatchInfo from the row data
			match := data.MatchInfo{
				HomeTeam:  strings.TrimSpace(row[1]),
				GuestTeam: strings.TrimSpace(row[2]),
				DateTime:  currentRoundDateTime,
				Address:   "", // Address not available in this Excel format
			}

			matches = append(matches, match)
		}
	}

	return matches, nil
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// ParseExcelForLeague downloads and parses Excel file for a given league
func ParseExcelForLeague(league *data.League) ([]data.MatchInfo, error) {
	// Download the Excel file
	filePath, err := DownloadExcelForLeague(league)
	if err != nil {
		return nil, fmt.Errorf("failed to download Excel file: %v", err)
	}

	// Parse the Excel file
	matches, err := ParseChessResultsExcel(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Excel file: %v", err)
	}

	PrintMatches(matches)
	return matches, nil
}

// PrintMatches prints the parsed matches data for debugging
func PrintMatches(matches []data.MatchInfo) {
	fmt.Printf("Found %d matches:\n", len(matches))
	fmt.Println(80 * '=')

	for i, match := range matches {
		fmt.Printf("Match %d:\n", i+1)
		fmt.Printf("  Home Team: %s\n", match.HomeTeam)
		fmt.Printf("  Guest Team: %s\n", match.GuestTeam)
		fmt.Printf("  Date/Time: %s\n", match.DateTime)
		fmt.Printf("  Address: %s\n", match.Address)
		fmt.Println('-' * 40)
	}
}
