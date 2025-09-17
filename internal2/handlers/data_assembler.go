package handlers

import (
	"fmt"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

// PreparePDFData creates PDF data from arbiter and league parameters
func PreparePDFData(arbiter *data.Arbiter, league *data.League) *data.PDFData {
	pdfData := data.NewPDFData()

	// Set league and director from selected league
	leagueInfo, directorInfo := data.FromLeague(*league)
	pdfData.SetLeague(leagueInfo)
	pdfData.SetDirector(directorInfo)

	// Set arbiter
	arbiterInfo := data.FromArbiter(*arbiter)
	pdfData.SetArbiter(arbiterInfo)

	// Set match details (you'll provide logic later)
	pdfData.SetMatch(data.MatchInfo{
		HomeTeam:  "Team A",
		GuestTeam: "Team B",
		DateTime:  "2024-10-15 14:00",
		Address:   "Chess Center Bratislava",
	})

	// Set contact person (you'll provide logic later)
	pdfData.SetContactPerson("John Doe")

	return pdfData
}

// PrintPDFData prints the PDF data in a readable format
func PrintPDFData(pdfData *data.PDFData) {
	fmt.Println("=== PDF DATA ===")
	fmt.Printf("League: %s (%s)\n", pdfData.League.Name, pdfData.League.Year)
	fmt.Printf("Director: %s\n", pdfData.Director.Contact)
	fmt.Printf("Arbiter: %s %s (ID: %s)\n", pdfData.Arbiter.FirstName, pdfData.Arbiter.LastName, pdfData.Arbiter.PlayerID)
	fmt.Printf("Match: %s vs %s\n", pdfData.Match.HomeTeam, pdfData.Match.GuestTeam)
	fmt.Printf("Date/Time: %s\n", pdfData.Match.DateTime)
	fmt.Printf("Address: %s\n", pdfData.Match.Address)
	fmt.Printf("Contact Person: %s\n", pdfData.ContactPerson)
	fmt.Println("================")
}
