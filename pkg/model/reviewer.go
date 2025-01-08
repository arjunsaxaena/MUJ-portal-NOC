package model

import "time"

type Reviewer struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	Department   string    `json:"department" db:"department"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
