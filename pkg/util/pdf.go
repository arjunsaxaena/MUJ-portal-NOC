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

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 6, "Sub: Recommendation for "+title+" "+submission.Name+" carrying out internship cum project in your esteemed Organization", "", "J", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 6, "Dear Sir/Madam,", "", "J", false)
	pdf.Ln(5)

	fullDept := getFullDepartmentName(submission.Department)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 6, "This is to certify that "+title+" ", "", "J", false)
	pdf.SetFont("Arial", "B", 12)
	pdf.MultiCell(0, 6, submission.Name+" ", "", "J", false)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 6, "(Reg No. "+submission.RegistrationNumber+") is a student of Manipal University Jaipur, India, studying in the "+submission.Semester+" semester of the four-year B.Tech Degree Program in the Department of "+fullDept+", Section "+submission.Section+".", "", "J", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 6, "This recommendation is issued with reference to the application for an internship/project in your esteemed organization for a duration from ", "", "J", false)
	pdf.SetFont("Arial", "B", 12)
	pdf.MultiCell(0, 6, startDateFormatted+" ", "", "J", false)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 6, "to ", "", "J", false)
	pdf.SetFont("Arial", "B", 12)
	pdf.MultiCell(0, 6, endDateFormatted+".", "", "J", false)
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 6, "This Internship/Project would add value to the academic career of the student. So I request you to kindly allow our student to undergo Internship/Project at your organization.", "", "J", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 6, "Manipal University Jaipur has no objection for "+title+" "+submission.Name+" in doing an internship at your organization and has been advised to abide by both MUJ's and the interning organization's ethics/rules/regulations/values and work culture without compromising on integrity and self-discipline.", "", "J", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 6, "Thanking you.", "", "J", false)
	pdf.MultiCell(0, 6, "Yours sincerely,", "", "J", false)
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
	pdf.CellFormat(95, 10, nocText, "", 0, "L", false, 0, "")

	currentDate := time.Now().Format("02-Jan-2006")
	pdf.CellFormat(0, 10, currentDate, "", 0, "R", false, 0, "")
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

	startDateStr := startDate.Format("02 Jan 2006")
	endDateStr := endDate.Format("02 Jan 2006")

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

	title := "Mr."
	if submission.Gender == "Female" {
		title = "Ms."
	}

	fullDept := getFullDepartmentName(submission.Department)

	content := title + " " + submission.Name + ", Reg No.- " + submission.RegistrationNumber +
		" is an undergraduate B.Tech " + strconv.Itoa(year) + "th year student in the Department of " +
		fullDept + ", Manipal University Jaipur. He wishes to apply for an Internship/ Industrial Training in your esteemed organization"

	if submission.CompanyName != nil && *submission.CompanyName != "" {
		content += ", " + *submission.CompanyName
	}

	content += ". The university has no objection to his undergoing an Internship Training Program from " +
		startDateStr + " to " + endDateStr + "."

	pdf.MultiCell(0, 6, content, "", "J", false)
	pdf.Ln(10)

	pdf.MultiCell(0, 6, "With Best Regards,", "", "J", false)
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.MultiCell(0, 6, "Prof (Dr) Neha Chaudhary", "", "L", false)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 6, "HoD, Department of Computer Science & Engineering", "", "L", false)
	pdf.MultiCell(0, 6, "School of Computer Science & Engineering (SCSE)", "", "L", false)
	pdf.MultiCell(0, 6, "Manipal University Jaipur, Rajasthan (INDIA)", "", "L", false)
	pdf.MultiCell(0, 6, "Ph.: 0141-3999100 (Ext No 768)", "", "L", false)
	pdf.MultiCell(0, 6, "Email: chaudhary.neha@jaipur.manipal.edu", "", "L", false)

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
