package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"errors"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func CreateHoD(name, email, passwordHash, appPassword, department string) (int, error) {
	query := `
		INSERT INTO hod (name, email, password_hash, app_password, department)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int
	err := database.DB.QueryRow(query, name, email, passwordHash, appPassword, department).Scan(&id)
	if err != nil {
		log.Printf("Error creating HoD: %v", err)
		return 0, err
	}
	return id, nil
}

func GetHoDs(filters model.GetHoDFilters) ([]model.HoD, error) {
	query := "SELECT * FROM hod WHERE 1=1"
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

	var hods []model.HoD
	err := database.DB.Select(&hods, query, args...)
	if err != nil {
		log.Printf("Error fetching HoDs with filters %v: %v", filters, err)
		return nil, err
	}

	return hods, nil
}

func ValidateHoDPassword(email, password string) (bool, error) {
	filters := model.GetHoDFilters{Email: email}
	hods, err := GetHoDs(filters)
	if err != nil {
		return false, err
	}

	if len(hods) != 1 {
		return false, errors.New("invalid credentials")
	}

	hod := hods[0]

	err = bcrypt.CompareHashAndPassword([]byte(hod.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("invalid credentials")
	}
	return true, nil
}

func UpdateHoD(id int, name, email, passwordHash, appPassword, department string) error {
	query := `
		UPDATE hod
		SET name = COALESCE($1, name),
		    email = COALESCE($2, email),
		    password_hash = COALESCE($3, password_hash),
		    app_password = COALESCE($4, app_password),
		    department = COALESCE($5, department)
		WHERE id = $6
	`
	_, err := database.DB.Exec(query, name, email, passwordHash, appPassword, department, id)
	if err != nil {
		log.Printf("Error updating HoD with ID %d: %v", id, err)
		return err
	}

	return nil
}

func DeleteHoD(id int) error {
	query := `
		DELETE FROM hod
		WHERE id = $1
	`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting HoD with ID %d: %v", id, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no HoD found with the given ID")
	}

	return nil
}
