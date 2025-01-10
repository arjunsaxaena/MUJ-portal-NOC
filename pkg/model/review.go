package model

import (
	"database/sql"
	"time"
)

type Review struct {
	ID           int            `json:"id" db:"id"`
	SubmissionID int            `json:"submission_id" db:"submission_id"`
	ReviewerID   int            `json:"reviewer_id" db:"reviewer_id"`
	Status       string         `json:"status" db:"status"`
	Comments     string         `json:"comments" db:"comments"`
	HodAction    sql.NullString `json:"hod_action" db:"hod_action"` // NullString represents a string that may be null
	HodRemarks   sql.NullString `json:"hod_remarks" db:"hod_remarks"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}
