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

// PDFData represents the complete data structure for PDF generation
type PDFData struct {
	League        LeagueInfo   `json:"league"`
	Director      DirectorInfo `json:"director"`
	Arbiter       ArbiterInfo  `json:"arbiter"`
	Match         MatchInfo    `json:"match"`
	ContactPerson string       `json:"contactPerson"`
}

// LeagueInfo contains league-specific information
type LeagueInfo struct {
	Name string `json:"name"`
	Year string `json:"year"`
}

// DirectorInfo contains director contact information (single string)
type DirectorInfo struct {
	Contact string `json:"contact"` // Combined name and contact info
}

// ArbiterInfo contains arbiter information for delegation
type ArbiterInfo struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PlayerID  string `json:"playerId"` // Using PlayerId instead of ArbiterId
}

// MatchInfo contains match-specific details
type MatchInfo struct {
	HomeTeam  string `json:"homeTeam"`
	GuestTeam string `json:"guestTeam"`
	DateTime  string `json:"dateTime"`
	Address   string `json:"address"`
}

// Round represents a single round of matches
type Round struct {
	Number   int         `json:"number"`
	DateTime string      `json:"dateTime"` // e.g., "2025/10/25 at 11:00"
	Matches  []MatchInfo `json:"matches"`
}

// TournamentData represents the complete tournament structure
type TournamentData struct {
	LeagueName string  `json:"leagueName"`
	Rounds     []Round `json:"rounds"`
}

// NewPDFData creates a new PDFData instance with default values
func NewPDFData() *PDFData {
	return &PDFData{
		League:        LeagueInfo{},
		Director:      DirectorInfo{},
		Arbiter:       ArbiterInfo{},
		Match:         MatchInfo{},
		ContactPerson: "",
	}
}

// SetLeague sets the league information
func (p *PDFData) SetLeague(league LeagueInfo) {
	p.League = league
}

// SetDirector sets the director information
func (p *PDFData) SetDirector(director DirectorInfo) {
	p.Director = director
}

// SetArbiter sets the arbiter information
func (p *PDFData) SetArbiter(arbiter ArbiterInfo) {
	p.Arbiter = arbiter
}

// SetMatch sets the match information
func (p *PDFData) SetMatch(match MatchInfo) {
	p.Match = match
}

// SetContactPerson sets the contact person
func (p *PDFData) SetContactPerson(contactPerson string) {
	p.ContactPerson = contactPerson
}

// ToMap converts PDFData to a map[string]string for PDF form filling
func (p *PDFData) ToMap() map[string]string {
	result := make(map[string]string)

	// Arbiter fields
	result["arbiter_first_name"] = p.Arbiter.FirstName
	result["arbiter_last_name"] = p.Arbiter.LastName
	result["arbiter_id"] = p.Arbiter.PlayerID

	// League fields
	result["league_name"] = p.League.Name
	result["league_year"] = p.League.Year

	// Match fields
	result["home_team"] = p.Match.HomeTeam
	result["guest_team"] = p.Match.GuestTeam
	result["date_time"] = p.Match.DateTime
	result["address"] = p.Match.Address

	// Director field (single combined string)
	result["league_director_contact"] = p.Director.Contact

	// Contact person
	result["contact_person"] = p.ContactPerson

	return result
}

// Validate checks if the PDFData has all required fields
func (p *PDFData) Validate() error {
	if p.League.Name == "" {
		return fmt.Errorf("league name is required")
	}
	if p.League.Year == "" {
		return fmt.Errorf("league year is required")
	}
	if p.Arbiter.FirstName == "" {
		return fmt.Errorf("arbiter first name is required")
	}
	if p.Arbiter.LastName == "" {
		return fmt.Errorf("arbiter last name is required")
	}
	if p.Arbiter.PlayerID == "" {
		return fmt.Errorf("arbiter player ID is required")
	}
	if p.Director.Contact == "" {
		return fmt.Errorf("league director contact is required")
	}
	return nil
}

// FromArbiter converts a chess Arbiter to ArbiterInfo
func FromArbiter(arbiter Arbiter) ArbiterInfo {
	return ArbiterInfo{
		FirstName: arbiter.FirstName,
		LastName:  arbiter.LastName,
		PlayerID:  arbiter.PlayerId, // Using PlayerId instead of ArbiterId
	}
}

// FromLeague converts a chess League to LeagueInfo and DirectorInfo
func FromLeague(league League) (LeagueInfo, DirectorInfo) {
	leagueInfo := LeagueInfo{
		Name: league.LeagueName,
		Year: league.SaisonName, // Assuming SaisonName contains the year
	}

	// Combine director name and email into a single contact string
	directorContact := fmt.Sprintf("%s %s", league.DirectorFirstName, league.DirectorSurname)
	if league.DirectorEmail != "" {
		directorContact += fmt.Sprintf(" (%s)", league.DirectorEmail)
	}

	directorInfo := DirectorInfo{
		Contact: directorContact,
	}

	return leagueInfo, directorInfo
}
