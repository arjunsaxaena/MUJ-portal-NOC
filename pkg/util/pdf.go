package util

import (
	"MUJ_AMG/pkg/model"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func getFullDepartmentName(dept string) string {
	switch dept {
	case "CSE":
		return "Computer Science and Engineering"
	case "IT":
		return "Information Technology"
	default:
		return dept
	}
}

func getRoman(semester string) string {
	switch semester {
	case "5":
		return "V"
	case "7":
		return "VII"
	default:
		return semester
	}
}

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
	nocText := fmt.Sprintf("MUJ/FoSTA/D%s/2025/%s%s/%s", submission.Department, submission.Semester, submission.Section, submission.RegistrationNumber[len(submission.RegistrationNumber)-4:])
	pdf.CellFormat(95, 10, nocText, "", 0, "L", false, 0, "") // left-aligned text

	currentDate := time.Now().Format("02-Jan-2006")
	pdf.CellFormat(0, 10, currentDate, "", 0, "R", false, 0, "") // right-aligned date on the same line
	pdf.Ln(20)

	pdf.SetFont("Arial", "BU", 16)
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

	title := "Mr."
	if submission.Gender == "Female" {
		title = "Ms."
	}

	fullDept := getFullDepartmentName(submission.Department)
	romanSemester := getRoman(submission.Semester)

	var subjectText string
	subjectText = fmt.Sprintf("Sub: Recommendation for %s %s carrying out internship cum project in your esteemed Organization.",
		title, submission.Name)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 6, subjectText, "", "J", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 6, "Dear Sir/Madam,", "", "J", false)
	pdf.Ln(5)

	bodyText := fmt.Sprintf("This is to certify that %s %s (Reg No. %s) is a student of Manipal University Jaipur, India, studying in the %s Semester of the four-year B.Tech Degree Programme in the Department of %s.", title, submission.Name, submission.RegistrationNumber, romanSemester, fullDept)

	pdf.MultiCell(0, 6, bodyText, "", "J", false)
	pdf.Ln(5)

	var internshipText string
	if submission.NocType != "generic" && submission.CompanyName != nil {
		internshipText = fmt.Sprintf("This recommendation is issued with reference to the application for an internship/project in your esteemed organization, %s for a duration from %s to %s.", *submission.CompanyName, startDateFormatted, endDateFormatted)
	} else {
		internshipText = fmt.Sprintf("This recommendation is issued with reference to the application for an internship/project in your esteemed organization for a duration from %s to %s.", startDateFormatted, endDateFormatted)
	}

	pdf.MultiCell(0, 6, internshipText, "", "J", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 6, "This Internship/Project would add value to the academic career of the student. So I request you to kindly allow our student to undergo Internship/Project at your organization.", "", "J", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 6, "Manipal University Jaipur has no objection for "+title+" "+submission.Name+" in doing an internship at your organization and has been advised to abide by both MUJ's and the interning organization's ethics/rules/regulations/values and work culture without compromising on integrity and self-discipline.", "", "J", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 6, "Thanking you.", "", "J", false)
	pdf.MultiCell(0, 6, "Yours sincerely,", "", "J", false)
	pdf.Ln(22)

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
