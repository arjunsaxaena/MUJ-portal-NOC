package controller

import (
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/pkg/util"
	"MUJ_AMG/portal_service/repository"
	submissionRepository "MUJ_AMG/submission_service/repository"
	"fmt"
	"net/http"
	"path/filepath"
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

	hodFilters := model.GetHoDFilters{
		ID: strconv.Itoa(input.HodID),
	}
	hods, err := repository.GetHoDs(hodFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch HoD details"})
		return
	}

	if len(hods) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "HoD not found"})
		return
	}

	hod := hods[0]

	submissionFilters := model.GetSubmissionFilters{
		ID: strconv.Itoa(input.SubmissionID),
	}
	submissions, err := submissionRepository.GetSubmissions(submissionFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submission"})
		return
	}

	if len(submissions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	submission := submissions[0]

	if input.Action == "Rejected" {
		err := submissionRepository.UpdateSubmission(input.SubmissionID, input.Action, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission status"})
			return
		}

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
		var nocPath string
		var err error
		if submission.NocType == "Generic" {
			nocPath, err = util.CreateGenericNocPdf(submission)
		} else {
			nocPath, err = util.CreateNocPdf(submission)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate NOC PDF"})
			return
		}

		nocPath = filepath.Base(nocPath)

		subject := "Your Placement Application Status - NOC Available"
		body := fmt.Sprintf("Dear %s,\n\nYour placement application has been approved. Please collect your No Objection Certificate (NOC) from the %s office \n\nBest regards,\nHOD",
			submission.Name, submission.Department)

		err = util.SendEmail(hod.Email, submission.OfficialMailID, subject, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email to student"})
			return
		}

		officeSubject := "NOC Generated"
		officeBody := fmt.Sprintf("NOC for %s", submission.Name)

		err = util.SendEmailWithAttachment(hod.Email, "arjunsaxena04@gmail.com", officeSubject, officeBody, nocPath) // change arjunsaxena04@gmail.com to office.cse@jaipur.manipal.edu
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send NOC to Office"})
			return
		}

		err = submissionRepository.UpdateSubmission(input.SubmissionID, "NOC ready", nocPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "NOC generated successfully", "noc_path": nocPath})
		return
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
