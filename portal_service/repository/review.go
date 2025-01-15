package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"fmt"
	"time"
)

func CreateReview(submissionID, reviewerID int, status, comments string) (int, error) {
	query := `
		INSERT INTO reviewer_reviews (submission_id, reviewer_id, status, comments, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var reviewID int
	err := database.DB.QueryRow(query, submissionID, reviewerID, status, comments, time.Now(), time.Now()).Scan(&reviewID)
	if err != nil {
		return 0, err
	}
	return reviewID, nil
}

func UpdateReview(reviewID int, status, comments string) error {
	query := `
		UPDATE reviewer_reviews
		SET status = $1, comments = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := database.DB.Exec(query, status, comments, time.Now(), reviewID)
	return err
}

func GetReviews(filters model.GetReviewFilters) ([]model.Review, error) {
	query := `
		SELECT id, submission_id, reviewer_id, status, comments, created_at, updated_at
		FROM reviewer_reviews
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
	if filters.ReviewerID != "" {
		query += fmt.Sprintf(" AND reviewer_id = $%d", argIndex)
		args = append(args, filters.ReviewerID)
		argIndex++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filters.Status)
		argIndex++
	}

	var reviews []model.Review
	err := database.DB.Select(&reviews, query, args...)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
