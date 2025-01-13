package model

import "time"

type StudentSubmission struct {
	ID                  int       `json:"id" db:"id"`
	RegistrationNumber  string    `json:"registration_number" form:"registrationNumber" binding:"required" db:"registration_number"`
	Name                string    `json:"name" form:"name" binding:"required" db:"name"`
	OfficialMailID      string    `json:"official_mail_id" form:"email" binding:"required,email" db:"official_mail_id"`
	MobileNumber        string    `json:"mobile_number" form:"mobile" binding:"required" db:"mobile_number"`
	Department          string    `json:"department" form:"department" binding:"required" db:"department"`
	Section             string    `json:"section" form:"section" binding:"required" db:"section"`
	OfferType           string    `json:"offer_type" form:"offerType" binding:"required" db:"offer_type"`
	CompanyName         string    `json:"company_name" form:"companyName" db:"company_name"`
	CompanyState        string    `json:"company_state" form:"companyState" binding:"required" db:"company_state"`
	CompanyCity         string    `json:"company_city" form:"companyCity" binding:"required" db:"company_city"`
	Pincode             string    `json:"pincode" form:"companyPin" binding:"required" db:"pincode"`
	OfferTypeDetail     string    `json:"offer_type_detail" form:"internshipType" db:"offer_type_detail"`
	PackagePPO          string    `json:"package_ppo" form:"ppoPackage" db:"package_ppo"`
	StipendAmount       string    `json:"stipend_amount" form:"stipend" db:"stipend_amount"`
	InternshipStartDate string    `json:"internship_start_date" form:"startDate" db:"internship_start_date"`
	InternshipEndDate   string    `json:"internship_end_date" form:"endDate" db:"internship_end_date"`
	OfferLetterPath     string    `json:"offer_letter_path" db:"offer_letter_path"`
	MailCopyPath        string    `json:"mail_copy_path" db:"mail_copy_path"`
	TermsAccepted       bool      `json:"terms_accepted" form:"termsAccepted" binding:"required" db:"terms_accepted"`
	Status              string    `json:"status" db:"status"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
