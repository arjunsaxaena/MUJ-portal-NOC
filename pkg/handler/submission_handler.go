package handler

import (
	"MUJ_automated_mail_generation/pkg/database"
	"MUJ_automated_mail_generation/pkg/model"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SubmitHandler(c *gin.Context) {
	var submission model.StudentSubmission

	// for debug: print the raw form data
	fmt.Println("Form data received:")

	if err := c.ShouldBind(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Struct before DB insert: %+v\n", submission)

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

	submission.Status = "Pending" // Set additional fields

	err = database.CreateSubmission(&submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission received successfully"})
}

func GetAllSubmissionsHandler(c *gin.Context) {
	submissions, err := database.GetAllSubmissions()
	if err != nil {
		log.Printf("Error fetching submissions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"submissions": submissions})
}
