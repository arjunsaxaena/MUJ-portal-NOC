package database

import (
	"MUJ_automated_mail_generation/pkg/model"
	"fmt"
	"time"
)

func CreateReview(submissionID, reviewerID int, status, comments string) (int, error) {
	query := `
		INSERT INTO reviews (submission_id, reviewer_id, status, comments, created_at, updated_at)
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

func UpdateReview(reviewID int, status, comments string) error {
	query := `
		UPDATE reviews
		SET status = $1, comments = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := DB.Exec(query, status, comments, time.Now(), reviewID)
	return err
}

func UpdateHodAction(reviewID int, hodAction, hodRemarks string) error {
	query := `
		UPDATE reviews
		SET hod_action = $1, hod_remarks = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := DB.Exec(query, hodAction, hodRemarks, time.Now(), reviewID)
	return err
}

func GetReviewsBySubmission(submissionID int) ([]model.Review, error) {
	query := `
		SELECT id, submission_id, reviewer_id, status, comments, hod_action, hod_remarks, created_at, updated_at
		FROM reviews
		WHERE submission_id = $1
	`
	var reviews []model.Review
	err := DB.Select(&reviews, query, submissionID)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func GetReviewsByReviewer(reviewerID int) ([]model.Review, error) {
	query := `
		SELECT id, submission_id, reviewer_id, status, comments, hod_action, hod_remarks, created_at, updated_at
		FROM reviews
		WHERE reviewer_id = $1
	`
	var reviews []model.Review
	err := DB.Select(&reviews, query, reviewerID)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func GetAllReviews() ([]model.Review, error) {
	query := `
        SELECT id, submission_id, reviewer_id, status, comments, hod_action, hod_remarks, created_at, updated_at
        FROM reviews
    `
	var reviews []model.Review
	err := DB.Select(&reviews, query)
	if err != nil {
		fmt.Printf("Error fetching reviews: %v\n", err)
		return nil, err
	}
	return reviews, nil
}
