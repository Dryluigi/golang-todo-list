package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func SetupTestDB() *sql.DB {
	godotenv.Load(".env.test")

	dbPort := os.Getenv("DB_PORT")
	portInt, err := strconv.Atoi(dbPort)

	if err != nil {
		log.Fatal("invalid port")
	}

	connectionStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		portInt,
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", connectionStr)

	if err != nil {
		log.Fatal("could not open test DB")
	}

	return db
}
