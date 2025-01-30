package util

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func ImportCSVToPostgres(csvFile string, db *sqlx.DB) error {
	file, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	insertSQL := `
	INSERT INTO students (registration_number, name, official_mail_id)
	VALUES ($1, $2, $3)
	ON CONFLICT (registration_number) DO NOTHING;
	`

	for i, row := range rows {
		if i == 0 {
			continue
		}

		registrationNumber := row[1]
		name := row[2]
		officialMailID := row[8]

		_, err := db.Exec(insertSQL, registrationNumber, name, officialMailID)
		if err != nil {
			return fmt.Errorf("error inserting data: %v", err)
		}
	}

	return nil
}
