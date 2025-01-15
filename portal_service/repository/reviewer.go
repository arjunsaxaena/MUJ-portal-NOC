package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"errors"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func CreateReviewer(name, email, passwordHash, department string) (int, error) {
	query := `
		INSERT INTO reviewers (name, email, password_hash, department)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err := database.DB.QueryRow(query, name, email, passwordHash, department).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetReviewers(filters model.GetReviewerFilters) ([]model.Reviewer, error) {
	query := "SELECT * FROM reviewers WHERE 1=1"
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

	var reviewers []model.Reviewer
	err := database.DB.Select(&reviewers, query, args...)
	if err != nil {
		log.Printf("Error fetching reviewers with filters %v: %v", filters, err)
		return nil, err
	}

	return reviewers, nil
}

func ValidateReviewerPassword(email, password string) (bool, error) {
	filters := model.GetReviewerFilters{Email: email}
	reviewers, err := GetReviewers(filters)
	if err != nil {
		return false, err
	}

	if len(reviewers) != 1 {
		return false, errors.New("invalid credentials")
	}

	reviewer := reviewers[0]

	err = bcrypt.CompareHashAndPassword([]byte(reviewer.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("invalid credentials")
	}
	return true, nil
}

func UpdateReviewer(id int, name, email, passwordHash, department string) error {
	query := `
		UPDATE reviewers
		SET name = COALESCE($1, name),
		    email = COALESCE($2, email),
		    password_hash = COALESCE($3, password_hash),
		    department = COALESCE($4, department)
		WHERE id = $5
	`
	_, err := database.DB.Exec(query, name, email, passwordHash, department, id)
	if err != nil {
		log.Printf("Error updating reviewer with ID %d: %v", id, err)
		return err
	}

	return nil
} // For testing

func DeleteReviewer(id int) error {
	query := `
		DELETE FROM reviewers
		WHERE id = $1
	`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting reviewer with ID %d: %v", id, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no reviewer found with the given ID")
	}

	return nil
} // For testing
