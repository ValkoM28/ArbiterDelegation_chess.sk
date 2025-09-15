package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// SessionData represents our simple in-memory data storage for the session
// Much simpler - no timestamps, no expiration, just load once and use
type SessionData struct {
	data  map[string]interface{} // The actual data
	mutex sync.RWMutex           // For thread safety
}

// NewSessionData creates a new SessionData instance
func NewSessionData() *SessionData {
	return &SessionData{
		data: make(map[string]interface{}),
	}
}

// Get retrieves data from session storage
func (sd *SessionData) Get(key string) (interface{}, bool) {
	sd.mutex.RLock()         // Lock for reading
	defer sd.mutex.RUnlock() // Unlock when function exits

	data, exists := sd.data[key]
	return data, exists
}

// Set stores data in session storage
func (sd *SessionData) Set(key string, value interface{}) {
	sd.mutex.Lock()         // Lock for writing
	defer sd.mutex.Unlock() // Unlock when function exits

	sd.data[key] = value
}

// LoadData loads data from external API and stores it
func (sd *SessionData) LoadData(key string, url string) error {
	// This will call our HTTP client function
	data, err := fetchFromAPI(url)
	if err != nil {
		return err
	}

	// Store the data
	sd.Set(key, data)
	return nil
}

// Clear removes all data from session storage
func (sd *SessionData) Clear() {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	sd.data = make(map[string]interface{})
}

// HasData checks if we have data for a key
func (sd *SessionData) HasData(key string) bool {
	sd.mutex.RLock()
	defer sd.mutex.RUnlock()

	_, exists := sd.data[key]
	return exists
}

// fetchFromAPI makes an HTTP request to get data
func fetchFromAPI(url string) (map[string]interface{}, error) {
	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the HTTP GET request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse the JSON response - it could be an array or an object
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Convert to map format for consistency
	resultMap := make(map[string]interface{})
	resultMap["data"] = result

	return resultMap, nil
}
