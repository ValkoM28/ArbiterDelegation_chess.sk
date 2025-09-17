package chess

import (
	"fmt"
	"net/url"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

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
func LoadArbiters(storage *data.SessionData) error {
	// Load arbiters data from your real API with hardcoded active status parameter
	// TODO: Remove client-side filtering when chess.sk API properly supports status=active parameter
	arbitersURL := buildURLWithParams("https://chess.sk/api/matrika.php/v1/arbiters", map[string]string{
		"status": "active", // Currently ignored by API, but kept for when it gets fixed
	})

	err := storage.LoadData("arbiters", arbitersURL)
	if err != nil {
		return fmt.Errorf("failed to load arbiters: %v", err)
	}

	// TEMPORARY: Client-side filtering for active arbiters until chess.sk API supports status=active
	arbitersData, exists := storage.Get("arbiters")
	if exists {
		filteredArbiters, err := filterActiveArbiters(arbitersData)
		if err != nil {
			return fmt.Errorf("failed to filter arbiters: %v", err)
		}
		storage.Set("arbiters", filteredArbiters)
	}

	return nil
}

// GetArbiters returns all arbiters from storage
func GetArbiters(storage *data.SessionData) ([]data.Arbiter, error) {
	return storage.GetAllArbiters()
}

// GetArbiterByID returns a specific arbiter by ID
func GetArbiterByID(storage *data.SessionData, arbiterID int) (*data.Arbiter, error) {
	return storage.GetArbiterByID(arbiterID)
}
