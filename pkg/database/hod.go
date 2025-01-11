package database

import (
	"MUJ_automated_mail_generation/pkg/model"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateHoD(name, email, passwordHash, department string) (int, error) {
	query := `
		INSERT INTO hod (name, email, password_hash, department)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err := DB.QueryRow(query, name, email, passwordHash, department).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetHoDByEmail(email string) (*model.HoD, error) {
	query := `
		SELECT id, name, email, password_hash, department, created_at
		FROM hod
		WHERE email = $1
	`
	var hod model.HoD
	err := DB.Get(&hod, query, email)
	if err != nil {
		fmt.Printf("Error fetching HoD: %v\n", err)
		return nil, err
	}
	return &hod, nil
}

func GetHoDsByDepartment(department string) ([]model.HoD, error) {
	query := `
		SELECT id, name, email, password_hash, department, created_at
		FROM hod
		WHERE department = $1
	`
	var hods []model.HoD
	err := DB.Select(&hods, query, department)
	if err != nil {
		fmt.Printf("Error fetching HoDs: %v\n", err)
		return nil, err
	}
	return hods, nil
}

func ValidateHoDPassword(email, password string) (bool, error) {
	hod, err := GetHoDByEmail(email)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hod.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("invalid credentials")
	}
	return true, nil
}

func GetAllHoDs() ([]model.HoD, error) {
	query := `
		SELECT id, name, email, password_hash, department, created_at
		FROM hod
	`
	var hods []model.HoD
	err := DB.Select(&hods, query)
	if err != nil {
		fmt.Printf("Error fetching all HoDs: %v\n", err)
		return nil, err
	}
	return hods, nil
}
