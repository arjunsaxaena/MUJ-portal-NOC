package model

import (
	"time"
)

type ReviewerReview struct {
	ID           int       `json:"id" db:"id"`
	SubmissionID int       `json:"submission_id" db:"submission_id"`
	ReviewerID   int       `json:"reviewer_id" db:"reviewer_id"`
	Status       string    `json:"status" db:"status"` // e.g., Accepted, Rejected, Rework
	Comments     string    `json:"comments" db:"comments"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// sql.NullString  // NullString represents a string that may be null
