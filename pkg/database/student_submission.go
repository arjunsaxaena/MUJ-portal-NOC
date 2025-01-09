package database

import (
	"MUJ_automated_mail_generation/pkg/model"
	"fmt"
)

func CreateSubmission(submission *model.StudentSubmission) error {
	_, err := DB.NamedExec(`
		INSERT INTO student_submissions (
			registration_number, name, official_mail_id, mobile_number, department, section, offer_type, 
			company_name, company_address, offer_type_detail, package_ppo, stipend_amount, 
			internship_start_date, internship_end_date, offer_letter_path, mail_copy_path, terms_accepted, 
			status, created_at
		) 
		VALUES (
			:registration_number, :name, :official_mail_id, :mobile_number, :department, :section, :offer_type, 
			:company_name, :company_address, :offer_type_detail, :package_ppo, :stipend_amount, 
			:internship_start_date, :internship_end_date, :offer_letter_path, :mail_copy_path, :terms_accepted, 
			:status, CURRENT_TIMESTAMP
		)
	`, submission)

	if err != nil {
		fmt.Println("Error inserting submission:", err)
		return err
	}
	return nil
}

func GetAllSubmissions() ([]model.StudentSubmission, error) {
	var submissions []model.StudentSubmission
	err := DB.Select(&submissions, "SELECT * FROM student_submissions")
	return submissions, err
}

func GetSubmissionsByDepartment(department string) ([]model.StudentSubmission, error) {
	var submissions []model.StudentSubmission
	err := DB.Select(&submissions, "SELECT * FROM student_submissions WHERE department = ?", department)
	return submissions, err
} // Have to implement handler

func UpdateSubmissionStatus(submissionID int, status, remarks string) error {
	query := `
		UPDATE student_submissions
		SET status = $1, remarks = $2
		WHERE id = $3
	`
	_, err := DB.Exec(query, status, remarks, submissionID)
	return err
} // Have to implement handler
