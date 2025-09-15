package data

import (
	"encoding/json"
	"fmt"
)

// Arbiter represents an arbiter from the chess.sk API
type Arbiter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Category string `json:"category"`
	License  string `json:"license"`
	Active   bool   `json:"active"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	// Add more fields as needed based on your API response
}

// League represents a league from the chess.sk API
type League struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Season    string `json:"season"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Status    string `json:"status"`
	Location  string `json:"location"`
	// Add more fields as needed based on your API response
}

// ProcessArbitersData processes raw API data and extracts arbiters
func ProcessArbitersData(rawData interface{}) ([]Arbiter, error) {
	// Convert raw data to JSON bytes
	jsonData, err := json.Marshal(rawData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling raw data: %v", err)
	}

	// Parse into structured data
	var arbiters []Arbiter
	if err := json.Unmarshal(jsonData, &arbiters); err != nil {
		return nil, fmt.Errorf("error unmarshaling arbiters data: %v", err)
	}

	return arbiters, nil
}

// ProcessLeaguesData processes raw API data and extracts leagues
func ProcessLeaguesData(rawData interface{}) ([]League, error) {
	// Convert raw data to JSON bytes
	jsonData, err := json.Marshal(rawData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling raw data: %v", err)
	}

	// Parse into structured data
	var leagues []League
	if err := json.Unmarshal(jsonData, &leagues); err != nil {
		return nil, fmt.Errorf("error unmarshaling leagues data: %v", err)
	}

	return leagues, nil
}

// GetArbiterByID finds an arbiter by ID from the loaded data
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
		if arbiter.ID == arbiterID {
			return &arbiter, nil
		}
	}

	return nil, fmt.Errorf("arbiter with ID %d not found", arbiterID)
}

// GetLeagueByID finds a league by ID from the loaded data
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
		if league.ID == leagueID {
			return &league, nil
		}
	}

	return nil, fmt.Errorf("league with ID %d not found", leagueID)
}

// GetAllArbiters returns all loaded arbiters
func (sd *SessionData) GetAllArbiters() ([]Arbiter, error) {
	rawData, exists := sd.Get("arbiters")
	if !exists {
		return nil, fmt.Errorf("arbiters data not loaded")
	}

	return ProcessArbitersData(rawData)
}

// GetAllLeagues returns all loaded leagues
func (sd *SessionData) GetAllLeagues() ([]League, error) {
	rawData, exists := sd.Get("leagues")
	if !exists {
		return nil, fmt.Errorf("leagues data not loaded")
	}

	return ProcessLeaguesData(rawData)
}
