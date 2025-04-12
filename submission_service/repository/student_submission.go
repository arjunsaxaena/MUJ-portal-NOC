package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"fmt"
	"strconv"
	"time"
)

func CreateSubmission(submission *model.StudentSubmission) error {
	_, err := database.DB.NamedExec(`
		INSERT INTO student_submissions (
			registration_number, name, gender, semester, official_mail_id, mobile_number, department, section, 
			offer_type, cgpa, backlogs, company_name, company_state, company_city, pincode, offer_type_detail, package_ppo, stipend_amount, 
			internship_start_date, internship_end_date, offer_letter_path, mail_copy_path, hrd_email, terms_accepted, 
			status, noc_type, noc_path, created_at, updated_at
		)
		VALUES (
			:registration_number, :name, :gender, :semester, :official_mail_id, :mobile_number, :department, :section, 
			:offer_type, :cgpa, :backlogs, :company_name, :company_state, :company_city, :pincode, :offer_type_detail, :package_ppo, :stipend_amount, 
			:internship_start_date, :internship_end_date, :offer_letter_path, :mail_copy_path, :hrd_email, :terms_accepted, 
			:status, :noc_type, :noc_path, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`, submission)

	return err
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
	if filters.NocType != "" {
		query += " AND noc_type = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.NocType)
		paramIndex++
	}

	var submissions []model.StudentSubmission
	err := database.DB.Select(&submissions, query, args...)
	return submissions, err
}

func UpdateSubmission(submissionID, status, nocPath string) error {
	query := `UPDATE student_submissions
	          SET updated_at = $1`
	args := []interface{}{time.Now()}

	if status != "" {
		query += `, status = $2`
		args = append(args, status)
	}

	if nocPath != "" {
		query += `, noc_path = $3`
		args = append(args, nocPath)
	}

	query += ` WHERE id = $` + fmt.Sprint(len(args)+1)

	args = append(args, submissionID)

	_, err := database.DB.Exec(query, args...)
	return err
}

func DeleteSubmission(submissionID string) error {
	query := `
		DELETE FROM student_submissions
		WHERE id = $1
	`

	_, err := database.DB.Exec(query, submissionID)
	return err
}
