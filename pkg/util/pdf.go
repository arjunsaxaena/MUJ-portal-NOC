package util

import (
	"MUJ_automated_mail_generation/pkg/model"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func CreateNocPdf(submission model.StudentSubmission) (string, error) {
	// Log the submission details for debugging
	fmt.Printf("Submission data: %+v\n", submission)

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10) // Set margins
	pdf.AddPage()

	// Add title
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "No Objection Certificate (NOC)", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Set font for body text
	pdf.SetFont("Arial", "", 12)

	// Add the body content
	body := fmt.Sprintf(`This is to certify that %s, a student of the %s Department, Section %s, has been offered a placement at %s.

The details of the placement offer are as follows:
Company Name: %s
Offer Type: %s
Package: %s
Stipend: %s
Internship Dates: %s to %s

This certificate is issued to facilitate the necessary processes.

Best regards,
The University Placement Office`,
		submission.Name,
		submission.Department,
		submission.Section,
		submission.CompanyName,
		submission.CompanyName,
		submission.OfferType,
		submission.PackagePPO,
		submission.StipendAmount,
		submission.InternshipStartDate,
		submission.InternshipEndDate)

	// Log body content for debugging
	fmt.Printf("Body content: \n%s\n", body)

	pdf.MultiCell(0, 10, body, "", "L", false)
	pdf.Ln(10)

	// Define output directory
	outputDir := "generated_pdfs"
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Save the PDF to a file
	fileName := fmt.Sprintf("NOC_%s.pdf", submission.RegistrationNumber)
	outputPath := filepath.Join(outputDir, fileName)

	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to save PDF: %v", err)
	}

	fmt.Printf("PDF successfully saved at: %s\n", outputPath)
	return outputPath, nil
}
