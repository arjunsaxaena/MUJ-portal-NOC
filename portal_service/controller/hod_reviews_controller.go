package controller

import (
	"MUJ_AMG/pkg/database"
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

	var hod model.HoD
	err = database.DB.Get(&hod, "SELECT * FROM hod WHERE id = $1", input.HodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch HoD details"})
		return
	}

	if input.Action == "Rejected" || input.Action == "Rework" {
		updateSubmissionURL := fmt.Sprintf("%s/submissions", cfg.SubmissionServiceURL)
		updateSubmissionBody := struct {
			Status string `json:"status"`
		}{
			Status: input.Action,
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

		subject := "Your Placement Application Status - HoD Review"
		body := fmt.Sprintf("Dear %s,\n\nYour placement application has been %s.\n\nHoD Comments: %s\n\nBest regards",
			submission.Name, input.Action, input.Remarks)

		err = util.SendEmail(hod.Email, submission.OfficialMailID, subject, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email to student"})
			return
		}
	}

	if input.Action == "Approved" {
		nocURL, err := util.CreateNocPdf(submission, "muj-student-data", "generated_pdfs")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate NOC PDF"})
			return
		}

		subject := "Your No Objection Certificate (NOC)"
		body := fmt.Sprintf("Dear %s,\n\nYour placement application has been approved. Please find attached the No Objection Certificate (NOC) for your reference.\n\nBest regards,\nHoD", submission.Name)

		err = util.SendEmailWithAttachment(hod.Email, submission.OfficialMailID, subject, body, nocURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send NOC email to student"})
			return
		}
		updateSubmissionURL := fmt.Sprintf("%s/submissions", cfg.SubmissionServiceURL)
		updateSubmissionBody := struct {
			Status string `json:"status"`
		}{
			Status: "NOC sent",
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

	reviewID, err := repository.CreateHodReview(input.SubmissionID, input.HodID, input.Action, input.Remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HoD review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"review_id": reviewID, "message": "HoD review created successfully"})
}

func GetHodReviewsHandler(c *gin.Context) {
	var filters model.GetHodReviewFilters
	err := c.ShouldBindQuery(&filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	reviews, err := repository.GetHodReviews(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch HoD reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
