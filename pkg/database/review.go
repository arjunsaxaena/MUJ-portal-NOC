package database

import (
	"MUJ_automated_mail_generation/pkg/model"
	"time"
)

func CreateReviewerReview(submissionID, reviewerID int, status, comments string) (int, error) {
	query := `
		INSERT INTO reviewer_reviews (submission_id, reviewer_id, status, comments, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var reviewID int
	err := DB.QueryRow(query, submissionID, reviewerID, status, comments, time.Now(), time.Now()).Scan(&reviewID)
	if err != nil {
		return 0, err
	}
	return reviewID, nil
}

func UpdateReviewerReview(reviewID int, status, comments string) error {
	query := `
		UPDATE reviewer_reviews
		SET status = $1, comments = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := DB.Exec(query, status, comments, time.Now(), reviewID)
	return err
}

func GetReviewerReviewsBySubmission(submissionID int) ([]model.ReviewerReview, error) {
	query := `
		SELECT id, submission_id, reviewer_id, status, comments, created_at, updated_at
		FROM reviewer_reviews
		WHERE submission_id = $1
	`
	var reviews []model.ReviewerReview
	err := DB.Select(&reviews, query, submissionID)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func GetReviewerReviewsByReviewer(reviewerID int) ([]model.ReviewerReview, error) {
	query := `
		SELECT id, submission_id, reviewer_id, status, comments, created_at, updated_at
		FROM reviewer_reviews
		WHERE reviewer_id = $1
	`
	var reviews []model.ReviewerReview
	err := DB.Select(&reviews, query, reviewerID)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func GetAllReviewerReviews() ([]model.ReviewerReview, error) {
	query := `
		SELECT id, submission_id, reviewer_id, status, comments, created_at, updated_at
		FROM reviewer_reviews
	`
	var reviews []model.ReviewerReview
	err := DB.Select(&reviews, query)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
