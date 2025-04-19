package model

import "time"

type Student struct {
	ID                 string    `json:"id" db:"id"`
	RegistrationNumber string    `json:"registration_number" form:"registrationNumber" binding:"required" db:"registration_number"`
	Name               string    `json:"name" form:"name" binding:"required" db:"name"`
	OfficialMailID     string    `json:"official_mail_id" form:"email" binding:"required,email" db:"official_mail_id"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type GetStudentFilters struct {
	ID                 string `form:"id"`
	RegistrationNumber string `form:"registration_number"`
	EmailID            string `form:"email_id"`
	Department         string `form:"department"`
}
