package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func CreateOffice(department, email string) (string, error) {
	query := `
		INSERT INTO office (department, email)
		VALUES ($1, $2)
		RETURNING id
	`
	var id string
	err := database.DB.QueryRow(query, department, email).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetOffices(filters model.GetOfficeFilters) ([]model.Office, error) {
	query := "SELECT * FROM office WHERE 1=1"
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
	if filters.Email != "" {
		query += " AND email = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.Email)
		paramIndex++
	}

	var offices []model.Office
	err := database.DB.Select(&offices, query, args...)
	if err != nil {
		log.Printf("Error fetching offices with filters %v: %v", filters, err)
		return nil, err
	}

	return offices, nil
}

func UpdateOffice(id string, department *string, email *string) error {
	query := "UPDATE office SET "
	args := []interface{}{}
	argCount := 1

	if department != nil {
		query += "department = $" + fmt.Sprint(argCount) + ", "
		args = append(args, *department)
		argCount++
	}

	if email != nil {
		query += "email = $" + fmt.Sprint(argCount) + ", "
		args = append(args, *email)
		argCount++
	}

	// Add updated_at timestamp
	query += "updated_at = NOW() "

	if argCount > 1 {
		query = strings.TrimSuffix(query, ", ")
		query += " WHERE id = $" + fmt.Sprint(argCount)
	} else {
		query += "WHERE id = $" + fmt.Sprint(argCount)
	}

	args = append(args, id)

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Error updating office with ID %s: %v", id, err)
		return err
	}

	return nil
}

func DeleteOffice(id string) error {
	query := `
		DELETE FROM office
		WHERE id = $1
	`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting office with ID %s: %v", id, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no office found with the given ID")
	}

	return nil
}
