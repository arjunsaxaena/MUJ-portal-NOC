package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

type ConfigInterface interface {
	GetDatabaseURL() string
}

func Connect(cfg ConfigInterface) {
	dsn := cfg.GetDatabaseURL()

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Database connection established")
}

func Close() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Fatalf("Error closing the database connection: %v", err)
		} else {
			fmt.Println("Database connection closed")
		}
	}
}
