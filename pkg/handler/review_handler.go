package handler

import (
	"MUJ_automated_mail_generation/pkg/database"
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

	// Log the inputs
	fmt.Printf("Creating review with submission_id: %d, reviewer_id: %d, status: %s, remarks: %s\n",
		input.SubmissionID, input.ReviewerID, input.Status, input.Comments)

	reviewID, err := database.CreateReview(input.SubmissionID, input.ReviewerID, input.Status, input.Comments)
	if err != nil {
		fmt.Printf("Error creating review: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"review_id": reviewID, "message": "Review created successfully"})
}

func UpdateReviewHandler(c *gin.Context) {
	var input struct {
		Status  string `json:"status" binding:"required"`
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

	err = database.UpdateReview(reviewID, input.Status, input.Remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
}

func GetReviewsBySubmissionHandler(c *gin.Context) {
	submissionID, err := strconv.Atoi(c.Param("submission_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	reviews, err := database.GetReviewsBySubmission(submissionID)
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

	reviews, err := database.GetReviewsByReviewer(reviewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func GetAllReviewsHandler(c *gin.Context) {
	reviews, err := database.GetAllReviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
