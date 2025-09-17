package pdf

import (
	"fmt"
	"os"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

// validateTemplate checks if the PDF template file exists and is accessible
func validateTemplate(templatePath string) error {
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("PDF template not found: %s", templatePath)
	}
	return nil
}

// validatePDFData checks if the PDFData has all required fields
func validatePDFData(pdfData data.PDFData) error {
	return pdfData.Validate()
}
