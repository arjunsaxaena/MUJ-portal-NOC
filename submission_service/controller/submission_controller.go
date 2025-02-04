package controller

import (
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/pkg/util"
	"MUJ_AMG/submission_service/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	submission.PackagePPO = util.FormatFloat(submission.PackagePPO)
	submission.StipendAmount = util.FormatFloat(submission.StipendAmount)

	if (submission.HRDEmail == nil || *submission.HRDEmail == "") && (submission.HRDNumber == nil || *submission.HRDNumber == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either HRDEmail or HRDNumber must be provided"})
		return
	}

	offerLetterUploaded := false
	if offerLetter, err := c.FormFile("offerLetter"); err == nil {
		offerPath := util.SaveFile(offerLetter, "offerLetters", submission.RegistrationNumber)
		if offerPath != "" {
			submission.OfferLetterPath = offerPath
			offerLetterUploaded = true
			fmt.Printf("Offer letter saved at: %s\n", offerPath)
		}
	} else {
		fmt.Println("No offer letter file received")
	}

	if mailCopy, err := c.FormFile("mailCopy"); err == nil {
		mailPath := util.SaveFile(mailCopy, "mailCopies", submission.RegistrationNumber)
		if mailPath != "" {
			submission.MailCopyPath = mailPath
			fmt.Printf("Mail copy saved at: %s\n", mailPath)
		}
	} else {
		fmt.Println("No mail copy file received")
	}

	if submission.NocType == "Specific" && !offerLetterUploaded {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offer letter is required for 'specific' NOC type"})
		return
	}

	if submission.NocType == "Generic" {
		fmt.Println("NOC type is 'Generic', offer letter is not required")
	} else if submission.NocType == "Specific" && !offerLetterUploaded {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offer letter is required for 'Specific' NOC type"})
		return
	} else if submission.NocType != "Generic" && submission.NocType != "Specific" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid NOC type. Must be 'Generic' or 'Specific'"})
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

///////////////////////////////////////////////////////////////////////////////////////////////

type OTPRequest struct {
	Email string `json:"email"`
}

func GenerateOTPHandler(c *gin.Context) {
	var request OTPRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	email := request.Email
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	generatedOTP := util.GenerateOTP(email)
	emailBody := fmt.Sprintf("Your OTP for email verification is: %s. It is valid for 5 minutes.", generatedOTP)

	err := util.SendEmail("arjunsaxena04@gmail.com", email, "Email Verification OTP", emailBody, "bjkwwhugjefvcdoa")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email"})
		fmt.Printf("Error sending OTP email: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func ValidateOTPHandler(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	email := request.Email
	otp := request.OTP

	if email == "" || otp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and OTP are required"})
		return
	}

	if !util.ValidateOTP(email, otp) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	util.MarkEmailValidated(email)

	c.JSON(http.StatusOK, gin.H{"message": "OTP validated successfully"})
}
