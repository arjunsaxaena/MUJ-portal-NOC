package model

import (
	"time"
)

type HodReview struct {
	ID           string    `json:"id" db:"id"`
	SubmissionID int       `json:"submission_id" db:"submission_id"`
	HodID        int       `json:"hod_id" db:"hod_id"`
	Action       string    `json:"action" db:"action"`
	Remarks      string    `json:"remarks" db:"remarks"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type GetHodReviewFilters struct {
	ID           string `json:"id" db:"id"`
	SubmissionID string `form:"submission_id"`
	HodID        string `form:"hod_id"`
	Action       string `form:"action"`
}
