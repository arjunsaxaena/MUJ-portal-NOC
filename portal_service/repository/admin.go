package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(name, email, passwordHash, appPassword string) (string, error) {
	query := `
		INSERT INTO admin (name, email, password_hash, app_password)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := database.DB.QueryRow(query, name, email, passwordHash, appPassword).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetAdmins(filters model.GetAdminFilters) ([]model.Admin, error) {
	query := "SELECT * FROM admin WHERE 1=1"
	var args []interface{}
	paramIndex := 1

	if filters.ID != "" {
		query += " AND id = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.ID)
		paramIndex++
	}
	if filters.Email != "" {
		query += " AND email = $" + strconv.Itoa(paramIndex)
		args = append(args, filters.Email)
		paramIndex++
	}

	var admins []model.Admin
	err := database.DB.Select(&admins, query, args...)
	if err != nil {
		log.Printf("Error fetching admins with filters %v: %v", filters, err)
		return nil, err
	}

	return admins, nil
}

func ValidateAdminPassword(email, password string) (bool, error) {
	filters := model.GetAdminFilters{Email: email}
	admins, err := GetAdmins(filters)
	if err != nil {
		return false, err
	}

	if len(admins) != 1 {
		return false, errors.New("invalid credentials")
	}

	admin := admins[0]

	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("invalid credentials")
	}
	return true, nil
}

func UpdateAdmin(id string, name *string, passwordHash *string) error {
	query := "UPDATE admin SET "
	args := []interface{}{}
	argCount := 1

	if name != nil {
		query += "name = $" + fmt.Sprint(argCount) + ", "
		args = append(args, *name)
		argCount++
	}

	if passwordHash != nil {
		query += "password_hash = $" + fmt.Sprint(argCount) + ", "
		args = append(args, *passwordHash)
		argCount++
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = $" + fmt.Sprint(argCount)
	args = append(args, id) // Changed id to string

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Error updating admin with ID %s: %v", id, err)
		return err
	}

	return nil
}

func DeleteAdmin(id string) error {
	query := `
		DELETE FROM admin
		WHERE id = $1
	`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting admin with ID %s: %v", id, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no admin found with the given ID")
	}

	return nil
}
