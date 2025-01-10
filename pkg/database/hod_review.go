package database

import (
	"MUJ_automated_mail_generation/pkg/model"
	"fmt"
	"time"
)

func CreateHodReview(submissionID, hodID int, action, remarks string) (int, error) {
	query := `
		INSERT INTO hod_reviews (submission_id, hod_id, action, remarks, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var reviewID int
	err := DB.QueryRow(query, submissionID, hodID, action, remarks, time.Now(), time.Now()).Scan(&reviewID)
	if err != nil {
		return 0, err
	}
	return reviewID, nil
}

func UpdateHodReview(reviewID int, action, remarks string) error {
	query := `
		UPDATE hod_reviews
		SET action = $1, remarks = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := DB.Exec(query, action, remarks, time.Now(), reviewID)
	return err
}

func GetHodReviewsBySubmission(submissionID int) ([]model.HodReview, error) {
	query := `
		SELECT id, submission_id, hod_id, action, remarks, created_at, updated_at
		FROM hod_reviews
		WHERE submission_id = $1
	`
	var reviews []model.HodReview
	err := DB.Select(&reviews, query, submissionID)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func GetHodReviewsByHod(hodID int) ([]model.HodReview, error) {
	query := `
		SELECT id, submission_id, hod_id, action, remarks, created_at, updated_at
		FROM hod_reviews
		WHERE hod_id = $1
	`
	var reviews []model.HodReview
	err := DB.Select(&reviews, query, hodID)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func GetAllHodReviews() ([]model.HodReview, error) {
	query := `
        SELECT id, submission_id, hod_id, action, remarks, created_at, updated_at
        FROM hod_reviews
    `
	var reviews []model.HodReview
	err := DB.Select(&reviews, query)
	if err != nil {
		fmt.Printf("Error fetching HoD reviews: %v\n", err)
		return nil, err
	}
	return reviews, nil
}
