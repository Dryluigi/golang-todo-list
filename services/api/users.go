package api

import "database/sql"

func VerifyUserEmailExist(db *sql.DB, email string) bool {
	var temp uint
	err := db.QueryRow("SELECT id FROM tbl_users WHERE email = ?", email).Scan(&temp)

	return err == nil
}
