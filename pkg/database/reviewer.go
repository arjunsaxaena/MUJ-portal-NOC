package database

import (
	"MUJ_automated_mail_generation/pkg/model"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateReviewer(email, passwordHash, department string) (int, error) {
	query := `
		INSERT INTO reviewers (email, password_hash, department)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int
	err := DB.QueryRow(query, email, passwordHash, department).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetReviewerByEmail(email string) (*model.Reviewer, error) {
	query := `
		SELECT id, email, password_hash, department, created_at
		FROM reviewers
		WHERE email = $1
	`
	var reviewer model.Reviewer
	err := DB.Get(&reviewer, query, email)
	if err != nil {
		fmt.Printf("Error fetching reviewer: %v\n", err)
		return nil, err
	}
	fmt.Printf("Fetched reviewer: %+v\n", reviewer)
	return &reviewer, nil
}

func GetReviewerByDepartment(department string) ([]model.Reviewer, error) {
	query := `
		SELECT id, email, password_hash, department, created_at
		FROM reviewers
		WHERE department = $1
	`
	var reviewers []model.Reviewer
	err := DB.Select(&reviewers, query, department)
	if err != nil {
		fmt.Printf("Error fetching reviewers: %v\n", err)
		return nil, err
	}
	fmt.Printf("Fetched reviewers: %+v\n", reviewers)
	return reviewers, nil
}

func ValidateReviewerPassword(email, password string) (bool, error) {
	reviewer, err := GetReviewerByEmail(email)
	if err != nil {
		return false, err
	}

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
