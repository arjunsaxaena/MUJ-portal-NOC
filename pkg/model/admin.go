package model

import "time"

type Admin struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type GetAdminFilters struct {
	ID    string `form:"id"`
	Email string `form:"Email"`
}
