package model

import (
	"time"
)

type FpcReview struct {
	ID           string    `json:"id" db:"id"`
	SubmissionID int       `json:"submission_id" db:"submission_id"`
	FpcID        int       `json:"fpc_id" db:"fpc_id"`
	Status       string    `json:"status" db:"status"`
	Comments     string    `json:"comments" db:"comments"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type GetFpcReviewFilters struct {
	ID           string `form:"id"`
	SubmissionID string `form:"submission_id"`
	FpcID        string `form:"fpc_id"`
	Status       string `form:"status"`
}
