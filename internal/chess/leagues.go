package chess

import (
	"fmt"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

// LoadLeagues loads leagues data from the API and stores it in storage
func LoadLeagues(storage *data.SessionData, seasonStartYear string) error {
	// Load leagues data from your real API with season parameter
	leaguesURL := fmt.Sprintf("https://chess.sk/api/leagues.php/v1/leagues?saisonStartYear=%s", seasonStartYear)
	err := storage.LoadData("leagues", leaguesURL)
	if err != nil {
		return fmt.Errorf("failed to load leagues: %v", err)
	}

	return nil
}

// GetLeagues returns all leagues from storage
func GetLeagues(storage *data.SessionData) ([]data.League, error) {
	return storage.GetAllLeagues()
}

// GetLeagueByID returns a specific league by ID
func GetLeagueByID(storage *data.SessionData, leagueID int) (*data.League, error) {
	return storage.GetLeagueByID(leagueID)
}
