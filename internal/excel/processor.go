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
	"github.com/xuri/excelize/v2"
)

// DownloadChessResultsExcel downloads an Excel file from chess-results.com for the given tournament ID.
// It constructs the appropriate URL and downloads the file to a temporary location.
// The file is saved with a timestamp to avoid conflicts.
// Returns the file path of the downloaded Excel file or an error if the download fails.
func DownloadChessResultsExcel(tournamentID string) (string, error) {
	fmt.Println("========== START DownloadChessResultsExcel ==========")
	fmt.Printf("[DOWNLOAD-EXCEL] Tournament ID: %s\n", tournamentID)

	// Construct the URL for the Excel download
	url := fmt.Sprintf("https://chess-results.com/tnr%s.aspx?lan=1&zeilen=0&art=2&prt=4&excel=2010", tournamentID)
	fmt.Printf("[DOWNLOAD-EXCEL] Target URL: %s\n", url)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Second, // Longer timeout for file download
	}
	fmt.Println("[DOWNLOAD-EXCEL] HTTP client created with 60s timeout")

	// Make the HTTP GET request
	fmt.Println("[DOWNLOAD-EXCEL] Making HTTP GET request to chess-results.com...")
	requestStartTime := time.Now()
	resp, err := client.Get(url)
	requestDuration := time.Since(requestStartTime)

	if err != nil {
		fmt.Printf("[DOWNLOAD-EXCEL] ✗ HTTP request failed (took %v): %v\n", requestDuration, err)
		return "", fmt.Errorf("failed to download Excel file: %v", err)
	}
	defer resp.Body.Close()
	fmt.Printf("[DOWNLOAD-EXCEL] ✓ HTTP request completed in %v\n", requestDuration)
	fmt.Printf("[DOWNLOAD-EXCEL] Response status: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("[DOWNLOAD-EXCEL] Response headers: %v\n", resp.Header)

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[DOWNLOAD-EXCEL] ✗ chess-results.com returned non-OK status: %d\n", resp.StatusCode)
		return "", fmt.Errorf("chess-results returned status code: %d", resp.StatusCode)
	}
	fmt.Println("[DOWNLOAD-EXCEL] ✓ Response status is OK (200)")

	// Create a dedicated directory for Excel files
	excelDir := "assets/tempfiles/"
	fmt.Printf("[DOWNLOAD-EXCEL] Creating directory: %s\n", excelDir)
	if err := os.MkdirAll(excelDir, 0755); err != nil {
		fmt.Printf("[DOWNLOAD-EXCEL] ✗ Failed to create directory: %v\n", err)
		return "", fmt.Errorf("failed to create downloads directory: %v", err)
	}
	fmt.Println("[DOWNLOAD-EXCEL] ✓ Directory created/verified")

	// Create a permanent file to store the Excel data
	fileName := fmt.Sprintf("chess_results_%s_%d.xlsx", tournamentID, time.Now().Unix())
	filePath := filepath.Join(excelDir, fileName)
	fmt.Printf("[DOWNLOAD-EXCEL] Target file path: %s\n", filePath)

	// Create the file
	fmt.Println("[DOWNLOAD-EXCEL] Creating file...")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("[DOWNLOAD-EXCEL] ✗ Failed to create file: %v\n", err)
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	fmt.Println("[DOWNLOAD-EXCEL] ✓ File created successfully")

	// Copy the response body to the file
	fmt.Println("[DOWNLOAD-EXCEL] Writing Excel data to file...")
	copyStartTime := time.Now()
	bytesWritten, err := io.Copy(file, resp.Body)
	copyDuration := time.Since(copyStartTime)

	if err != nil {
		fmt.Printf("[DOWNLOAD-EXCEL] ✗ Failed to write Excel file (took %v): %v\n", copyDuration, err)
		return "", fmt.Errorf("failed to write Excel file: %v", err)
	}
	fmt.Printf("[DOWNLOAD-EXCEL] ✓ Excel file written successfully in %v\n", copyDuration)
	fmt.Printf("[DOWNLOAD-EXCEL] Bytes written: %d\n", bytesWritten)

	fmt.Printf("[DOWNLOAD-EXCEL] ✓ Excel file saved to: %s\n", filePath)
	fmt.Println("========== END DownloadChessResultsExcel (success) ==========")
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
	fmt.Println("========== START ExtractTournamentIDFromLeague ==========")
	fmt.Printf("[EXTRACT-ID] League: %s\n", league.LeagueName)
	fmt.Printf("[EXTRACT-ID] ChessResultsLink: %s\n", league.ChessResultsLink)

	if league.ChessResultsLink == "" {
		fmt.Println("[EXTRACT-ID] ✗ League has no ChessResultsLink")
		return "", fmt.Errorf("league has no ChessResultsLink")
	}

	// Regular expression to match tournament ID in chess-results URLs
	// Pattern: tnr followed by digits, e.g., tnr123456.aspx
	fmt.Println("[EXTRACT-ID] Applying regex pattern: tnr(\\d+)\\.aspx")
	re := regexp.MustCompile(`tnr(\d+)\.aspx`)
	matches := re.FindStringSubmatch(league.ChessResultsLink)

	fmt.Printf("[EXTRACT-ID] Regex matches: %v\n", matches)
	fmt.Printf("[EXTRACT-ID] Number of matches: %d\n", len(matches))

	if len(matches) < 2 {
		fmt.Printf("[EXTRACT-ID] ✗ Could not extract tournament ID from link: %s\n", league.ChessResultsLink)
		return "", fmt.Errorf("could not extract tournament ID from link: %s", league.ChessResultsLink)
	}

	tournamentID := matches[1]
	fmt.Printf("[EXTRACT-ID] ✓ Extracted tournament ID: %s\n", tournamentID)
	fmt.Println("========== END ExtractTournamentIDFromLeague (success) ==========")
	return tournamentID, nil
}

// DownloadExcelForLeague downloads Excel file for a given league
func DownloadExcelForLeague(league *data.League) (string, error) {
	fmt.Println("========== START DownloadExcelForLeague ==========")
	fmt.Printf("[DOWNLOAD-LEAGUE-EXCEL] League: %s (ID: %d)\n", league.LeagueName, league.LeagueId)

	// Extract tournament ID from league's ChessResultsLink
	fmt.Println("[DOWNLOAD-LEAGUE-EXCEL] Extracting tournament ID from ChessResultsLink")
	extractStartTime := time.Now()
	tournamentID, err := ExtractTournamentIDFromLeague(league)
	extractDuration := time.Since(extractStartTime)

	if err != nil {
		fmt.Printf("[DOWNLOAD-LEAGUE-EXCEL] ✗ Failed to extract tournament ID (took %v): %v\n", extractDuration, err)
		return "", fmt.Errorf("failed to extract tournament ID: %v", err)
	}
	fmt.Printf("[DOWNLOAD-LEAGUE-EXCEL] ✓ Tournament ID extracted in %v: %s\n", extractDuration, tournamentID)

	// Download the Excel file
	fmt.Printf("[DOWNLOAD-LEAGUE-EXCEL] Downloading Excel file for tournament ID: %s\n", tournamentID)
	downloadStartTime := time.Now()
	filePath, err := DownloadChessResultsExcel(tournamentID)
	downloadDuration := time.Since(downloadStartTime)

	if err != nil {
		fmt.Printf("[DOWNLOAD-LEAGUE-EXCEL] ✗ Failed to download Excel file (took %v): %v\n", downloadDuration, err)
		return "", fmt.Errorf("failed to download Excel file: %v", err)
	}
	fmt.Printf("[DOWNLOAD-LEAGUE-EXCEL] ✓ Excel file downloaded successfully in %v\n", downloadDuration)
	fmt.Printf("[DOWNLOAD-LEAGUE-EXCEL] File path: %s\n", filePath)

	fmt.Println("========== END DownloadExcelForLeague (success) ==========")
	return filePath, nil
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// ParseChessResultsExcelToRounds parses an Excel file and returns rounds with matches
func ParseChessResultsExcelToRounds(filePath string) ([]data.Round, error) {
	//fmt.Println("========== START ParseChessResultsExcelToRounds ==========")
	//fmt.Printf("[PARSE-EXCEL] File path: %s\n", filePath)

	// Open the Excel file
	//fmt.Println("[PARSE-EXCEL] Opening Excel file...")
	openStartTime := time.Now()
	f, err := excelize.OpenFile(filePath)
	openDuration := time.Since(openStartTime)

	if err != nil {
		fmt.Printf("[PARSE-EXCEL] ✗ Failed to open Excel file (took %v): %v\n", openDuration, err)
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer f.Close()
	//fmt.Printf("[PARSE-EXCEL] ✓ Excel file opened successfully in %v\n", openDuration)

	// Get the first sheet name
	sheetName := f.GetSheetName(0)
	//fmt.Printf("[PARSE-EXCEL] First sheet name: '%s'\n", sheetName)

	if sheetName == "" {
		fmt.Println("[PARSE-EXCEL] ✗ No sheets found in Excel file")
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	// Get all rows from the sheet
	//fmt.Println("[PARSE-EXCEL] Reading rows from sheet...")
	rowsStartTime := time.Now()
	rows, err := f.GetRows(sheetName)
	rowsDuration := time.Since(rowsStartTime)

	if err != nil {
		fmt.Printf("[PARSE-EXCEL] ✗ Failed to get rows from sheet (took %v): %v\n", rowsDuration, err)
		return nil, fmt.Errorf("failed to get rows from sheet: %v", err)
	}
	//fmt.Printf("[PARSE-EXCEL] ✓ Rows read successfully in %v\n", rowsDuration)
	//fmt.Printf("[PARSE-EXCEL] Total rows in sheet: %d\n", len(rows))

	var rounds []data.Round
	var currentRound *data.Round
	roundHeaderCount := 0
	matchRowCount := 0
	emptyRowCount := 0

	//fmt.Println("[PARSE-EXCEL] Starting row-by-row processing...")
	for rowIndex, row := range rows {
		if len(row) == 0 {
			emptyRowCount++
			continue
		}

		// Check if this is a round header
		// Format 1: "Round 1" (simple format - date/time from match rows)
		// Format 2: "Round 1 on 2025/10/25 at 11:00" (date/time embedded in header)
		if len(row) == 1 && strings.HasPrefix(row[0], "Round ") {
			roundHeaderCount++
			//fmt.Printf("[PARSE-EXCEL] Row %d: Found round header: '%s'\n", rowIndex+1, row[0])

			// If we have a previous round, add it to the rounds slice
			if currentRound != nil {
				//fmt.Printf("[PARSE-EXCEL]   Finalizing previous round %d with %d matches\n", currentRound.Number, len(currentRound.Matches))
				rounds = append(rounds, *currentRound)
			}

			// Try Format 2 first: "Round N on YYYY/MM/DD at HH:MM"
			reWithDate := regexp.MustCompile(`Round (\d+) on (\d{4}/\d{2}/\d{2}) at (\d{2}:\d{2})`)
			matchesWithDate := reWithDate.FindStringSubmatch(row[0])

			if len(matchesWithDate) >= 4 {
				// Format 2: Extract round number and date/time from header
				roundNumber, _ := strconv.Atoi(matchesWithDate[1])
				dateTime := fmt.Sprintf("%s %s", matchesWithDate[2], matchesWithDate[3])

				//fmt.Printf("[PARSE-EXCEL]   Extracted (Format 2): Round=%d, DateTime='%s'\n", roundNumber, dateTime)

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

					//fmt.Printf("[PARSE-EXCEL]   Extracted (Format 1): Round=%d (DateTime will be set from matches)\n", roundNumber)

					currentRound = &data.Round{
						Number:   roundNumber,
						DateTime: "", // Will be set from first match
						Matches:  []data.MatchInfo{},
					}
				} //else {
				//fmt.Printf("[PARSE-EXCEL]   ⚠ Warning: Could not parse round header format\n")
				//}
			}
			continue
		}

		// Check if this is the column header row (skip it)
		if len(row) >= 3 && row[0] == "No." && row[1] == "Team" && row[2] == "Team" {
			fmt.Printf("[PARSE-EXCEL] Row %d: Skipping column header row\n", rowIndex+1)
			continue
		}

		// Check if this is a match row
		// Format: [No.] [HomeTeam] [GuestTeam] [Res1] [:] [Res2] [Date] [Time] [Location]
		// Need at least 3 columns (some formats may not have all 9 columns)
		if len(row) >= 3 && isNumeric(row[0]) && row[1] != "" && row[2] != "" {
			// Only add match if we have a current round
			if currentRound != nil {
				matchRowCount++

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
						//fmt.Printf("[PARSE-EXCEL]   Set round %d DateTime to: %s\n", currentRound.Number, dateTime)
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

				//fmt.Printf("[PARSE-EXCEL] Row %d: Adding match to round %d: %s vs %s at %s (%s)\n",
				//rowIndex+1, currentRound.Number, match.HomeTeam, match.GuestTeam, match.DateTime, match.Address)

				currentRound.Matches = append(currentRound.Matches, match)
			} //else {
			//fmt.Printf("[PARSE-EXCEL] Row %d: ⚠ Warning: Found match row but no current round\n", rowIndex+1)
			//}
		}
	}

	// Add the last round if it exists
	if currentRound != nil {
		//fmt.Printf("[PARSE-EXCEL] Finalizing last round %d with %d matches\n", currentRound.Number, len(currentRound.Matches))
		rounds = append(rounds, *currentRound)
	}

	//fmt.Println("[PARSE-EXCEL] Processing complete - Summary:")
	//fmt.Printf("[PARSE-EXCEL]   Total rows processed: %d\n", len(rows))
	//fmt.Printf("[PARSE-EXCEL]   Empty rows skipped: %d\n", emptyRowCount)
	//fmt.Printf("[PARSE-EXCEL]   Round headers found: %d\n", roundHeaderCount)
	//fmt.Printf("[PARSE-EXCEL]   Match rows found: %d\n", matchRowCount)
	//fmt.Printf("[PARSE-EXCEL]   Total rounds created: %d\n", len(rounds))

	// Log detailed information about each round
	//for i, round := range rounds {
	//fmt.Printf("[PARSE-EXCEL]   Round %d: Number=%d, DateTime='%s', Matches=%d\n",
	//i+1, round.Number, round.DateTime, len(round.Matches))
	//}

	//fmt.Println("========== END ParseChessResultsExcelToRounds (success) ==========")
	return rounds, nil
}

// ParseExcelForLeagueToRounds downloads and parses Excel file for a given league, returning rounds
func ParseExcelForLeagueToRounds(league *data.League) ([]data.Round, error) {
	//fmt.Println("========== START ParseExcelForLeagueToRounds ==========")
	//fmt.Printf("[PARSE-LEAGUE-EXCEL] League: %s (ID: %d)\n", league.LeagueName, league.LeagueId)
	//fmt.Printf("[PARSE-LEAGUE-EXCEL] ChessResultsLink: %s\n", league.ChessResultsLink)

	// Download the Excel file
	//fmt.Println("[PARSE-LEAGUE-EXCEL] Step 1: Downloading Excel file")
	downloadStartTime := time.Now()
	filePath, err := DownloadExcelForLeague(league)
	downloadDuration := time.Since(downloadStartTime)

	if err != nil {
		fmt.Printf("[PARSE-LEAGUE-EXCEL] ✗ Failed to download Excel file (took %v): %v\n", downloadDuration, err)
		return nil, fmt.Errorf("failed to download Excel file: %v", err)
	}
	//fmt.Printf("[PARSE-LEAGUE-EXCEL] ✓ Excel file downloaded successfully in %v\n", downloadDuration)
	//fmt.Printf("[PARSE-LEAGUE-EXCEL] Downloaded file path: %s\n", filePath)

	// Parse the Excel file
	//fmt.Println("[PARSE-LEAGUE-EXCEL] Step 2: Parsing Excel file")
	//parseStartTime := time.Now()
	rounds, err := ParseChessResultsExcelToRounds(filePath)
	//parseDuration := time.Since(parseStartTime)

	if err != nil {
		// Clean up Excel file even if parsing fails
		//fmt.Printf("[PARSE-LEAGUE-EXCEL] ✗ Failed to parse Excel file (took %v): %v\n", parseDuration, err)
		//fmt.Println("[PARSE-LEAGUE-EXCEL] Cleaning up Excel file after parse error")
		CleanupTempFile(filePath)
		return nil, fmt.Errorf("failed to parse Excel file: %v", err)
	}
	//fmt.Printf("[PARSE-LEAGUE-EXCEL] ✓ Excel file parsed successfully in %v\n", parseDuration)
	//fmt.Printf("[PARSE-LEAGUE-EXCEL] Parsed %d rounds\n", len(rounds))

	// Clean up Excel file immediately after parsing
	//fmt.Println("[PARSE-LEAGUE-EXCEL] Step 3: Cleaning up temporary Excel file")
	cleanupStartTime := time.Now()
	if err := CleanupTempFile(filePath); err != nil {
		cleanupDuration := time.Since(cleanupStartTime)
		fmt.Printf("[PARSE-LEAGUE-EXCEL] ⚠ Warning: failed to cleanup Excel file %s (took %v): %v\n", filePath, cleanupDuration, err)
	} else {
		cleanupDuration := time.Since(cleanupStartTime)
		fmt.Printf("[PARSE-LEAGUE-EXCEL] ✓ Excel file cleaned up successfully in %v\n", cleanupDuration)
	}

	totalDuration := time.Since(downloadStartTime)
	fmt.Printf("[PARSE-LEAGUE-EXCEL] Total operation time: %v\n", totalDuration)
	fmt.Println("========== END ParseExcelForLeagueToRounds (success) ==========")
	return rounds, nil
}
