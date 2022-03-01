package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func migrateDB(db *sql.DB) {
	log.Println("Migration is running")
	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		log.Fatal("there is error: " + err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migration",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatal("there is error: " + err.Error())
	}

	m.Up()
}
