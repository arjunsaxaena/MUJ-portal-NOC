package model

import "time"

type Review struct {
	ID           int       `json:"id" db:"id"`
	SubmissionID int       `json:"submission_id" db:"submission_id"`
	ReviewerID   int       `json:"reviewer_id" db:"reviewer_id"`
	Status       string    `json:"status" db:"status"`
	Comments     string    `json:"comments" db:"comments"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
