// Package pdf provides functionality for generating and filling PDF forms for chess arbiter delegations.
// It handles PDF form filling, data mapping, validation, and file generation.
package pdf

import (
	"fmt"
	"os"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form"
)

// FillForm fills a PDF form with the provided data and saves it to a new file.
// It reads the PDF template, fills in the form fields with the provided data map,
// and saves the result to a new file with a unique name.
// Returns the path to the filled PDF file or an error if the operation fails.
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

	// Generate unique output filename with UUID
	outputPath := fmt.Sprintf("assets/results/%s_%s_%s.pdf", data["text_1nzhs"], data["text_2qqiu"], uuid.New().String()[:8])

	// DEBUG: Log the generated filename
	fmt.Printf("DEBUG GENERATOR: Generated filename: %s\n", outputPath)

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

// generateSinglePDF generates a single PDF from PDFData
func generateSinglePDF(pdfData data.PDFData, templatePath string, index int) (string, error) {
	// Validate the PDF data
	if err := validatePDFData(pdfData); err != nil {
		return "", fmt.Errorf("validation failed for item %d: %v", index, err)
	}

	// Map data to fields using the same logic as the original
	fieldData := MapDataToFields(pdfData, DefaultFieldMapping)

	// Generate PDF for this data - same as original line 161
	outputPath, err := FillForm(templatePath, fieldData)
	if err != nil {
		return "", fmt.Errorf("error generating PDF for item %d: %v", index, err)
	}

	return outputPath, nil
}

// GeneratePDFsFromDelegateArbiters generates PDF files for each delegate-arbiter data
func GeneratePDFsFromDelegateArbiters(pdfDataArray []data.PDFData, templatePath string) ([]string, error) {
	if err := validateTemplate(templatePath); err != nil {
		return nil, err
	}

	// Process each PDF data item
	var generatedFiles []string
	for i, pdfData := range pdfDataArray {
		filePath, err := generateSinglePDF(pdfData, templatePath, i)
		if err != nil {
			return nil, err
		}

		generatedFiles = append(generatedFiles, filePath)
		fmt.Printf("Generated PDF: %s\n", filePath)
	}

	return generatedFiles, nil
}
