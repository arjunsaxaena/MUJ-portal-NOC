package util

import (
	"MUJ_automated_mail_generation/pkg/model"
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func CreateNocPdf(submission model.StudentSubmission, bucketName, keyPrefix string) (string, error) {
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

	// Write PDF data to a buffer
	var pdfBuffer bytes.Buffer
	err := pdf.Output(&pdfBuffer)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %v", err)
	}

	// Define the S3 object key
	fileName := fmt.Sprintf("NOC_%s.pdf", submission.RegistrationNumber)
	s3Key := fmt.Sprintf("%s/%s", keyPrefix, fileName)

	// Upload the PDF to S3
	err = UploadFileToS3(bucketName, &pdfBuffer, s3Key)
	if err != nil {
		return "", fmt.Errorf("failed to upload NOC PDF to S3: %v", err)
	}

	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, s3Key)
	fmt.Printf("PDF successfully uploaded to S3: %s\n", s3URL)

	return s3URL, nil
}
