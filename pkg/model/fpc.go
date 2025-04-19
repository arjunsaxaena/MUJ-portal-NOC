package model

import "time"

type FpC struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	AppPassword  string    `json:"app_password" db:"app_password"`
	Department   string    `json:"department" db:"department"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type GetFpCFilters struct {
	ID         string `form:"id"`
	Department string `form:"department"`
	Email      string `form:"Email"`
}
