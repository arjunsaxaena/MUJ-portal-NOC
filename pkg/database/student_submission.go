package database

import (
	"MUJ_automated_mail_generation/pkg/model"
	"fmt"
	"log"
	"time"
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
	query := "SELECT * FROM student_submissions WHERE department = $1"
	err := DB.Select(&submissions, query, department)
	if err != nil {
		log.Printf("Error executing query: %v", err)
	}
	return submissions, err
}

func GetApprovedSubmissionsByDepartment(department string) ([]model.StudentSubmission, error) {
	var submissions []model.StudentSubmission
	query := `
		SELECT * 
		FROM student_submissions 
		WHERE department = $1 
		AND status = 'Approved'
	`

	err := DB.Select(&submissions, query, department)
	if err != nil {
		log.Printf("Error fetching approved submissions for department %s: %v", department, err)
		return nil, err
	}

	return submissions, nil
}

func UpdateSubmissionStatus(submissionID int, status string) error {
	query := `
		UPDATE student_submissions
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	// Execute the query to update the status
	_, err := DB.Exec(query, status, time.Now(), submissionID)
	if err != nil {
		fmt.Printf("Error updating submission status: %v\n", err)
		return err
	}

	return nil
}
