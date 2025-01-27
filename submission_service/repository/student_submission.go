package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"fmt"
	"log"
	"strconv"
	"time"
)

func CreateSubmission(submission *model.StudentSubmission) error {
	_, err := database.DB.NamedExec(`
		INSERT INTO student_submissions (
			registration_number, name, gender, semester, official_mail_id, mobile_number, department, section, 
			offer_type, company_name, company_state, company_city, pincode, offer_type_detail, package_ppo, stipend_amount, 
			internship_start_date, internship_end_date, has_offer_letter, offer_letter_path, mail_copy_path, hrd_email, terms_accepted, 
			status, created_at, updated_at
		)
		VALUES (
			:registration_number, :name, :gender, :semester, :official_mail_id, :mobile_number, :department, :section, 
			:offer_type, :company_name, :company_state, :company_city, :pincode, :offer_type_detail, :package_ppo, :stipend_amount, 
			:internship_start_date, :internship_end_date, :has_offer_letter, :offer_letter_path, :mail_copy_path, :hrd_email, :terms_accepted, 
			:status, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`, submission)

	if err != nil {
		fmt.Println("Error inserting submission:", err)
		return err
	}
	return nil
}

func GetSubmissions(filters model.GetSubmissionFilters) ([]model.StudentSubmission, error) {
	query := "SELECT * FROM student_submissions WHERE 1=1"
	var args []interface{}
	paramIndex := 1

	if filters.ID != "" {
		query += " AND id = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.ID)
		paramIndex++
	}
	if filters.Department != "" {
		query += " AND department = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.Department)
		paramIndex++
	}
	if filters.Status != "" {
		query += " AND status = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.Status)
		paramIndex++
	}

	var submissions []model.StudentSubmission
	err := database.DB.Select(&submissions, query, args...)
	if err != nil {
		log.Printf("Error fetching submissions with filters %v: %v", filters, err)
		return nil, err
	}

	return submissions, nil
}

func UpdateSubmission(submissionID int, status string, nocPath string) error {
	// Construct the query dynamically based on the provided fields
	query := `UPDATE student_submissions
	          SET updated_at = $1`
	args := []interface{}{time.Now()}

	// If status is provided, include it in the query
	if status != "" {
		query += `, status = $2`
		args = append(args, status)
	}

	// If noc_path is provided, include it in the query
	if nocPath != "" {
		query += `, noc_path = $3`
		args = append(args, nocPath)
	}

	// Ensure the WHERE condition is always included
	query += ` WHERE id = $` + fmt.Sprint(len(args)+1)

	// Append submissionID to the arguments
	args = append(args, submissionID)

	// Execute the query
	_, err := database.DB.Exec(query, args...)
	if err != nil {
		fmt.Printf("Error updating submission: %v\n", err)
		return err
	}

	return nil
}

func DeleteSubmission(submissionID int) error {
	query := `
		DELETE FROM student_submissions
		WHERE id = $1
	`

	_, err := database.DB.Exec(query, submissionID)
	if err != nil {
		fmt.Printf("Error deleting submission: %v\n", err)
		return err
	}

	return nil
}
