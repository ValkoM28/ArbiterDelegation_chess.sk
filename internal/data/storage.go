// Package data provides data models and structures for the chess arbiter delegation generator.
// It includes models for arbiters, leagues, matches, and PDF generation data.
package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// SessionData represents a simple in-memory data storage for the application session.
// It provides thread-safe storage for arbiters, leagues, and other data loaded from external APIs.
// This is a simplified implementation without timestamps or expiration - data is loaded once and used.
type SessionData struct {
	data  map[string]interface{} // The actual data storage map
	mutex sync.RWMutex           // Read-write mutex for thread safety
}

// NewSessionData creates a new SessionData instance with an empty data map.
// Returns a pointer to a new SessionData ready for use.
func NewSessionData() *SessionData {
	return &SessionData{
		data: make(map[string]interface{}),
	}
}

// Get retrieves data from session storage by key.
// Returns the data and a boolean indicating whether the key exists.
// This method is thread-safe and uses read locking.
func (sd *SessionData) Get(key string) (interface{}, bool) {
	sd.mutex.RLock()         // Lock for reading
	defer sd.mutex.RUnlock() // Unlock when function exits

	data, exists := sd.data[key]
	return data, exists
}

// Set stores data in session storage with the given key.
// This method is thread-safe and uses write locking.
func (sd *SessionData) Set(key string, value interface{}) {
	sd.mutex.Lock()         // Lock for writing
	defer sd.mutex.Unlock() // Unlock when function exits

	sd.data[key] = value
}

// LoadData loads data from an external API URL and stores it with the given key.
// This method makes an HTTP GET request to the provided URL and stores the response.
// Returns an error if the HTTP request fails or if the response cannot be processed.
func (sd *SessionData) LoadData(key string, url string) error {
	fmt.Println("========== START SessionData.LoadData ==========")
	fmt.Printf("[LOAD-DATA] Key: %s\n", key)
	fmt.Printf("[LOAD-DATA] URL: %s\n", url)

	// This will call our HTTP client function
	fmt.Println("[LOAD-DATA] Calling fetchFromAPI()")
	fetchStartTime := time.Now()
	data, err := fetchFromAPI(url)
	fetchDuration := time.Since(fetchStartTime)

	if err != nil {
		fmt.Printf("[LOAD-DATA] ✗ fetchFromAPI failed (took %v): %v\n", fetchDuration, err)
		return err
	}
	fmt.Printf("[LOAD-DATA] ✓ fetchFromAPI succeeded in %v\n", fetchDuration)

	// Store the data
	fmt.Printf("[LOAD-DATA] Storing data with key '%s'\n", key)
	sd.Set(key, data)
	fmt.Printf("[LOAD-DATA] ✓ Data stored successfully\n")

	fmt.Println("========== END SessionData.LoadData (success) ==========")
	return nil
}

// Clear removes all data from session storage.
// This method is thread-safe and resets the data map to empty.
func (sd *SessionData) Clear() {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	sd.data = make(map[string]interface{})
}

// HasData checks if data exists for the given key in session storage.
// Returns true if the key exists, false otherwise.
// This method is thread-safe and uses read locking.
func (sd *SessionData) HasData(key string) bool {
	sd.mutex.RLock()
	defer sd.mutex.RUnlock()

	_, exists := sd.data[key]
	return exists
}

// fetchFromAPI makes an HTTP GET request to the specified URL and returns the response data.
// It creates an HTTP client with a 30-second timeout and handles the response parsing.
// The response is expected to be JSON and will be wrapped in a map with a "data" key.
// Returns an error if the request fails, the response status is not OK, or JSON parsing fails.
func fetchFromAPI(url string) (map[string]interface{}, error) {
	fmt.Println("========== START fetchFromAPI ==========")
	fmt.Printf("[FETCH-API] Target URL: %s\n", url)

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	fmt.Println("[FETCH-API] HTTP client created with 30s timeout")

	// Make the HTTP GET request
	fmt.Println("[FETCH-API] Making HTTP GET request...")
	requestStartTime := time.Now()
	resp, err := client.Get(url)
	requestDuration := time.Since(requestStartTime)

	if err != nil {
		fmt.Printf("[FETCH-API] ✗ HTTP request failed (took %v): %v\n", requestDuration, err)
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()
	fmt.Printf("[FETCH-API] ✓ HTTP request completed in %v\n", requestDuration)
	fmt.Printf("[FETCH-API] Response status: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("[FETCH-API] Response headers: %v\n", resp.Header)

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[FETCH-API] ✗ API returned non-OK status: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}
	fmt.Println("[FETCH-API] ✓ Response status is OK (200)")

	// Read the response body
	fmt.Println("[FETCH-API] Reading response body...")
	readStartTime := time.Now()
	body, err := io.ReadAll(resp.Body)
	readDuration := time.Since(readStartTime)

	if err != nil {
		fmt.Printf("[FETCH-API] ✗ Failed to read response body (took %v): %v\n", readDuration, err)
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	fmt.Printf("[FETCH-API] ✓ Response body read successfully in %v\n", readDuration)
	fmt.Printf("[FETCH-API] Response body size: %d bytes\n", len(body))

	// Show a preview of the response body (first 500 characters)
	bodyPreview := string(body)
	if len(bodyPreview) > 500 {
		bodyPreview = bodyPreview[:500] + "... (truncated)"
	}
	fmt.Printf("[FETCH-API] Response body preview: %s\n", bodyPreview)

	// Parse the JSON response - it could be an array or an object
	fmt.Println("[FETCH-API] Parsing JSON response...")
	parseStartTime := time.Now()
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("[FETCH-API] ✗ Failed to parse JSON (took %v): %v\n", time.Since(parseStartTime), err)
		fmt.Printf("[FETCH-API] Body that failed to parse: %s\n", string(body))
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}
	parseDuration := time.Since(parseStartTime)
	fmt.Printf("[FETCH-API] ✓ JSON parsed successfully in %v\n", parseDuration)
	fmt.Printf("[FETCH-API] Parsed data type: %T\n", result)

	// Convert to map format for consistency
	resultMap := make(map[string]interface{})
	resultMap["data"] = result

	// Log statistics about the data
	switch v := result.(type) {
	case []interface{}:
		fmt.Printf("[FETCH-API] Result is an array with %d items\n", len(v))
	case map[string]interface{}:
		fmt.Printf("[FETCH-API] Result is an object with %d keys\n", len(v))
	default:
		fmt.Printf("[FETCH-API] Result is of type: %T\n", result)
	}

	fmt.Println("========== END fetchFromAPI (success) ==========")
	return resultMap, nil
}

// ProcessData is a generic function that processes raw API data and extracts structured data.
// It takes raw data (typically from an API response) and converts it to a slice of the specified type T.
// The function expects the raw data to be wrapped in a map with a "data" key containing the actual array.
// It uses JSON marshaling/unmarshaling to convert the data to the target type.
// Returns an error if the data structure is invalid or conversion fails.
func ProcessData[T any](rawData interface{}) ([]T, error) {
	// Extract the actual data array from our wrapped structure
	dataMap, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("raw data is not a map")
	}

	dataArray, ok := dataMap["data"]
	if !ok {
		return nil, fmt.Errorf("no 'data' field in raw data")
	}

	// Convert to JSON bytes
	jsonData, err := json.Marshal(dataArray)
	if err != nil {
		return nil, fmt.Errorf("error marshaling raw data: %v", err)
	}

	// Parse into structured data
	var result []T
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling data: %v", err)
	}

	return result, nil
}

// ProcessArbitersData processes raw API data and extracts a slice of Arbiter structs.
// This is a convenience function that calls ProcessData with the Arbiter type.
// It's used specifically for processing arbiters data from the chess.sk API.
func ProcessArbitersData(rawData interface{}) ([]Arbiter, error) {
	return ProcessData[Arbiter](rawData)
}

// ProcessLeaguesData processes raw API data and extracts a slice of League structs.
// This is a convenience function that calls ProcessData with the League type.
// It's used specifically for processing leagues data from the chess.sk API.
func ProcessLeaguesData(rawData interface{}) ([]League, error) {
	return ProcessData[League](rawData)
}

// GetArbiterByID finds an arbiter by ID from the loaded arbiters data.
// It searches through all loaded arbiters and returns the one with the matching PlayerId.
// The arbiterID parameter is converted to string for comparison with the PlayerId field.
// Returns a pointer to the found Arbiter or an error if not found or data not loaded.
func (sd *SessionData) GetArbiterByID(arbiterID int) (*Arbiter, error) {
	rawData, exists := sd.Get("arbiters")
	if !exists {
		return nil, fmt.Errorf("arbiters data not loaded")
	}

	arbiters, err := ProcessArbitersData(rawData)
	if err != nil {
		return nil, err
	}

	for _, arbiter := range arbiters {
		if arbiter.PlayerId == fmt.Sprintf("%d", arbiterID) {
			return &arbiter, nil
		}
	}

	return nil, fmt.Errorf("arbiter with ID %d not found", arbiterID)
}

// GetLeagueByID finds a league by ID from the loaded leagues data.
// It searches through all loaded leagues and returns the one with the matching LeagueId.
// The leagueID parameter is converted to string for comparison with the LeagueId field.
// Returns a pointer to the found League or an error if not found or data not loaded.
func (sd *SessionData) GetLeagueByID(leagueID int) (*League, error) {
	rawData, exists := sd.Get("leagues")
	if !exists {
		return nil, fmt.Errorf("leagues data not loaded")
	}

	leagues, err := ProcessLeaguesData(rawData)
	if err != nil {
		return nil, err
	}

	for _, league := range leagues {
		if league.LeagueId == fmt.Sprintf("%d", leagueID) {
			return &league, nil
		}
	}

	return nil, fmt.Errorf("league with ID %d not found", leagueID)
}

// GetAllArbiters returns all loaded arbiters from the session storage.
// It retrieves the raw arbiters data and processes it into a slice of Arbiter structs.
// Returns an error if arbiters data has not been loaded yet.
func (sd *SessionData) GetAllArbiters() ([]Arbiter, error) {
	rawData, exists := sd.Get("arbiters")
	if !exists {
		return nil, fmt.Errorf("arbiters data not loaded")
	}

	return ProcessArbitersData(rawData)
}

// GetAllLeagues returns all loaded leagues from the session storage.
// It retrieves the raw leagues data and processes it into a slice of League structs.
// Returns an error if leagues data has not been loaded yet.
func (sd *SessionData) GetAllLeagues() ([]League, error) {
	rawData, exists := sd.Get("leagues")
	if !exists {
		return nil, fmt.Errorf("leagues data not loaded")
	}

	return ProcessLeaguesData(rawData)
}
