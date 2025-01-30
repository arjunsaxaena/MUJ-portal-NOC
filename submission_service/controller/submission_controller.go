package controller

import (
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/submission_service/repository"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

const uploadPath = "../../uploads"

func SubmitHandler(c *gin.Context) {
	var submission model.StudentSubmission

	form, _ := c.MultipartForm()
	fmt.Println("Raw Form Data Received:", form.Value)
	fmt.Println("Files Received:", form.File)

	if err := c.ShouldBind(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("Error binding form data: %v\n", err)
		return
	}

	studentFilters := model.GetStudentFilters{
		RegistrationNumber: submission.RegistrationNumber,
		EmailID:            submission.OfficialMailID,
	}

	students, err := repository.GetStudents(studentFilters)
	if err != nil || len(students) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found with given Registration Number and Email ID"})
		fmt.Printf("Student validation failed: %v\n", err)
		return
	}

	submission.PackagePPO = formatFloat(submission.PackagePPO)
	submission.StipendAmount = formatFloat(submission.StipendAmount)

	if (submission.HRDEmail == nil || *submission.HRDEmail == "") && (submission.HRDNumber == nil || *submission.HRDNumber == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either HRDEmail or HRDNumber must be provided"})
		return
	}

	if offerLetter, err := c.FormFile("offerLetter"); err == nil {
		offerPath := saveFileLocally(offerLetter, "offerLetters", submission.RegistrationNumber)
		if offerPath != "" {
			submission.OfferLetterPath = offerPath
			fmt.Printf("Offer letter saved at: %s\n", offerPath)
		}
	} else {
		fmt.Println("No offer letter file received")
	}

	if mailCopy, err := c.FormFile("mailCopy"); err == nil {
		mailPath := saveFileLocally(mailCopy, "mailCopies", submission.RegistrationNumber)
		if mailPath != "" {
			submission.MailCopyPath = mailPath
			fmt.Printf("Mail copy saved at: %s\n", mailPath)
		}
	} else {
		fmt.Println("No mail copy file received")
	}

	if submission.OfferLetterPath == "" && submission.HasOfferLetter {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offer letter file path is missing"})
		return
	}
	if submission.MailCopyPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mail copy file path is missing"})
		return
	}

	submission.Status = "Pending"
	fmt.Printf("Submission ready for DB insert: %+v\n", submission)

	if err := repository.CreateSubmission(&submission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
		fmt.Printf("Error saving submission to DB: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission received successfully"})
	fmt.Println("Submission received successfully")
}

func formatFloat(value string) string {
	if value == "" {
		return "0.00"
	}
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return fmt.Sprintf("%.2f", floatVal)
	}
	return "0.00"
}

func saveFileLocally(fileHeader *multipart.FileHeader, folder, regNumber string) string {
	dstDir := filepath.Join(uploadPath, folder)
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return ""
	}

	dstPath := filepath.Join(dstDir, fmt.Sprintf("%s_%s", regNumber, fileHeader.Filename))
	dstFile, err := os.Create(dstPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return ""
	}
	defer dstFile.Close()

	srcFile, err := fileHeader.Open()
	if err != nil {
		fmt.Printf("Error opening uploaded file: %v\n", err)
		return ""
	}
	defer srcFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return ""
	}

	return "/files/" + folder + "/" + filepath.Base(dstPath)
}
