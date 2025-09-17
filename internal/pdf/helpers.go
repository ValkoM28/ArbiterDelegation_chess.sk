package pdf

import (
	"fmt"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

// PreparePDFDataFromArbiterAndLeague creates PDF data from arbiter and league
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

// fromArbiter converts a chess Arbiter to ArbiterData
func fromArbiter(arbiter data.Arbiter) data.ArbiterData {
	return data.ArbiterData{
		FirstName: arbiter.FirstName,
		LastName:  arbiter.LastName,
		PlayerID:  arbiter.PlayerId,
	}
}

// fromLeague converts a chess League to LeagueData and DirectorData
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

