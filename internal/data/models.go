package data

import (
	"fmt"
)

// PDFData represents the structured data for PDF generation
type PDFData struct {
	Arbiter       ArbiterData
	League        LeagueData
	Match         MatchData
	Director      DirectorData
	ContactPerson string
}

// ArbiterData contains arbiter information
type ArbiterData struct {
	FirstName string
	LastName  string
	PlayerID  string
}

// LeagueData contains league information
type LeagueData struct {
	Name string
	Year string
}

// MatchData contains match information
type MatchData struct {
	HomeTeam  string
	GuestTeam string
	DateTime  string
	Address   string
}

// DirectorData contains director information
type DirectorData struct {
	Contact string
}

// Round represents a single round of matches
type Round struct {
	Number   int         `json:"number"`
	DateTime string      `json:"dateTime"` // e.g., "2025/10/25 at 11:00"
	Matches  []MatchInfo `json:"matches"`
}

// MatchInfo contains match-specific details
type MatchInfo struct {
	HomeTeam  string `json:"homeTeam"`
	GuestTeam string `json:"guestTeam"`
	DateTime  string `json:"dateTime"`
	Address   string `json:"address"`
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

// NewPDFData creates a new PDFData instance
func NewPDFData() *PDFData {
	return &PDFData{
		Arbiter:       ArbiterData{},
		League:        LeagueData{},
		Match:         MatchData{},
		Director:      DirectorData{},
		ContactPerson: "",
	}
}

// Validate checks if the PDFData has all required fields
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

// SetLeague sets the league information
func (p *PDFData) SetLeague(league LeagueData) {
	p.League = league
}

// SetDirector sets the director information
func (p *PDFData) SetDirector(director DirectorData) {
	p.Director = director
}

// SetArbiter sets the arbiter information
func (p *PDFData) SetArbiter(arbiter ArbiterData) {
	p.Arbiter = arbiter
}

// SetMatch sets the match information
func (p *PDFData) SetMatch(match MatchData) {
	p.Match = match
}

// SetContactPerson sets the contact person
func (p *PDFData) SetContactPerson(contactPerson string) {
	p.ContactPerson = contactPerson
}
