package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"log"
	"strconv"
)

func GetStudents(filters model.GetStudentFilters) ([]model.Student, error) {
	query := "SELECT * FROM students WHERE 1=1"
	var args []interface{}
	paramIndex := 1

	if filters.RegistrationNumber != "" {
		query += " AND registration_number = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.RegistrationNumber)
		paramIndex++
	}

	if filters.EmailID != "" {
		query += " AND email_id = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.EmailID)
		paramIndex++
	}

	if filters.Department != "" {
		query += " AND department = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.Department)
		paramIndex++
	}

	var students []model.Student
	err := database.DB.Select(&students, query, args...)
	if err != nil {
		log.Printf("Error fetching students with filters %v: %v", filters, err)
		return nil, err
	}

	return students, nil
}
