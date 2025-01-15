package controller

import (
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/pkg/util"
	"MUJ_AMG/portal_service/config"
	"MUJ_AMG/portal_service/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

	reviewID, err := repository.CreateReview(input.SubmissionID, input.ReviewerID, input.Status, input.Comments)
	if err != nil {
		fmt.Printf("Error creating reviewer review: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reviewer review"})
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
		return
	}

	submissionURL := fmt.Sprintf("%s/submissions?id=%d", cfg.SubmissionServiceURL, input.SubmissionID)
	submissionResp, err := http.Get(submissionURL)
	if err != nil || submissionResp.StatusCode != http.StatusOK {
		log.Printf("Error fetching submission via API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submission"})
		return
	}
	defer submissionResp.Body.Close()

	reviewerFilters := model.GetReviewerFilters{
		ID: strconv.Itoa(input.ReviewerID),
	}
	reviewers, err := repository.GetReviewers(reviewerFilters)
	if err != nil {
		log.Printf("Error fetching reviewer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviewer details"})
		return
	}

	if len(reviewers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reviewer not found"})
		return
	}

	reviewer := reviewers[0]

	if input.Status == "Approved" {
		updateSubmissionURL := fmt.Sprintf("%s/submissions", cfg.SubmissionServiceURL)
		updateSubmissionBody := struct {
			Status string `json:"status"`
		}{
			Status: "Approved",
		}
		updateSubmissionJSON, err := json.Marshal(updateSubmissionBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal update submission body"})
			return
		}

		updateSubmissionReq, err := http.NewRequest("PUT", updateSubmissionURL, bytes.NewBuffer(updateSubmissionJSON))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create update submission request"})
			return
		}
		updateSubmissionReq.Header.Set("Content-Type", "application/json")

		q := updateSubmissionReq.URL.Query()
		q.Add("id", strconv.Itoa(input.SubmissionID))
		updateSubmissionReq.URL.RawQuery = q.Encode()

		updateSubmissionResp, err := http.DefaultClient.Do(updateSubmissionReq)
		if err != nil || updateSubmissionResp.StatusCode != http.StatusOK {
			log.Printf("Error updating submission status via API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission status"})
			return
		}
		defer updateSubmissionResp.Body.Close()
	}

	if input.Status == "Rejected" || input.Status == "Rework" {
		updateSubmissionURL := fmt.Sprintf("%s/submissions", cfg.SubmissionServiceURL)
		updateSubmissionBody := struct {
			Status string `json:"status"`
		}{
			Status: input.Status,
		}
		updateSubmissionJSON, err := json.Marshal(updateSubmissionBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal update submission body"})
			return
		}

		updateSubmissionReq, err := http.NewRequest("PUT", updateSubmissionURL, bytes.NewBuffer(updateSubmissionJSON))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create update submission request"})
			return
		}
		updateSubmissionReq.Header.Set("Content-Type", "application/json")

		q := updateSubmissionReq.URL.Query()
		q.Add("id", strconv.Itoa(input.SubmissionID))
		updateSubmissionReq.URL.RawQuery = q.Encode()

		updateSubmissionResp, err := http.DefaultClient.Do(updateSubmissionReq)
		if err != nil || updateSubmissionResp.StatusCode != http.StatusOK {
			log.Printf("Error updating submission status via API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission status"})
			return
		}
		defer updateSubmissionResp.Body.Close()

		submissionURL := fmt.Sprintf("%s/submissions?id=%d", cfg.SubmissionServiceURL, input.SubmissionID)
		submissionResp, err := http.Get(submissionURL)
		if err != nil || submissionResp.StatusCode != http.StatusOK {
			log.Printf("Error fetching submission via API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submission"})
			return
		}
		defer submissionResp.Body.Close()

		var submissionResponse struct {
			Submissions []model.StudentSubmission `json:"submissions"`
		}
		err = json.NewDecoder(submissionResp.Body).Decode(&submissionResponse)
		if err != nil {
			log.Printf("Error decoding submission data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode submission data"})
			return
		}

		if len(submissionResponse.Submissions) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
			return
		}

		submission := submissionResponse.Submissions[0]

		subject := "Your NOC Application Status"
		body := fmt.Sprintf("Dear %s,\n\nYour NOC application has been %s.\n\nComments: %s\n\nBest regards",
			submission.Name, input.Status, input.Comments)

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

	reviewID, err := strconv.Atoi(c.DefaultQuery("id", ""))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	reviewFilters := model.GetReviewFilters{
		ID: strconv.Itoa(reviewID),
	}
	reviews, err := repository.GetReviews(reviewFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch review"})
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	review := reviews[0]

	err = repository.UpdateReview(reviewID, input.Status, input.Comments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reviewer review status/comments"})
		return
	}

	if review.Status == "Approved" && (input.Status == "Rejected" || input.Status == "Rework") {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Printf("Error loading config: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
			return
		}

		reviewerFilters := model.GetReviewerFilters{
			ID: strconv.Itoa(review.ReviewerID),
		}
		reviewers, err := repository.GetReviewers(reviewerFilters)
		if err != nil {
			log.Printf("Error fetching reviewer details: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviewer details"})
			return
		}

		if len(reviewers) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reviewer not found"})
			return
		}

		var reviewer = reviewers[0]

		submissionURL := fmt.Sprintf("%s/submissions?id=%d", cfg.SubmissionServiceURL, review.SubmissionID)
		submissionResp, err := http.Get(submissionURL)
		if err != nil || submissionResp.StatusCode != http.StatusOK {
			log.Printf("Error fetching submission via API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submission"})
			return
		}
		defer submissionResp.Body.Close()

		var submissionResponse struct {
			Submissions []model.StudentSubmission `json:"submissions"`
		}
		err = json.NewDecoder(submissionResp.Body).Decode(&submissionResponse)
		if err != nil {
			log.Printf("Error decoding submission data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode submission data"})
			return
		}

		if len(submissionResponse.Submissions) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
			return
		}

		submission := submissionResponse.Submissions[0]

		subject := "Your NOC Application Status"
		body := fmt.Sprintf("Dear %s,\n\nYour NOC application has been %s.\n\nComments: %s\n\nBest regards",
			submission.Name, input.Status, input.Comments)

		err = util.SendEmail(reviewer.Email, submission.OfficialMailID, subject, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email to student"})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reviewer review updated successfully",
	})
}

func GetReviewsHandler(c *gin.Context) {
	submissionID := c.DefaultQuery("submission_id", "")
	reviewerID := c.DefaultQuery("reviewer_id", "")
	status := c.DefaultQuery("status", "")
	reviewID := c.DefaultQuery("review_id", "")

	filters := model.GetReviewFilters{
		ID:           reviewID,
		SubmissionID: submissionID,
		ReviewerID:   reviewerID,
		Status:       status,
	}

	reviews, err := repository.GetReviews(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
