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

func CreateReviewHandler(c *gin.Context) {
	var input struct {
		SubmissionID int    `json:"submission_id" binding:"required"`
		ReviewerID   int    `json:"reviewer_id" binding:"required"`
		Status       string `json:"status" binding:"required"`
		Comments     string `json:"comments"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Printf("Creating reviewer review with submission_id: %d, reviewer_id: %d, status: %s, comments: %s\n",
		input.SubmissionID, input.ReviewerID, input.Status, input.Comments)

	reviewID, err := database.CreateReviewerReview(input.SubmissionID, input.ReviewerID, input.Status, input.Comments)
	if err != nil {
		fmt.Printf("Error creating reviewer review: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reviewer review"})
		return
	}

	if input.Status == "Approved" {
		err := database.UpdateSubmissionStatus(input.SubmissionID, "Approved")
		if err != nil {
			fmt.Printf("Error updating submission status: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission status"})
			return
		}
	}

	// If the submission is rejected or rework, send an email from the reviewer's email
	if input.Status == "Rejected" || input.Status == "Rework" {
		// Fetch the student submission from the database
		var submission model.StudentSubmission
		err := database.DB.Get(&submission, "SELECT * FROM student_submissions WHERE id = $1", input.SubmissionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student submission"})
			return
		}

		// Fetch the reviewer's details (including email) from the database
		var reviewer model.Reviewer
		err = database.DB.Get(&reviewer, "SELECT * FROM reviewers WHERE id = $1", input.ReviewerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviewer details"})
			return
		}

		// Create the email subject and body
		subject := "Your Placement Application Status"
		body := fmt.Sprintf("Dear %s,\n\nYour placement application has been %s.\n\nComments: %s\n\nBest regards",
			submission.Name, input.Status, input.Comments) // Add reviewer's name in the email

		// Send email using the reviewer's email ID as the sender
		err = util.SendEmail(reviewer.Email, submission.OfficialMailID, subject, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email to student"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"review_id": reviewID,
		"message":   "Reviewer review created successfully",
	})
}

func UpdateReviewHandler(c *gin.Context) {
	var input struct {
		Status   string `json:"status" binding:"required"`
		Comments string `json:"comments"`
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

	err = database.UpdateReviewerReview(reviewID, input.Status, input.Comments)
	if err != nil {
		fmt.Printf("Error updating reviewer review: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reviewer review status/comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reviewer review updated successfully",
	})
}

func GetReviewsBySubmissionHandler(c *gin.Context) {
	submissionID, err := strconv.Atoi(c.Param("submission_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	reviews, err := database.GetReviewerReviewsBySubmission(submissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func GetReviewsByReviewerHandler(c *gin.Context) {
	reviewerID, err := strconv.Atoi(c.Param("reviewer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reviewer ID"})
		return
	}

	reviews, err := database.GetReviewerReviewsByReviewer(reviewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func GetAllReviewerReviewsHandler(c *gin.Context) {
	reviews, err := database.GetAllReviewerReviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	// Return the reviews in a response
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
