package util

import (
	"MUJ_automated_mail_generation/pkg/model"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func CreateNocPdf(submission model.StudentSubmission) (string, error) {
	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font for the document
	pdf.SetFont("Arial", "B", 16)

	// Title of the document (NOC)
	pdf.Cell(190, 10, "No Objection Certificate (NOC)")

	// Set the font for the body text
	pdf.SetFont("Arial", "", 12)

	// Add the body of the NOC
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

	// Add the body content (use MultiCell for multiline text)
	// Here, we're providing values for border, align, and fill
	pdf.MultiCell(0, 10, body, "", "", false)

	// Save the PDF to a file
	fileName := fmt.Sprintf("NOC_%s.pdf", submission.RegistrationNumber)
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// Something might be wrong here

// NOC was saved in project root
