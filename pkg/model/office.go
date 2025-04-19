package model

import (
	"time"

	"github.com/google/uuid"
)

type Office struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Department string    `json:"department" db:"department"`
	Email      string    `json:"email" db:"email"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type GetOfficeFilters struct {
	ID         string `form:"id"`
	Department string `form:"department"`
	Email      string `form:"email"`
}
