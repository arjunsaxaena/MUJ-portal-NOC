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

	fmt.Printf("Creating review with submission_id: %d, reviewer_id: %d, status: %s, comments: %s\n",
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

	// Update the review status and comments
	err = database.UpdateReview(reviewID, input.Status, input.Comments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review status/comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
}

func UpdateHodReviewHandler(c *gin.Context) {
	var input struct {
		HodAction  string `json:"hod_action" binding:"required"`
		HodRemarks string `json:"hod_remarks"`
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

	err = database.UpdateHodAction(reviewID, input.HodAction, input.HodRemarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update HoD action/remarks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated with HoD action successfully"})
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
		fmt.Printf("Error in GetAllReviews: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	// Convert sql.NullString to regular string for JSON response
	for i := range reviews {
		if !reviews[i].HodAction.Valid {
			reviews[i].HodAction.String = "" // Set to empty string if NULL
		}
		if !reviews[i].HodRemarks.Valid {
			reviews[i].HodRemarks.String = ""
		}
	}

	c.JSON(http.StatusOK, reviews)
}
