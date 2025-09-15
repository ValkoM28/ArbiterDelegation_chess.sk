package pdf

import (
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form"
)

// ListFillableFields extracts and lists all fillable form fields from a PDF document
func ListFillableFields(pdfPath string) error {
	// Read the PDF file into a context
	ctx, err := api.ReadContextFile(pdfPath)
	if err != nil {
		return fmt.Errorf("error reading PDF file: %v", err)
	}

	// List form fields
	fields, err := form.ListFormFields(ctx)
	if err != nil {
		return fmt.Errorf("error listing form fields: %v", err)
	}

	// Process the form fields
	for _, field := range fields {
		fmt.Printf("Field: %+v\n", field)
	}

	return nil
}

// FillForm fills the PDF form with provided data
func FillForm(pdfPath string, data map[string]string) (string, error) {
	// This is a placeholder - you'll need to implement the actual form filling logic
	// using the pdfcpu library
	return "", fmt.Errorf("FillForm not implemented yet")
}
