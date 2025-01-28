package repository

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/model"
	"errors"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(name, email, passwordHash, appPassword string) (int, error) {
	query := `
		INSERT INTO admin (name, email, password_hash, app_password)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err := database.DB.QueryRow(query, name, email, passwordHash, appPassword).Scan(&id)
	if err != nil {
		return 0, err
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

func UpdateAdmin(id int, name, email, passwordHash, appPassword string) error {
	query := `
		UPDATE admin
		SET name = COALESCE($1, name),
		    email = COALESCE($2, email),
		    password_hash = COALESCE($3, password_hash),
		    app_password = COALESCE($4, app_password)
		WHERE id = $5
	`
	_, err := database.DB.Exec(query, name, email, passwordHash, appPassword, id)
	if err != nil {
		log.Printf("Error updating admin with ID %d: %v", id, err)
		return err
	}

	return nil
}

func DeleteAdmin(id int) error {
	query := `
		DELETE FROM admin
		WHERE id = $1
	`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting admin with ID %d: %v", id, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no admin found with the given ID")
	}

	return nil
}
