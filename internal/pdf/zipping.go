package pdf

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

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
func GeneratePDFsAndZip(pdfDataArray []data.PDFData, templatePath string, zipName string) (string, error) {
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
