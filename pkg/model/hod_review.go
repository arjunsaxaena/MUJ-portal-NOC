package model

import (
	"time"
)

type HodReview struct {
	ID           int       `json:"id" db:"id"`
	SubmissionID int       `json:"submission_id" db:"submission_id"`
	HodID        int       `json:"hod_id" db:"hod_id"`
	Action       string    `json:"action" db:"action"` // e.g., Approved, Rejected
	Remarks      string    `json:"remarks" db:"remarks"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// sql.NullString  // NullString represents a string that may be null
