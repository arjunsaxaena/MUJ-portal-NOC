package handler

import (
	"MUJ_automated_mail_generation/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateHodReviewHandler handles the request to create a new HoD review.
func CreateHodReviewHandler(c *gin.Context) {
	var input struct {
		SubmissionID int    `json:"submission_id" binding:"required"`
		HodID        int    `json:"hod_id" binding:"required"`
		Action       string `json:"action" binding:"required"`
		Remarks      string `json:"remarks"`
	}

	// Bind the request body to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the database function to create the review
	reviewID, err := database.CreateHodReview(input.SubmissionID, input.HodID, input.Action, input.Remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HoD review"})
		return
	}

	// Return a success response with the review ID
	c.JSON(http.StatusCreated, gin.H{"review_id": reviewID, "message": "HoD review created successfully"})
}

// UpdateHodReviewHandler handles the request to update an existing HoD review.
func UpdateHodReviewHandler(c *gin.Context) {
	var input struct {
		Action  string `json:"action" binding:"required"`
		Remarks string `json:"remarks"`
	}

	// Get the review ID from the URL parameter
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// Bind the request body to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the database function to update the review
	err = database.UpdateHodReview(reviewID, input.Action, input.Remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update HoD review"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "HoD review updated successfully"})
}

// GetHodReviewsBySubmissionHandler handles the request to get all HoD reviews for a submission.
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

	// Return the list of reviews
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// GetHodReviewsByHodHandler handles the request to get all HoD reviews by a specific HoD.
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

	// Return the list of reviews
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// GetAllHodReviewsHandler handles the request to get all HoD reviews.
func GetAllHodReviewsHandler(c *gin.Context) {
	reviews, err := database.GetAllHodReviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all HoD reviews"})
		return
	}

	// Return the list of reviews
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
