package handler

import (
	"MUJ_automated_mail_generation/pkg/database"
	"MUJ_automated_mail_generation/pkg/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SubmitHandler(c *gin.Context) {
	var submission model.StudentSubmission
	if err := c.ShouldBind(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the struct before database insertion (for debugging)
	fmt.Printf("Struct before DB insert: %+v\n", submission)

	// Parse internship_start_date (which is a string in form data) into time.Time
	startDate, err := time.Parse("2006-01-02", submission.InternshipStartDate.Format("2006-01-02"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid internship_start_date format, use YYYY-MM-DD"})
		return
	}

	// Parse internship_end_date (which is a string in form data) into time.Time
	endDate, err := time.Parse("2006-01-02", submission.InternshipEndDate.Format("2006-01-02"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid internship_end_date format, use YYYY-MM-DD"})
		return
	}

	// Assign the parsed time values back to the struct (which expects time.Time)
	submission.InternshipStartDate = startDate
	submission.InternshipEndDate = endDate

	// Convert PackagePPO and StipendAmount to float64
	packagePPO, err := strconv.ParseFloat(submission.PackagePPO, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package_ppo value"})
		return
	}
	stipendAmount, err := strconv.ParseFloat(submission.StipendAmount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stipend_amount value"})
		return
	}

	// Assign parsed values back to submission
	submission.PackagePPO = fmt.Sprintf("%.2f", packagePPO)
	submission.StipendAmount = fmt.Sprintf("%.2f", stipendAmount)

	// Handle offer letter upload
	offerLetter, _ := c.FormFile("offer_letter")
	if offerLetter != nil {
		offerLetterPath := fmt.Sprintf("./uploads/offer_letters/%s_%s", submission.RegistrationNumber, offerLetter.Filename)
		err := c.SaveUploadedFile(offerLetter, offerLetterPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save offer letter"})
			return
		}
		submission.OfferLetterPath = offerLetterPath
	}

	// Handle mail copy upload
	mailCopy, _ := c.FormFile("mail_copy")
	if mailCopy != nil {
		mailCopyPath := fmt.Sprintf("./uploads/mail_copies/%s_%s", submission.RegistrationNumber, mailCopy.Filename)
		err := c.SaveUploadedFile(mailCopy, mailCopyPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save mail copy"})
			return
		}
		submission.MailCopyPath = mailCopyPath
	}

	// Set additional fields
	submission.Status = "Pending"

	// Insert into database
	_, err = database.DB.NamedExec(`
		INSERT INTO student_submissions (
			registration_number, name, official_mail_id, mobile_number, department, section, offer_type, 
			company_name, company_address, offer_type_detail, package_ppo, stipend_amount, 
			internship_start_date, internship_end_date, offer_letter_path, terms_accepted, 
			status, created_at
		) 
		VALUES (
			:registration_number, :name, :official_mail_id, :mobile_number, :department, :section, :offer_type, 
			:company_name, :company_address, :offer_type_detail, :package_ppo, :stipend_amount, 
			:internship_start_date, :internship_end_date, :offer_letter_path, :terms_accepted, 
			:status, CURRENT_TIMESTAMP
		)
	`, &submission)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission received successfully"})
}
