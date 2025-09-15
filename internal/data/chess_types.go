package data

import (
	"encoding/json"
	"fmt"
)

// Arbiter represents an arbiter from the chess.sk API
type Arbiter struct {
	ArbiterId    string `json:"ArbiterId"`
	PlayerId     string `json:"PlayerId"`
	FideId       string `json:"FideId"`
	LastName     string `json:"LastName"`
	FirstName    string `json:"FirstName"`
	ValidTo      string `json:"ValidTo"`
	Licencia     string `json:"Licencia"`
	KlubId       string `json:"KlubId"`
	KlubName     string `json:"KlubName"`
	IsActive     bool   `json:"IsActive"`
	ArbiterLevel string `json:"ArbiterLevel"`
}

// League represents a league from the chess.sk API
type League struct {
	LeagueId          string `json:"leagueId"`
	SaisonName        string `json:"saisonName"`
	LeagueName        string `json:"leagueName"`
	ChessResultsLink  string `json:"chessResultsLink"`
	DirectorId        string `json:"directorId"`
	DirectorSurname   string `json:"directorSurname"`
	DirectorFirstName string `json:"directorFirstName"`
	DirectorEmail     string `json:"directorEmail"`
}

// ProcessArbitersData processes raw API data and extracts arbiters
func ProcessArbitersData(rawData interface{}) ([]Arbiter, error) {
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
	var arbiters []Arbiter
	if err := json.Unmarshal(jsonData, &arbiters); err != nil {
		return nil, fmt.Errorf("error unmarshaling arbiters data: %v", err)
	}

	return arbiters, nil
}

// ProcessLeaguesData processes raw API data and extracts leagues
func ProcessLeaguesData(rawData interface{}) ([]League, error) {
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
		if arbiter.ArbiterId == fmt.Sprintf("%d", arbiterID) {
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
		if league.LeagueId == fmt.Sprintf("%d", leagueID) {
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
