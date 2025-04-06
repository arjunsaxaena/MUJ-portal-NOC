package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"fmt"
	"time"
)

func CreateFpcReview(submissionID, fpcID, status, comments string) (string, error) {
	query := `
		INSERT INTO fpc_reviews (submission_id, fpc_id, status, comments, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var reviewID string
	err := database.DB.QueryRow(query, submissionID, fpcID, status, comments, time.Now(), time.Now()).Scan(&reviewID)
	if err != nil {
		return "", err
	}
	return reviewID, err
}

func UpdateFpcReview(reviewID, status, comments string) error {
	query := `
		UPDATE fpc_reviews
		SET status = $1, comments = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := database.DB.Exec(query, status, comments, time.Now(), reviewID)
	return err
}

func GetFpcReviews(filters model.GetFpcReviewFilters) ([]model.FpcReview, error) {
	query := `
		SELECT id, submission_id, fpc_id, status, comments, created_at, updated_at
		FROM fpc_reviews
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
	if filters.FpcID != "" {
		query += fmt.Sprintf(" AND fpc_id = $%d", argIndex)
		args = append(args, filters.FpcID)
		argIndex++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filters.Status)
		argIndex++
	}

	var reviews []model.FpcReview
	err := database.DB.Select(&reviews, query, args...)
	return reviews, err
}
