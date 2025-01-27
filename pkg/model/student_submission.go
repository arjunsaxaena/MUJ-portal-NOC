package model

import "time"

type StudentSubmission struct {
	ID                  int       `json:"id" db:"id"`
	RegistrationNumber  string    `json:"registration_number" form:"registrationNumber" binding:"required" db:"registration_number"`
	Name                string    `json:"name" form:"name" binding:"required" db:"name"`
	Gender              string    `json:"gender" form:"gender" db:"gender"`
	Semester            string    `json:"semester" form:"semester" db:"semester"`
	OfficialMailID      string    `json:"official_mail_id" form:"email" binding:"required,email" db:"official_mail_id"`
	MobileNumber        string    `json:"mobile_number" form:"mobile" binding:"required" db:"mobile_number"`
	Department          string    `json:"department" form:"department" binding:"required" db:"department"`
	Section             string    `json:"section" form:"section" binding:"required" db:"section"`
	OfferType           string    `json:"offer_type" form:"offerType" binding:"required" db:"offer_type"`
	CompanyName         string    `json:"company_name" form:"companyName" binding:"required" db:"company_name"`
	CompanyState        string    `json:"company_state" form:"companyState" binding:"required" db:"company_state"`
	CompanyCity         string    `json:"company_city" form:"companyCity" binding:"required" db:"company_city"`
	Pincode             string    `json:"pincode" form:"companyPin" binding:"required" db:"pincode"`
	HRDEmail            *string   `json:"hrd_email" form:"hrdEmail" binding:"omitempty,email" db:"hrd_email"`
	HRDNumber           *string   `json:"hrd_number" form:"hrdNumber" db:"hrd_number"`
	HasOfferLetter      bool      `json:"has_offer_letter" form:"hasOfferLetter" db:"has_offer_letter"`
	OfferTypeDetail     string    `json:"offer_type_detail" form:"internshipType" binding:"required" db:"offer_type_detail"`
	PackagePPO          string    `json:"package_ppo" form:"ppoPackage" binding:"omitempty" db:"package_ppo"`
	StipendAmount       string    `json:"stipend_amount" form:"stipend" binding:"required" db:"stipend_amount"`
	InternshipStartDate string    `json:"internship_start_date" form:"startDate" binding:"required" db:"internship_start_date"`
	InternshipEndDate   string    `json:"internship_end_date" form:"endDate" binding:"required" db:"internship_end_date"`
	OfferLetterPath     string    `json:"offer_letter_path" form:"offerLetterPath" binding:"omitempty" db:"offer_letter_path"`
	MailCopyPath        string    `json:"mail_copy_path" form:"mailCopyPath" binding:"omitempty" db:"mail_copy_path"`
	TermsAccepted       bool      `json:"terms_accepted" form:"termsAccepted" binding:"required" db:"terms_accepted"`
	Status              string    `json:"status" db:"status"`
	NocPath             *string   `json:"noc_path" db:"noc_path"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type GetSubmissionFilters struct {
	ID         string `form:"id"`
	Department string `form:"department"`
	Status     string `form:"status"`
}
