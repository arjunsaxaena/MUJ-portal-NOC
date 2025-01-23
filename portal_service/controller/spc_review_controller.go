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

func CreateSpcReviewHandler(c *gin.Context) {
	var input struct {
		SubmissionID int    `json:"submission_id" binding:"required"`
		SpcID        int    `json:"spc_id" binding:"required"`
		Status       string `json:"status" binding:"required"`
		Comments     string `json:"comments"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Printf("Creating spc review with submission_id: %d, spc_id: %d, status: %s, comments: %s\n",
		input.SubmissionID, input.SpcID, input.Status, input.Comments)

	reviewID, err := repository.CreateSpcReview(input.SubmissionID, input.SpcID, input.Status, input.Comments)
	if err != nil {
		fmt.Printf("Error creating spc review: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create spc review"})
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

	spcFilters := model.GetSpCFilters{
		ID: strconv.Itoa(input.SpcID),
	}
	spcs, err := repository.GetSpCs(spcFilters)
	if err != nil {
		log.Printf("Error fetching spc: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spc details"})
		return
	}

	if len(spcs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Spc not found"})
		return
	}

	spc := spcs[0]

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

		err = util.SendEmail(spc.Email, submission.OfficialMailID, subject, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email to student"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"review_id": reviewID,
		"message":   "Spc review created successfully",
	})
}

func UpdateSpcReviewHandler(c *gin.Context) {
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

	spcReviewFilters := model.GetSpcReviewFilters{
		ID: strconv.Itoa(reviewID),
	}
	spcReviews, err := repository.GetSpcReviews(spcReviewFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spc review"})
		return
	}

	if len(spcReviews) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Spc review not found"})
		return
	}

	spcReview := spcReviews[0]

	err = repository.UpdateSpcReview(reviewID, input.Status, input.Comments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update spc review status/comments"})
		return
	}

	if spcReview.Status == "Approved" && (input.Status == "Rejected" || input.Status == "Rework") {
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Printf("Error loading config: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
			return
		}

		spcFilters := model.GetSpCFilters{
			ID: strconv.Itoa(spcReview.SpcID),
		}
		spcs, err := repository.GetSpCs(spcFilters)
		if err != nil {
			log.Printf("Error fetching spc details: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spc details"})
			return
		}

		if len(spcs) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Spc not found"})
			return
		}

		var spc = spcs[0]

		submissionURL := fmt.Sprintf("%s/submissions?id=%d", cfg.SubmissionServiceURL, spcReview.SubmissionID)
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

		err = util.SendEmail(spc.Email, submission.OfficialMailID, subject, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Spc review updated successfully",
	})
}

func GetSpcReviewsHandler(c *gin.Context) {
	submissionID := c.DefaultQuery("submission_id", "")
	spcID := c.DefaultQuery("spc_id", "")
	status := c.DefaultQuery("status", "")
	reviewID := c.DefaultQuery("review_id", "")

	filters := model.GetSpcReviewFilters{}
	if reviewID != "" {
		filters.ID = reviewID
	}
	if submissionID != "" {
		filters.SubmissionID = submissionID
	}
	if spcID != "" {
		filters.SpcID = spcID
	}
	if status != "" {
		filters.Status = status
	}

	reviews, err := repository.GetSpcReviews(filters)
	if err != nil {
		log.Printf("Error fetching reviews: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
