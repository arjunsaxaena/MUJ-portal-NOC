package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"fmt"
	"log"
	"time"
)

func CreateHodReview(submissionID, hodID, action, remarks string) (string, error) {
	query := `
		INSERT INTO hod_reviews (submission_id, hod_id, action, remarks, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var reviewID string
	now := time.Now()
	log.Printf("Inserting HoD review: submissionID=%s, hodID=%s, action=%s", submissionID, hodID, action)

	err := database.DB.QueryRow(query, submissionID, hodID, action, remarks, now, now).Scan(&reviewID)
	if err != nil {
		log.Printf("Error creating HoD review: %v", err)
		return "", err
	}

	log.Printf("Successfully created HoD review with ID: %s", reviewID)
	return reviewID, nil
}

func UpdateHodReview(reviewID int, action, remarks string) error {
	query := `
		UPDATE hod_reviews
		SET action = $1, remarks = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := database.DB.Exec(query, action, remarks, time.Now(), reviewID)
	return err
}

func GetHodReviews(filters model.GetHodReviewFilters) ([]model.HodReview, error) {
	query := `
		SELECT id, submission_id, hod_id, action, remarks, created_at, updated_at
		FROM hod_reviews
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if filters.ID != "" {
		query += fmt.Sprintf(" AND id = $%d", argIndex)
		args = append(args, filters.ID)
		argIndex++
	}
	if filters.SubmissionID != "" {
		query += fmt.Sprintf(" AND submission_id = $%d", argIndex)
		args = append(args, filters.SubmissionID)
		argIndex++
	}
	if filters.HodID != "" {
		query += fmt.Sprintf(" AND hod_id = $%d", argIndex)
		args = append(args, filters.HodID)
		argIndex++
	}
	if filters.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIndex)
		args = append(args, filters.Action)
		argIndex++
	}

	var reviews []model.HodReview
	err := database.DB.Select(&reviews, query, args...)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
