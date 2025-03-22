package controller

import (
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/pkg/util"
	"MUJ_AMG/portal_service/repository"
	submissionRepository "MUJ_AMG/submission_service/repository"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateFpcReviewHandler(c *gin.Context) {
	var input struct {
		SubmissionID int    `json:"submission_id" binding:"required"`
		FpcID        int    `json:"fpc_id" binding:"required"`
		Status       string `json:"status" binding:"required"`
		Comments     string `json:"comments"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fpcFilters := model.GetFpCFilters{
		ID: strconv.Itoa(input.FpcID),
	}

	fpcs, err := repository.GetFpCs(fpcFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch FPC details"})
		return
	}

	if len(fpcs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "FPC not found"})
		return
	}

	fpc := fpcs[0]

	reviewID, err := repository.CreateFpcReview(input.SubmissionID, input.FpcID, input.Status, input.Comments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create FPC review"})
		return
	}

	if input.Status == "Approved" {
		err := submissionRepository.UpdateSubmission(input.SubmissionID, "Approved", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission status"})
			return
		}
	}

	if input.Status == "Rejected" {
		err := submissionRepository.UpdateSubmission(input.SubmissionID, "Rejected", "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission status"})
			return
		}

		submissions, err := submissionRepository.GetSubmissions(model.GetSubmissionFilters{
			ID: strconv.Itoa(input.SubmissionID),
		})
		if err != nil || len(submissions) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
			return
		}

		submission := submissions[0]

		subject := "Your NOC Application Status"
		body := fmt.Sprintf("Dear %s,\n\nYour NOC application has been %s.\n\nComments: %s\n\nBest regards",
			submission.Name, input.Status, input.Comments)

		err = util.SendEmail(fpc.Email, submission.OfficialMailID, subject, body, fpc.AppPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email to student"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"review_id": reviewID,
		"message":   "FPC review created successfully",
	})
}

func UpdateFpcReviewHandler(c *gin.Context) {
	var input struct {
		Status   string `json:"status" binding:"required"`
		Comments string `json:"comments"`
	}

	reviewID, err := strconv.Atoi(c.DefaultQuery("id", ""))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fpcReviewFilters := model.GetFpcReviewFilters{
		ID: strconv.Itoa(reviewID),
	}
	fpcReviews, err := repository.GetFpcReviews(fpcReviewFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fpc review"})
		return
	}

	if len(fpcReviews) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "FPC review not found"})
		return
	}

	fpcReview := fpcReviews[0]

	err = repository.UpdateFpcReview(reviewID, input.Status, input.Comments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fpc review status/comments"})
		return
	}

	if fpcReview.Status == "Approved" && input.Status == "Rejected" {
		fpcFilters := model.GetFpCFilters{
			ID: strconv.Itoa(fpcReview.FpcID),
		}
		fpcs, err := repository.GetFpCs(fpcFilters)
		if err != nil {
			log.Printf("Error fetching fpc details: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fpc details"})
			return
		}

		if len(fpcs) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "FPC not found"})
			return
		}

		fpc := fpcs[0]

		submissionFilters := model.GetSubmissionFilters{
			ID: strconv.Itoa(fpcReview.SubmissionID),
		}
		submissions, err := submissionRepository.GetSubmissions(submissionFilters)
		if err != nil {
			log.Printf("Error fetching submission with ID %d: %v", fpcReview.SubmissionID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submission"})
			return
		}

		if len(submissions) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
			return
		}

		submission := submissions[0]

		subject := "Your NOC Application Status"
		body := fmt.Sprintf("Dear %s,\n\nYour NOC application has been %s.\n\nComments: %s\n\nBest regards",
			submission.Name, input.Status, input.Comments)

		err = util.SendEmail(fpc.Email, submission.OfficialMailID, subject, body, fpc.AppPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "FPC review updated successfully",
	})
}

func GetFpcReviewsHandler(c *gin.Context) {
	submissionID := c.DefaultQuery("submission_id", "")
	fpcID := c.DefaultQuery("fpc_id", "")
	status := c.DefaultQuery("status", "")
	reviewID := c.DefaultQuery("review_id", "")

	filters := model.GetFpcReviewFilters{}
	if reviewID != "" {
		filters.ID = reviewID
	}
	if submissionID != "" {
		filters.SubmissionID = submissionID
	}
	if fpcID != "" {
		filters.FpcID = fpcID
	}
	if status != "" {
		filters.Status = status
	}

	reviews, err := repository.GetFpcReviews(filters)
	if err != nil {
		log.Printf("Error fetching reviews: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
