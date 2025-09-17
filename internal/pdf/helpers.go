// Package pdf provides functionality for generating and filling PDF forms for chess arbiter delegations.
// It handles PDF form filling, data mapping, validation, and file generation.
package pdf

import (
	"fmt"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

// PreparePDFDataFromArbiterAndLeague creates a PDFData structure from arbiter and league information.
// It converts the raw arbiter and league data into the format needed for PDF generation.
// The function sets up placeholder data for match details and contact person.
// Returns a pointer to a new PDFData instance ready for PDF generation.
func PreparePDFDataFromArbiterAndLeague(arbiter *data.Arbiter, league *data.League) *data.PDFData {
	pdfData := data.NewPDFData()

	// Set league and director from selected league
	leagueInfo, directorInfo := fromLeague(*league)
	pdfData.SetLeague(leagueInfo)
	pdfData.SetDirector(directorInfo)

	// Set arbiter
	arbiterInfo := fromArbiter(*arbiter)
	pdfData.SetArbiter(arbiterInfo)

	// Set match details (you'll provide logic later)
	pdfData.SetMatch(data.MatchData{
		HomeTeam:  "Team A",
		GuestTeam: "Team B",
		DateTime:  "2024-10-15 14:00",
		Address:   "Chess Center Bratislava",
	})

	// Set contact person (you'll provide logic later)
	pdfData.SetContactPerson("John Doe")

	return pdfData
}

// fromArbiter converts a chess Arbiter to ArbiterData for PDF generation.
// It extracts the necessary fields from the API response format to the PDF data format.
// Returns an ArbiterData struct with the arbiter's name and player ID.
func fromArbiter(arbiter data.Arbiter) data.ArbiterData {
	return data.ArbiterData{
		FirstName: arbiter.FirstName,
		LastName:  arbiter.LastName,
		PlayerID:  arbiter.PlayerId,
	}
}

// fromLeague converts a chess League to LeagueData and DirectorData for PDF generation.
// It extracts league information and director contact details from the API response format.
// The director contact combines name and email into a single string.
// Returns both LeagueData and DirectorData structs.
func fromLeague(league data.League) (data.LeagueData, data.DirectorData) {
	leagueInfo := data.LeagueData{
		Name: league.LeagueName,
		Year: league.SaisonName,
	}

	// Combine director name and email into a single contact string
	directorContact := fmt.Sprintf("%s %s", league.DirectorFirstName, league.DirectorSurname)
	if league.DirectorEmail != "" {
		directorContact += fmt.Sprintf(" (%s)", league.DirectorEmail)
	}

	directorInfo := data.DirectorData{
		Contact: directorContact,
	}

	return leagueInfo, directorInfo
}
