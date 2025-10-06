package pdf

import (
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

// FieldMapping defines the mapping between data fields and PDF form fields
type FieldMapping struct {
	ArbiterFirstName string
	ArbiterLastName  string
	ArbiterPlayerID  string
	LeagueAndYear    string
	HomeTeam         string
	GuestTeam        string
	DateTime         string
	Address          string
	DirectorContact  string
	ContactPerson    string
}

// DefaultFieldMapping provides the standard field mapping for the PDF template
var DefaultFieldMapping = FieldMapping{
	ArbiterFirstName: "text_2qqiu",
	ArbiterLastName:  "text_1nzhs",
	ArbiterPlayerID:  "text_3bxac",
	LeagueAndYear:    "text_4ab",
	GuestTeam:        "text_7ubi",  // Fixed: GuestTeam goes to date field
	DateTime:         "text_5ohxu", // Fixed: DateTime goes to guest team field
	HomeTeam:         "text_6wdxk",
	Address:          "text_8hipe",
	DirectorContact:  "text_9lqnq",
	ContactPerson:    "text_10cjzk",
}

// MapDataToFields converts PDFData to the field mapping format used by the PDF form
// This preserves the exact same logic as the original GeneratePDFsFromDelegateArbiters function
func MapDataToFields(pdfData data.PDFData, mapping FieldMapping) map[string]string {
	stringData := make(map[string]string)

	// Extract arbiter data - same logic as original lines 104-114
	stringData[mapping.ArbiterFirstName] = pdfData.Arbiter.FirstName
	stringData[mapping.ArbiterLastName] = pdfData.Arbiter.LastName
	stringData[mapping.ArbiterPlayerID] = pdfData.Arbiter.PlayerID

	// Extract league data - same logic as original lines 117-130
	leagueAndYear := pdfData.League.Name
	if pdfData.League.Year != "" {
		if leagueAndYear != "" {
			leagueAndYear += " " + pdfData.League.Year
		} else {
			leagueAndYear = pdfData.League.Year
		}
	}
	stringData[mapping.LeagueAndYear] = leagueAndYear

	// Extract match data - same logic as original lines 133-146
	stringData[mapping.HomeTeam] = pdfData.Match.HomeTeam
	stringData[mapping.GuestTeam] = pdfData.Match.GuestTeam
	stringData[mapping.DateTime] = pdfData.Match.DateTime
	stringData[mapping.Address] = pdfData.Match.Address

	// Extract director data - same logic as original lines 149-153
	stringData[mapping.DirectorContact] = pdfData.Director.Contact

	// Extract contact person - same logic as original lines 156-158
	stringData[mapping.ContactPerson] = pdfData.ContactPerson

	return stringData
}
