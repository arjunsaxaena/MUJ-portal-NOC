package controller

import (
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/pkg/util"
	"MUJ_AMG/submission_service/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SubmitHandler(c *gin.Context) {
	var submission model.StudentSubmission

	// for debug: print the raw form data
	form, _ := c.MultipartForm()
	fmt.Println("Raw Form Data Received:", form.Value)
	fmt.Println("Files Received:", form.File)

	// Log the form data for debugging purposes
	if err := c.ShouldBind(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("Error binding form data: %v\n", err)
		return
	}

	// logging for debug
	fmt.Printf("Struct before DB insert: %+v\n", submission)

	if submission.PackagePPO != "" {
		packagePPO, err := strconv.ParseFloat(submission.PackagePPO, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package_ppo value"})
			fmt.Printf("Error parsing PackagePPO: %v\n", err)
			return
		}
		submission.PackagePPO = fmt.Sprintf("%.2f", packagePPO)
		fmt.Printf("Formatted PackagePPO: %s\n", submission.PackagePPO)
	} else {
		submission.PackagePPO = ""
	}

	if submission.StipendAmount != "" {
		stipendAmount, err := strconv.ParseFloat(submission.StipendAmount, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stipend_amount value"})
			fmt.Printf("Error parsing StipendAmount: %v\n", err)
			return
		}
		submission.StipendAmount = fmt.Sprintf("%.2f", stipendAmount)
		fmt.Printf("Formatted StipendAmount: %s\n", submission.StipendAmount)
	} else {
		submission.StipendAmount = "0.00"
	}

	offerLetter, _ := c.FormFile("offerLetter")
	if offerLetter != nil {
		fmt.Printf("Received offer letter file: %s\n", offerLetter.Filename)

		file, err := offerLetter.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open offer letter file"})
			fmt.Printf("Error opening offer letter file: %v\n", err)
			return
		}
		defer file.Close()

		offerLetterKey := fmt.Sprintf("offerLetters/%s_%s", submission.RegistrationNumber, offerLetter.Filename)
		err = util.UploadFileToS3("muj-student-data", file, offerLetterKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload offer letter to S3"})
			fmt.Printf("Error uploading offer letter to S3: %v\n", err)
			return
		}
		offerLetterURL := fmt.Sprintf("https://muj-student-data.s3.amazonaws.com/%s", offerLetterKey)
		submission.OfferLetterPath = offerLetterURL
		fmt.Printf("Offer letter uploaded successfully. URL: %s\n", offerLetterURL)
	} else {
		submission.HasOfferLetter = false
		fmt.Println("No offer letter file received")
	}

	mailCopy, _ := c.FormFile("mailCopy")
	if mailCopy != nil {
		fmt.Printf("Received mail copy file: %s\n", mailCopy.Filename)

		file, err := mailCopy.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open mail copy file"})
			fmt.Printf("Error opening mail copy file: %v\n", err)
			return
		}
		defer file.Close()

		mailCopyKey := fmt.Sprintf("mailCopies/%s_%s", submission.RegistrationNumber, mailCopy.Filename)
		err = util.UploadFileToS3("muj-student-data", file, mailCopyKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload mail copy to S3"})
			fmt.Printf("Error uploading mail copy to S3: %v\n", err)
			return
		}
		mailCopyURL := fmt.Sprintf("https://muj-student-data.s3.amazonaws.com/%s", mailCopyKey)
		submission.MailCopyPath = mailCopyURL
		fmt.Printf("Mail copy uploaded successfully. URL: %s\n", mailCopyURL)
	} else {
		fmt.Println("No mail copy file received")
	}

	if submission.OfferLetterPath == "" && submission.HasOfferLetter {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offer letter file path is missing"})
		fmt.Println("Offer letter file path is missing")
		return
	}
	if submission.MailCopyPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mail copy file path is missing"})
		fmt.Println("Mail copy file path is missing")
		return
	}

	submission.Status = "Pending"
	fmt.Printf("Submission ready for DB insert: %+v\n", submission)

	err := repository.CreateSubmission(&submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
		fmt.Printf("Error saving submission to DB: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission received successfully"})
	fmt.Println("Submission received successfully")
}
