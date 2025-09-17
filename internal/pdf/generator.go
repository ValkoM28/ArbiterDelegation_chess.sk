package pdf

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form"
)

// ListFillableFields extracts and lists all fillable form fields from a PDF document
func ListFillableFields(pdfPath string) ([]string, error) {
	// Read the PDF file into a context
	ctx, err := api.ReadContextFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("error reading PDF file: %v", err)
	}

	// List form fields
	fields, err := form.ListFormFields(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing form fields: %v", err)
	}

	// Extract field names - fields is a slice of strings
	var fieldNames []string
	for _, fieldName := range fields {
		if fieldName != "" {
			fieldNames = append(fieldNames, fieldName)
			fmt.Printf("Field: %s\n", fieldName)
		}
	}

	return fieldNames, nil
}

// FillForm fills the PDF form with provided data
func FillForm(pdfPath string, data map[string]string) (string, error) {
	// Read the PDF file into a context
	ctx, err := api.ReadContextFile(pdfPath)
	if err != nil {
		return "", fmt.Errorf("error reading PDF file: %v", err)
	}

	// Create a field processor function
	fieldProcessor := func(id string, name string, fieldType form.FieldType, format form.DataFormat) ([]string, bool, bool) {
		if value, exists := data[name]; exists {
			return []string{value}, true, true
		}
		return []string{}, false, false
	}

	// Fill the form fields using the correct API
	_, _, err = form.FillForm(ctx, fieldProcessor, nil, form.DataFormat(0))
	if err != nil {
		return "", fmt.Errorf("error filling form fields: %v", err)
	}

	// Generate unique output filename with microsecond precision
	outputPath := fmt.Sprintf("assets/results/%s_%s_%d.pdf", data["Text1"], data["Text2"], time.Now().Nanosecond())

	// Ensure the results directory exists
	resultsDir := "assets/results"
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create results directory: %v", err)
	}

	// Write the filled PDF
	err = api.WriteContextFile(ctx, outputPath)
	if err != nil {
		return "", fmt.Errorf("error writing filled PDF: %v", err)
	}

	return outputPath, nil
}

// GeneratePDFsFromDelegateArbiters generates PDF files for each delegate-arbiter data
func GeneratePDFsFromDelegateArbiters(pdfDataArray []interface{}, templatePath string) ([]string, error) {
	var generatedFiles []string

	// First, get the field names from the template
	fieldNames, err := ListFillableFields(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error getting field names: %v", err)
	}

	fmt.Printf("Available fields in PDF: %v\n", fieldNames)

	// Process each PDF data item
	for i, pdfData := range pdfDataArray {
		// Convert interface{} to map[string]interface{}
		dataMap, ok := pdfData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid data format for item %d", i)
		}

		// Map the data to PDF field names based on the required order:
		// arbiter firstName, arbiterLastName, arbiter playerId, League and Year, Home team, away team, date and time, address of the venue, director+contact, contactperson
		// PDF fields are: Text1, Text2, Text3, Text4, Text5, Text6, Text7, Text8, Text9, Text10
		stringData := make(map[string]string)

		// Extract arbiter data
		if arbiter, ok := dataMap["arbiter"].(map[string]interface{}); ok {
			if firstName, ok := arbiter["firstName"].(string); ok {
				stringData["Text1"] = firstName // arbiter firstName
			}
			if lastName, ok := arbiter["lastName"].(string); ok {
				stringData["Text2"] = lastName // arbiter lastName
			}
			if playerId, ok := arbiter["playerId"].(string); ok {
				stringData["Text3"] = playerId // arbiter playerId
			}
		}

		// Extract league data
		if league, ok := dataMap["league"].(map[string]interface{}); ok {
			leagueAndYear := ""
			if name, ok := league["name"].(string); ok {
				leagueAndYear = name
			}
			if year, ok := league["year"].(string); ok {
				if leagueAndYear != "" {
					leagueAndYear += " " + year
				} else {
					leagueAndYear = year
				}
			}
			stringData["Text4"] = leagueAndYear // League and Year
		}

		// Extract match data
		if match, ok := dataMap["match"].(map[string]interface{}); ok {
			if homeTeam, ok := match["homeTeam"].(string); ok {
				stringData["Text5"] = homeTeam // Home team
			}
			if guestTeam, ok := match["guestTeam"].(string); ok {
				stringData["Text6"] = guestTeam // Away team
			}
			if dateTime, ok := match["dateTime"].(string); ok {
				stringData["Text7"] = dateTime // Date and time
			}
			if address, ok := match["address"].(string); ok {
				stringData["Text8"] = address // Address of the venue
			}
		}

		// Extract director data
		if director, ok := dataMap["director"].(map[string]interface{}); ok {
			if contact, ok := director["contact"].(string); ok {
				stringData["Text9"] = contact // Director+contact
			}
		}

		// Extract contact person
		if contactPerson, ok := dataMap["contactPerson"].(string); ok {
			stringData["Text10"] = contactPerson // Contact person
		}

		// Generate PDF for this data
		outputPath, err := FillForm(templatePath, stringData)
		if err != nil {
			return nil, fmt.Errorf("error generating PDF for item %d: %v", i, err)
		}

		generatedFiles = append(generatedFiles, outputPath)
		fmt.Printf("Generated PDF: %s\n", outputPath)
	}

	return generatedFiles, nil
}

// CreateZipFromFiles creates a zip file containing all the provided PDF files
func CreateZipFromFiles(pdfFiles []string, zipName string) (string, error) {
	// Ensure the results directory exists
	resultsDir := "assets/results"
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create results directory: %v", err)
	}

	// Create zip file path
	zipPath := filepath.Join(resultsDir, zipName)

	// Create the zip file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add each PDF file to the zip
	for _, pdfFile := range pdfFiles {
		// Open the PDF file
		file, err := os.Open(pdfFile)
		if err != nil {
			return "", fmt.Errorf("failed to open PDF file %s: %v", pdfFile, err)
		}
		defer file.Close()

		// Get file info for the zip entry
		fileInfo, err := file.Stat()
		if err != nil {
			return "", fmt.Errorf("failed to get file info for %s: %v", pdfFile, err)
		}

		// Create a zip file header
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return "", fmt.Errorf("failed to create zip header for %s: %v", pdfFile, err)
		}

		// Set the name in the zip to just the filename (not the full path)
		header.Name = filepath.Base(pdfFile)
		header.Method = zip.Deflate

		// Create the zip file entry
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return "", fmt.Errorf("failed to create zip entry for %s: %v", pdfFile, err)
		}

		// Copy the file content to the zip entry
		_, err = io.Copy(writer, file)
		if err != nil {
			return "", fmt.Errorf("failed to copy file %s to zip: %v", pdfFile, err)
		}
	}

	return zipPath, nil
}

// GeneratePDFsAndZip generates PDF files and creates a zip file containing all of them
func GeneratePDFsAndZip(pdfDataArray []interface{}, templatePath string, zipName string) (string, error) {
	// Generate PDFs first
	generatedFiles, err := GeneratePDFsFromDelegateArbiters(pdfDataArray, templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDFs: %v", err)
	}

	// Create zip file
	zipPath, err := CreateZipFromFiles(generatedFiles, zipName)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %v", err)
	}

	// Clean up individual PDF files after creating zip
	for _, file := range generatedFiles {
		if err := os.Remove(file); err != nil {
			fmt.Printf("Warning: failed to remove temporary PDF file %s: %v\n", file, err)
		}
	}

	return zipPath, nil
}
