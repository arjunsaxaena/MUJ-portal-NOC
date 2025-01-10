package model

import "time"

type StudentSubmission struct {
	ID                  int       `json:"id" db:"id"`
	RegistrationNumber  string    `json:"registration_number" form:"registration_number" binding:"required" db:"registration_number"`
	Name                string    `json:"name" form:"name" binding:"required" db:"name"`
	OfficialMailID      string    `json:"official_mail_id" form:"official_mail_id" binding:"required,email" db:"official_mail_id"`
	MobileNumber        string    `json:"mobile_number" form:"mobile_number" binding:"required" db:"mobile_number"`
	Department          string    `json:"department" form:"department" binding:"required" db:"department"`
	Section             string    `json:"section" form:"section" binding:"required" db:"section"`
	OfferType           string    `json:"offer_type" form:"offer_type" binding:"required" db:"offer_type"`
	CompanyName         string    `json:"company_name" form:"company_name" db:"company_name"`
	CompanyAddress      string    `json:"company_address" form:"company_address" db:"company_address"`
	OfferTypeDetail     string    `json:"offer_type_detail" form:"offer_type_detail" db:"offer_type_detail"`
	PackagePPO          string    `json:"package_ppo" form:"package_ppo" db:"package_ppo"`
	StipendAmount       string    `json:"stipend_amount" form:"stipend_amount" db:"stipend_amount"`
	InternshipStartDate string    `json:"internship_start_date" form:"internship_start_date" db:"internship_start_date"`
	InternshipEndDate   string    `json:"internship_end_date"  form:"internship_end_date" db:"internship_end_date"`
	OfferLetterPath     string    `json:"offer_letter_path" db:"offer_letter_path"`
	MailCopyPath        string    `json:"mail_copy_path" db:"mail_copy_path"`
	TermsAccepted       bool      `json:"terms_accepted" form:"terms_accepted" binding:"required" db:"terms_accepted"`
	Status              string    `json:"status" db:"status"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
