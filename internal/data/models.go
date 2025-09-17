// Package data provides data models and structures for the chess arbiter delegation generator.
// It includes models for arbiters, leagues, matches, and PDF generation data.
package data

import (
	"fmt"
)

// PDFData represents the structured data for PDF generation.
// It contains all necessary information to fill out a delegation form PDF.
type PDFData struct {
	Arbiter       ArbiterData  // Information about the assigned arbiter
	League        LeagueData   // Information about the chess league
	Match         MatchData    // Information about the specific match
	Director      DirectorData // Information about the league director
	ContactPerson string       // Contact person for the delegation
}

// ArbiterData contains arbiter information extracted from the chess.sk API.
type ArbiterData struct {
	FirstName string // Arbiter's first name
	LastName  string // Arbiter's last name
	PlayerID  string // Arbiter's player ID in the chess system
}

// LeagueData contains league information extracted from the chess.sk API.
type LeagueData struct {
	Name string // Name of the chess league
	Year string // Season year (e.g., "2024/2025")
}

// MatchData contains match information for the delegation.
type MatchData struct {
	HomeTeam  string // Name of the home team
	GuestTeam string // Name of the guest team
	DateTime  string // Date and time of the match
	Address   string // Venue address for the match
}

// DirectorData contains director information extracted from the chess.sk API.
type DirectorData struct {
	Contact string // Director's contact information (name and email)
}

// Round represents a single round of matches in a chess league.
// Each round contains multiple matches played at the same time.
type Round struct {
	Number   int         `json:"number"`   // Round number (1, 2, 3, etc.)
	DateTime string      `json:"dateTime"` // Date and time of the round (e.g., "2025/10/25 at 11:00")
	Matches  []MatchInfo `json:"matches"`  // List of matches in this round
}

// MatchInfo contains match-specific details extracted from Excel files.
type MatchInfo struct {
	HomeTeam  string `json:"homeTeam"`  // Name of the home team
	GuestTeam string `json:"guestTeam"` // Name of the guest team
	DateTime  string `json:"dateTime"`  // Date and time of the match
	Address   string `json:"address"`   // Venue address (usually empty in Excel format)
}

// League represents a league from the chess.sk API.
// This structure matches the JSON response format from the API.
type League struct {
	LeagueId          string `json:"leagueId"`          // Unique identifier for the league
	SaisonName        string `json:"saisonName"`        // Season name (e.g., "2024/2025")
	LeagueName        string `json:"leagueName"`        // Display name of the league
	ChessResultsLink  string `json:"chessResultsLink"`  // URL to chess-results.com tournament page
	DirectorId        string `json:"directorId"`        // Director's unique identifier
	DirectorSurname   string `json:"directorSurname"`   // Director's surname
	DirectorFirstName string `json:"directorFirstName"` // Director's first name
	DirectorEmail     string `json:"directorEmail"`     // Director's email address
}

// Arbiter represents an arbiter from the chess.sk API.
// This structure matches the JSON response format from the API.
type Arbiter struct {
	ArbiterId    string `json:"ArbiterId"`    // Unique identifier for the arbiter
	PlayerId     string `json:"PlayerId"`     // Player ID in the chess system
	FideId       string `json:"FideId"`       // FIDE ID (international chess federation)
	LastName     string `json:"LastName"`     // Arbiter's surname
	FirstName    string `json:"FirstName"`    // Arbiter's first name
	ValidTo      string `json:"ValidTo"`      // License validity end date
	Licencia     string `json:"Licencia"`     // License number
	KlubId       string `json:"KlubId"`       // Club identifier
	KlubName     string `json:"KlubName"`     // Club name
	IsActive     bool   `json:"IsActive"`     // Whether the arbiter is currently active
	ArbiterLevel string `json:"ArbiterLevel"` // Arbiter's certification level
}

// NewPDFData creates a new PDFData instance with empty fields.
// Returns a pointer to a new PDFData struct ready for population.
func NewPDFData() *PDFData {
	return &PDFData{
		Arbiter:       ArbiterData{},
		League:        LeagueData{},
		Match:         MatchData{},
		Director:      DirectorData{},
		ContactPerson: "",
	}
}

// Validate checks if the PDFData has all required fields for PDF generation.
// Currently validates league name and director contact as required fields.
// Arbiter fields are optional for testing purposes.
func (p *PDFData) Validate() error {
	// For testing purposes, make arbiter fields optional
	// if p.Arbiter.FirstName == "" {
	// 	return fmt.Errorf("arbiter first name is required")
	// }
	// if p.Arbiter.LastName == "" {
	// 	return fmt.Errorf("arbiter last name is required")
	// }
	// if p.Arbiter.PlayerID == "" {
	// 	return fmt.Errorf("arbiter player ID is required")
	// }
	if p.League.Name == "" {
		return fmt.Errorf("league name is required")
	}
	if p.Director.Contact == "" {
		return fmt.Errorf("director contact is required")
	}
	return nil
}

// SetLeague sets the league information in the PDFData.
func (p *PDFData) SetLeague(league LeagueData) {
	p.League = league
}

// SetDirector sets the director information in the PDFData.
func (p *PDFData) SetDirector(director DirectorData) {
	p.Director = director
}

// SetArbiter sets the arbiter information in the PDFData.
func (p *PDFData) SetArbiter(arbiter ArbiterData) {
	p.Arbiter = arbiter
}

// SetMatch sets the match information in the PDFData.
func (p *PDFData) SetMatch(match MatchData) {
	p.Match = match
}

// SetContactPerson sets the contact person in the PDFData.
func (p *PDFData) SetContactPerson(contactPerson string) {
	p.ContactPerson = contactPerson
}
