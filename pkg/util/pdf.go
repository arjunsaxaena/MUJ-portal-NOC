package util

import (
	"MUJ_AMG/pkg/model"
	"bytes"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func CreateNocPdf(submission model.StudentSubmission) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	uploadsDir := filepath.Join("../uploads", "noc")

	letterheadPath := "../pkg/images/muj_header.png"
	resolvedLetterheadPath, err := filepath.Abs(letterheadPath)
	if err != nil {
		fmt.Printf("Error resolving path for letterhead: %v\n", err)
		return "", fmt.Errorf("failed to resolve path for letterhead: %v", err)
	}

	imageWidth := 100.0
	imageHeight := 50.0
	pdf.ImageOptions(resolvedLetterheadPath, 0, 10, imageWidth, imageHeight, false, gofpdf.ImageOptions{}, 0, "")
	pdf.Ln(imageHeight + 3)

	pdf.SetFont("Arial", "", 10)
	nocText := fmt.Sprintf("MUJ/FoSTA/DCSE/2024/8H/%s", submission.RegistrationNumber)
	pdf.CellFormat(95, 10, nocText, "", 0, "L", false, 0, "") // on left side

	currentDate := time.Now().Format("02-Jan-2006")

	pdf.CellFormat(0, 10, currentDate, "", 1, "R", false, 0, "") // date on right side

	pdf.Ln(10)

	pdf.SetFont("Arial", "BU", 14) // underlined and bold
	pdf.CellFormat(0, 10, "To Whomsoever It May Concern", "", 1, "C", false, 0, "")

	pdf.Ln(10)

	startDate, err := time.Parse(time.RFC3339, submission.InternshipStartDate)
	if err != nil {
		return "", fmt.Errorf("invalid internship start date format: %v", err)
	}
	endDate, err := time.Parse(time.RFC3339, submission.InternshipEndDate)
	if err != nil {
		return "", fmt.Errorf("invalid internship end date format: %v", err)
	}
	startDateFormatted := startDate.Format("02-Jan-2006")
	endDateFormatted := endDate.Format("02-Jan-2006")

	pdf.SetFont("Arial", "", 10)
	title := "Mr."
	if submission.Gender == "Female" {
		title = "Ms."
	}
	body := fmt.Sprintf(`Sub: Recommendation for %s %s carrying out internship cum project in your esteemed Organization

Dear Sir/Madam,

This is to certify that %s %s (Reg No. %s) is a student of Manipal University Jaipur, India, studying in the %s semester of the four-year B.Tech Degree Program in the Department of %s, Section %s.

This recommendation is issued with reference to the application for an internship/project in your esteemed organization for a duration from %s to %s.

This Internship/Project would add value to the academic career of the student. So I request you to kindly allow our student to undergo Internship/Project at your organization.

Manipal University Jaipur has no objection for %s %s in doing an internship at your organization and has been advised to abide by both MUJ's and the interning organization's ethics/rules/regulations/values and work culture without compromising on integrity and self-discipline.

Thanking you.

Yours sincerely,`,
		title,
		submission.Name,
		title,
		submission.Name,
		submission.RegistrationNumber,
		submission.Semester,
		submission.Department,
		submission.Section,
		startDateFormatted,
		endDateFormatted,
		title,
		submission.Name,
	)

	pdf.MultiCell(0, 6, body, "", "L", false)

	pdf.Ln(4)

	pdf.SetFont("Arial", "", 10)
	footer := `Head of the Department & Professor (CSE)
Department of Computer Science & Engineering | FoSTA
Manipal University Jaipur
Dehmi Kalan, Off Jaipur-Ajmer Expressway, Jaipur - 303007, Rajasthan, India
chaudhary.neha@jaipur.manipal.edu
Phone: +91 141 3999100 (Extn:768) | Mobile: +91 9785500056`

	pdf.SetFont("Arial", "B", 11)
	pdf.MultiCell(0, 5, "Dr. Neha Chaudhary", "", "L", false)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, footer, "", "L", false)

	pdf.Ln(6)

	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 6, "This is a system-generated PDF.", "", 1, "C", false, 0, "")

	fileName := fmt.Sprintf("NOC_%s.pdf", submission.RegistrationNumber)
	localFilePath := filepath.Join(uploadsDir, fileName)

	err = pdf.OutputFileAndClose(localFilePath)
	if err != nil {
		fmt.Printf("Error saving PDF: %v\n", err)
		return "", fmt.Errorf("failed to save NOC PDF: %v", err)
	}

	fmt.Printf("PDF successfully saved at: %s\n", localFilePath)
	return localFilePath, nil
}

func CreateGenericNocPdf(submission model.StudentSubmission) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	pdf.SetFont("Arial", "", 10)
	nocText := fmt.Sprintf("MUJ/FoSTA/DCSE/2024/8H/%s", submission.RegistrationNumber)
	pdf.CellFormat(0, 10, nocText, "", 1, "L", false, 0, "")

	pdf.Ln(10)

	pdf.SetFont("Arial", "BU", 14)
	pdf.CellFormat(0, 10, "To Whomsoever It May Concern", "", 1, "C", false, 0, "")

	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	title := "Mr."
	if submission.Gender == "Female" {
		title = "Ms."
	}

	body := fmt.Sprintf(`This is to certify that %s %s (Reg No. %s) is a student of Manipal University Jaipur, India, studying in the %s semester of the four-year B.Tech Degree Program in the Department of %s, Section %s.

Manipal University Jaipur has no objection to %s %s undertaking professional engagements and has been advised to abide by ethical standards and work culture.

This document is issued upon request and holds no liability for further obligations.

Yours sincerely,`,
		title,
		submission.Name,
		submission.RegistrationNumber,
		submission.Semester,
		submission.Department,
		submission.Section,
		title,
		submission.Name,
	)

	pdf.MultiCell(0, 6, body, "", "L", false)

	pdf.Ln(6)

	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 6, "This is a system-generated PDF.", "", 1, "C", false, 0, "")

	var pdfBuffer bytes.Buffer
	err := pdf.Output(&pdfBuffer)
	if err != nil {
		fmt.Printf("Error generating PDF: %v\n", err)
		return "", fmt.Errorf("failed to generate PDF: %v", err)
	}

	uploadsDir := filepath.Join("../uploads", "noc")
	fileName := fmt.Sprintf("NOC_%s.pdf", submission.RegistrationNumber)
	localFilePath := filepath.Join(uploadsDir, fileName)

	err = pdf.OutputFileAndClose(localFilePath)
	if err != nil {
		fmt.Printf("Error saving PDF: %v\n", err)
		return "", fmt.Errorf("failed to save NOC PDF: %v", err)
	}

	fmt.Printf("PDF successfully saved at: %s\n", localFilePath)
	return localFilePath, nil
}
