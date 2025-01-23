package model

import "time"

type SpC struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	Department   string    `json:"department" db:"department"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type GetSpCFilters struct {
	ID         string `form:"id"`
	Department string `form:"department"`
	Email      string `form:"Email"`
}
