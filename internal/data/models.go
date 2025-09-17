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
	if p.Arbiter.FirstName == "" {
		return fmt.Errorf("arbiter first name is required")
	}
	if p.Arbiter.LastName == "" {
		return fmt.Errorf("arbiter last name is required")
	}
	if p.Arbiter.PlayerID == "" {
		return fmt.Errorf("arbiter player ID is required")
	}
	if p.League.Name == "" {
		return fmt.Errorf("league name is required")
	}
	if p.Director.Contact == "" {
		return fmt.Errorf("director contact is required")
	}
	return nil
}
