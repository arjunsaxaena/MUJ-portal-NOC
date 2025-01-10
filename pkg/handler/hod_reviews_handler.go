package handler

import (
	"MUJ_automated_mail_generation/pkg/database"
	"MUJ_automated_mail_generation/pkg/model"
	"MUJ_automated_mail_generation/pkg/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateHodReviewHandler(c *gin.Context) {
	var input struct {
		SubmissionID int    `json:"submission_id" binding:"required"`
		HodID        int    `json:"hod_id" binding:"required"`
		Action       string `json:"action" binding:"required"`
		Remarks      string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the student submission from the database
	var submission model.StudentSubmission
	err := database.DB.Get(&submission, "SELECT * FROM student_submissions WHERE id = $1", input.SubmissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student submission"})
		return
	}

	// Fetch the HoD's email from the database
	var hod model.Reviewer
	err = database.DB.Get(&hod, "SELECT * FROM reviewers WHERE id = $1", input.HodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch HoD details"})
		return
	}

	// Handle review actions: Rejected or Rework
	if input.Action == "Rejected" || input.Action == "Rework" {
		// Create email subject and body
		subject := "Your Placement Application Status - HoD Review"
		body := fmt.Sprintf("Dear %s,\n\nYour placement application has been %s.\n\nHoD Comments: %s\n\nBest regards",
			submission.Name, input.Action, input.Remarks)

		// Send email from HoD's email ID to student's email
		err = util.SendEmail(hod.Email, submission.OfficialMailID, subject, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email to student"})
			return
		}
	}

	reviewID, err := database.CreateHodReview(input.SubmissionID, input.HodID, input.Action, input.Remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HoD review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"review_id": reviewID, "message": "HoD review created successfully"})
}

func UpdateHodReviewHandler(c *gin.Context) {
	var input struct {
		Action  string `json:"action" binding:"required"`
		Remarks string `json:"remarks"`
	}

	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = database.UpdateHodReview(reviewID, input.Action, input.Remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update HoD review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "HoD review updated successfully"})
}

func GetHodReviewsBySubmissionHandler(c *gin.Context) {
	submissionID, err := strconv.Atoi(c.Param("submission_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	reviews, err := database.GetHodReviewsBySubmission(submissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch HoD reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func GetHodReviewsByHodHandler(c *gin.Context) {
	hodID, err := strconv.Atoi(c.Param("hod_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid HoD ID"})
		return
	}

	reviews, err := database.GetHodReviewsByHod(hodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch HoD reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func GetAllHodReviewsHandler(c *gin.Context) {
	reviews, err := database.GetAllHodReviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all HoD reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
