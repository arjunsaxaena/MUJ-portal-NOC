package database

import (
	"MUJ_automated_mail_generation/pkg/model"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateReviewer(username, passwordHash, department string) (int, error) {
	query := `
		INSERT INTO reviewers (username, password_hash, department)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int
	err := DB.QueryRow(query, username, passwordHash, department).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetReviewerByUsername(username string) (*model.Reviewer, error) {
	query := `
		SELECT id, username, password_hash, department, created_at
		FROM reviewers
		WHERE username = $1
	`
	var reviewer model.Reviewer
	err := DB.Get(&reviewer, query, username)
	if err != nil {
		fmt.Printf("Error fetching reviewer: %v\n", err)
		return nil, err
	}
	fmt.Printf("Fetched reviewer: %+v\n", reviewer)
	return &reviewer, nil
}

func ValidateReviewerPassword(username, password string) (bool, error) {
	reviewer, err := GetReviewerByUsername(username)
	if err != nil {
		return false, err
	}

	// Compare the hashed password.
	err = bcrypt.CompareHashAndPassword([]byte(reviewer.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("invalid credentials")
	}
	return true, nil
}

func GetAllReviewers() ([]model.Reviewer, error) {
	var reviewers []model.Reviewer
	err := DB.Select(&reviewers, "SELECT * FROM reviewers")
	return reviewers, err
}
