package util

import (
	"MUJ_AMG/pkg/model"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func CreateNocPdf(submission model.StudentSubmission) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	uploadsDir := filepath.Join("../uploads", "NOC")

	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadsDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// letterheadPath := "../pkg/images/muj_header.png"
	// resolvedLetterheadPath, err := filepath.Abs(letterheadPath)
	// if err != nil {
	// 	fmt.Printf("Error resolving path for letterhead: %v\n", err)
	// 	return "", fmt.Errorf("failed to resolve path for letterhead: %v", err)
	// }

	// imageWidth := 100.0
	// imageHeight := 50.0
	// pdf.ImageOptions(resolvedLetterheadPath, 0, 10, imageWidth, imageHeight, false, gofpdf.ImageOptions{}, 0, "")
	// pdf.Ln(imageHeight + 3)

	pdf.Ln(20) // remove this if letterhead

	pdf.SetFont("Arial", "", 12)
	nocText := fmt.Sprintf("MUJ/FoSTA/DCSE/2025/%s%s/%s", submission.Semester, submission.Section, submission.RegistrationNumber[len(submission.RegistrationNumber)-4:])
	pdf.CellFormat(95, 10, nocText, "", 0, "L", false, 0, "") // left-aligned text

	currentDate := time.Now().Format("02-Jan-2006")
	pdf.CellFormat(0, 10, currentDate, "", 0, "R", false, 0, "") // right-aligned date on the same line
	pdf.Ln(20)

	pdf.SetFont("Arial", "BU", 16) // underlined and bold
	pdf.CellFormat(0, 10, "To Whomsoever It May Concern", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)

	pdf.Ln(15)

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

	pdf.SetFont("Arial", "", 12)
	title := "Mr."
	if submission.Gender == "Female" {
		title = "Ms."
	}

	pdf.SetFont("Arial", "", 12)
	pdf.Write(6, "Sub: Recommendation for ")
	pdf.SetFont("Arial", "B", 12)
	pdf.Write(6, fmt.Sprintf("%s %s", title, submission.Name))
	pdf.SetFont("Arial", "", 12)
	pdf.Write(6, " carrying out internship cum project in your esteemed Organization\n\n")

	pdf.Write(6, "Dear Sir/Madam,\n\n")

	pdf.Write(6, "This is to certify that ")
	pdf.SetFont("Arial", "B", 12)
	pdf.Write(6, fmt.Sprintf("%s %s", title, submission.Name))
	pdf.SetFont("Arial", "", 12)
	pdf.Write(6, fmt.Sprintf(" (Reg No. %s) is a student of Manipal University Jaipur, India, studying in the %s semester of the four-year B.Tech Degree Program in the Department of %s, Section %s.\n\n",
		submission.RegistrationNumber,
		submission.Semester,
		submission.Department,
		submission.Section))

	pdf.Write(6, "This recommendation is issued with reference to the application for an internship/project in your esteemed organization for a duration from ")
	pdf.SetFont("Arial", "B", 12) // bold font for start date
	pdf.Write(6, startDateFormatted)
	pdf.SetFont("Arial", "", 12) // reset to normal font
	pdf.Write(6, " to ")
	pdf.SetFont("Arial", "B", 12) // bold font for end date
	pdf.Write(6, endDateFormatted)
	pdf.SetFont("Arial", "", 12) // reset to normal font
	pdf.Write(6, ".\n\n")

	pdf.Write(6, "This Internship/Project would add value to the academic career of the student. So I request you to kindly allow our student to undergo Internship/Project at your organization.\n\n")

	pdf.Write(6, "Manipal University Jaipur has no objection for ")
	pdf.SetFont("Arial", "B", 12)
	pdf.Write(6, fmt.Sprintf("%s %s", title, submission.Name))
	pdf.SetFont("Arial", "", 12)
	pdf.Write(6, " in doing an internship at your organization and has been advised to abide by both MUJ's and the interning organization's ethics/rules/regulations/values and work culture without compromising on integrity and self-discipline.\n\n")

	pdf.Write(6, "Thanking you.\n\n")
	pdf.Write(6, "Yours sincerely,\n\n\n")

	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	footer := `Head of the Department & Professor (CSE)
Department of Computer Science & Engineering | FoSTA
Manipal University Jaipur
Dehmi Kalan, Off Jaipur-Ajmer Expressway, Jaipur - 303007, Rajasthan, India
chaudhary.neha@jaipur.manipal.edu
Phone: +91 141 3999100 (Extn:768) | Mobile: +91 9785500056`

	pdf.SetFont("Arial", "B", 12)
	pdf.MultiCell(0, 5, "Dr. Neha Chaudhary", "", "L", false)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 5, footer, "", "L", false)

	pdf.Ln(6)

	// pdf.SetFont("Arial", "I", 8)
	// pdf.CellFormat(0, 6, "This is a system-generated PDF.", "", 1, "C", false, 0, "")

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

	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	nocText := fmt.Sprintf("MUJ/FoSTA/DCSE/2025/%s%s/%s", submission.Semester, submission.Section, submission.RegistrationNumber[len(submission.RegistrationNumber)-4:])
	pdf.CellFormat(95, 10, nocText, "", 0, "L", false, 0, "") // left-aligned text

	currentDate := time.Now().Format("02-Jan-2006")
	pdf.CellFormat(0, 10, currentDate, "", 0, "R", false, 0, "") // right-aligned date on the same line
	pdf.Ln(20)

	pdf.SetFont("Arial", "BU", 16) // underlined and bold
	pdf.CellFormat(0, 10, "To Whomsoever It May Concern", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)

	pdf.Ln(15)

	startDate, err := time.Parse(time.RFC3339, submission.InternshipStartDate)
	if err != nil {
		return "", fmt.Errorf("invalid internship start date format: %v", err)
	}
	endDate, err := time.Parse(time.RFC3339, submission.InternshipEndDate)
	if err != nil {
		return "", fmt.Errorf("invalid internship end date format: %v", err)
	}

	startDateStr := startDate.Format("02 Jan 2006")
	endDateStr := endDate.Format("02 Jan 2006")

	// pdf.Ln(20)

	semesterInt, err := strconv.Atoi(submission.Semester)
	if err != nil {
		return "", fmt.Errorf("invalid semester format: %v", err)
	}

	year := 1
	switch semesterInt {
	case 3, 4:
		year = 2
	case 5, 6:
		year = 3
	case 7, 8:
		year = 4
	}

	pdf.SetFont("Arial", "", 12)
	title := "Mr."
	if submission.Gender == "Female" {
		title = "Ms."
	}

	///////////////////////////////////////////////////////////////////////////////////////////////

	pdf.SetFont("Arial", "B", 12)
	pdf.Write(6, fmt.Sprintf(" %s ", title))

	pdf.SetFont("Arial", "B", 12)
	pdf.Write(6, submission.Name)

	pdf.SetFont("Arial", "B", 12)
	pdf.Write(6, ", Reg No.- ")

	pdf.SetFont("Arial", "B", 12)
	pdf.Write(6, submission.RegistrationNumber)

	pdf.SetFont("Arial", "", 12)
	pdf.Write(6, fmt.Sprintf(" is an undergraduate B.Tech %dth year student in the Department of ", year))

	pdf.SetFont("Arial", "B", 12)
	pdf.Write(6, submission.Department)

	pdf.SetFont("Arial", "", 12)
	pdf.Write(6, ", Manipal University Jaipur.")

	pdf.Write(6, "He wishes to apply for an Internship/ Industrial Training in your esteemed organization")
	if submission.CompanyName != nil && *submission.CompanyName != "" {
		pdf.Write(6, ", ")
		pdf.SetFont("Arial", "B", 12)
		pdf.Write(6, *submission.CompanyName)
	}
	pdf.SetFont("Arial", "", 12)
	pdf.Write(6, ". The university has no objection to his undergoing an ")

	pdf.SetFont("Arial", "B", 12)
	pdf.Write(6, "Internship Training Program")

	pdf.SetFont("Arial", "", 12)
	pdf.Write(6, fmt.Sprintf(" from %s to %s.", startDateStr, endDateStr))

	pdf.SetFont("Arial", "", 12)
	pdf.Write(6, "\n\nWith Best Regards,\n\n\n")

	pdf.Ln(12)

	/////////////////////////////////////////////////////////////////////////////////////////////////////

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 6, "Prof (Dr) Neha Chaudhary", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 6, "HoD, Department of Computer Science & Engineering", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, "School of Computer Science & Engineering (SCSE)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Manipal University Jaipur, Rajasthan (INDIA)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Ph.: 0141-3999100 (Ext No 768)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Email: chaudhary.neha@jaipur.manipal.edu", "", 1, "L", false, 0, "")

	///////////////////////////////////////////////////////////////////////////////////////////////////

	uploadsDir := filepath.Join("../uploads", "NOC")
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
