package model

import "time"

type HoD struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	AppPassword  string    `json:"app_password" db:"app_password"`
	RoleType     string    `json:"role_type" db:"role_type"`
	Department   string    `json:"department" db:"department"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type GetHoDFilters struct {
	ID         string `form:"id"`
	Department string `form:"department"`
	RoleType   string `form:"role_type"`
	Email      string `form:"Email"`
}
