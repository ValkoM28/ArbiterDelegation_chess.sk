// Package excel provides functionality for downloading and processing Excel files from chess-results.com.
// It handles the extraction of tournament data, round information, and match details from Excel files.
package excel

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
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/logger"
	"github.com/xuri/excelize/v2"
)

// DownloadChessResultsExcel downloads an Excel file from chess-results.com for the given tournament ID.
// It constructs the appropriate URL and downloads the file to a temporary location.
// The file is saved with a timestamp to avoid conflicts.
// Returns the file path of the downloaded Excel file or an error if the download fails.
func DownloadChessResultsExcel(tournamentID string) (string, error) {
	// Construct the URL for the Excel download
	url := fmt.Sprintf("https://chess-results.com/tnr%s.aspx?lan=1&zeilen=0&art=2&prt=4&excel=2010", tournamentID)

	logger.Debug("Downloading Excel for tournament ID: %s from %s", tournamentID, url)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Second,
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
	bytesWritten, err := io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write Excel file: %v", err)
	}

	logger.Debug("Excel file downloaded: %s (%d bytes)", filePath, bytesWritten)
	return filePath, nil
}

// CleanupTempFile removes the temporary Excel file from the filesystem.
// This function should be called after processing the Excel file to free up disk space.
// Returns an error if the file cannot be removed.
func CleanupTempFile(filePath string) error {
	return os.Remove(filePath)
}

// ExtractTournamentIDFromLeague extracts the tournament ID from a league's ChessResultsLink.
// It parses the URL to find the tournament ID which is used for downloading Excel files.
// The tournament ID is typically found in the URL path after "tnr".
// Returns the tournament ID as a string or an error if the URL format is invalid.
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

	tournamentID := matches[1]
	logger.Debug("Extracted tournament ID %s from league '%s'", tournamentID, league.LeagueName)
	return tournamentID, nil
}

// DownloadExcelForLeague downloads Excel file for a given league
func DownloadExcelForLeague(league *data.League) (string, error) {
	logger.Debug("Downloading Excel for league '%s' (ID: %d)", league.LeagueName, league.LeagueId)

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

	logger.Debug("Excel file downloaded for league '%s': %s", league.LeagueName, filePath)
	return filePath, nil
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// ParseChessResultsExcelToRounds parses an Excel file and returns rounds with matches
func ParseChessResultsExcelToRounds(filePath string) ([]data.Round, error) {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		logger.Error("Failed to open Excel file %s: %v", filePath, err)
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer f.Close()

	// Get the first sheet name
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		logger.Error("No sheets found in Excel file: %s", filePath)
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	// Get all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		logger.Error("Failed to get rows from sheet %s: %v", sheetName, err)
		return nil, fmt.Errorf("failed to get rows from sheet: %v", err)
	}

	var rounds []data.Round
	var currentRound *data.Round

	for _, row := range rows {
		if len(row) == 0 {
			continue
		}

		// Check if this is a round header
		// Format 1: "Round 1" (simple format - date/time from match rows)
		// Format 2: "Round 1 on 2025/10/25 at 11:00" (date/time embedded in header)
		if len(row) == 1 && strings.HasPrefix(row[0], "Round ") {
			// If we have a previous round, add it to the rounds slice
			if currentRound != nil {
				rounds = append(rounds, *currentRound)
			}

			// Try Format 2 first: "Round N on YYYY/MM/DD at HH:MM"
			reWithDate := regexp.MustCompile(`Round (\d+) on (\d{4}/\d{2}/\d{2}) at (\d{2}:\d{2})`)
			matchesWithDate := reWithDate.FindStringSubmatch(row[0])

			if len(matchesWithDate) >= 4 {
				// Format 2: Extract round number and date/time from header
				roundNumber, _ := strconv.Atoi(matchesWithDate[1])
				dateTime := fmt.Sprintf("%s %s", matchesWithDate[2], matchesWithDate[3])

				currentRound = &data.Round{
					Number:   roundNumber,
					DateTime: dateTime,
					Matches:  []data.MatchInfo{},
				}
			} else {
				// Format 1: Extract only round number, date/time will come from match rows
				reSimple := regexp.MustCompile(`Round (\d+)`)
				matchesSimple := reSimple.FindStringSubmatch(row[0])

				if len(matchesSimple) >= 2 {
					roundNumber, _ := strconv.Atoi(matchesSimple[1])

					currentRound = &data.Round{
						Number:   roundNumber,
						DateTime: "", // Will be set from first match
						Matches:  []data.MatchInfo{},
					}
				}
			}
			continue
		}

		// Check if this is the column header row (skip it)
		if len(row) >= 3 && row[0] == "No." && row[1] == "Team" && row[2] == "Team" {
			continue
		}

		// Check if this is a match row
		// Format: [No.] [HomeTeam] [GuestTeam] [Res1] [:] [Res2] [Date] [Time] [Location]
		// Need at least 3 columns (some formats may not have all 9 columns)
		if len(row) >= 3 && isNumeric(row[0]) && row[1] != "" && row[2] != "" {
			// Only add match if we have a current round
			if currentRound != nil {
				var dateTime string
				var address string

				// Check if we have date/time in columns (Format 1)
				if len(row) >= 9 {
					// Extract date and time from columns 6 and 7
					date := strings.TrimSpace(row[6])    // Format: YYYY/MM/DD
					timeStr := strings.TrimSpace(row[7]) // Format: HH:MM
					address = strings.TrimSpace(row[8])  // Location

					// Combine date and time
					dateTime = fmt.Sprintf("%s %s", date, timeStr)

					// Set round's DateTime from first match if not set (Format 1)
					if currentRound.DateTime == "" {
						currentRound.DateTime = dateTime
					}
				} else {
					// Use round's DateTime (Format 2 - date/time from header)
					dateTime = currentRound.DateTime
					address = "" // Address not available in this format
				}

				match := data.MatchInfo{
					HomeTeam:  strings.TrimSpace(row[1]),
					GuestTeam: strings.TrimSpace(row[2]),
					DateTime:  dateTime,
					Address:   address,
				}

				currentRound.Matches = append(currentRound.Matches, match)
			}
		}
	}

	// Add the last round if it exists
	if currentRound != nil {
		rounds = append(rounds, *currentRound)
	}

	logger.Debug("Parsed %d rounds from Excel file %s", len(rounds), filePath)
	return rounds, nil
}

// ParseExcelForLeagueToRounds downloads and parses Excel file for a given league, returning rounds
func ParseExcelForLeagueToRounds(league *data.League) ([]data.Round, error) {
	// Download the Excel file
	filePath, err := DownloadExcelForLeague(league)
	if err != nil {
		return nil, fmt.Errorf("failed to download Excel file: %v", err)
	}

	// Parse the Excel file
	rounds, err := ParseChessResultsExcelToRounds(filePath)
	if err != nil {
		// Clean up Excel file even if parsing fails
		CleanupTempFile(filePath)
		return nil, fmt.Errorf("failed to parse Excel file: %v", err)
	}

	// Clean up Excel file immediately after parsing
	if err := CleanupTempFile(filePath); err != nil {
		logger.Error("Failed to cleanup Excel file %s: %v", filePath, err)
	}

	logger.Info("Parsed %d rounds from Excel for league '%s'", len(rounds), league.LeagueName)
	return rounds, nil
}
