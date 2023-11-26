package handlers

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func InitializeDB() error {

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file: %s", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbSource := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)

	var openErr error
	db, openErr = sql.Open("mysql", dbSource)
	if openErr != nil {
		return openErr
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return pingErr
	}

	return nil
}
