package auth

import (
	"database/sql"
	"errors"

	"github.com/Dryluigi/golang-todo-list/models"
	"github.com/Dryluigi/golang-todo-list/requests"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredential error = errors.New("invalid credential")

func Login(db *sql.DB, request *requests.AuthLogin) (*IssuedToken, error) {
	userModel := models.User{}

	row := db.QueryRow("SELECT * FROM tbl_users WHERE email = ?", request.Email)
	err := row.Scan(&userModel.ID, &userModel.Email, &userModel.Password, &userModel.Username)

	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredential
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(request.Password)); err != nil {
		return nil, ErrInvalidCredential
	}

	accessToken, err := IssueAccessToken(userModel.ID)

	if err != nil {
		return nil, err
	}

	return &IssuedToken{AccessToken: accessToken}, nil
}

func Register(db *sql.DB, request *requests.AuthRegister) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)

	if err != nil {
		return nil, err
	}

	result, err := db.Exec(
		"INSERT INTO tbl_users (email, password, username) VALUES (?, ?, ?)",
		request.Email,
		string(hashedPassword[:]),
		request.Username,
	)

	if err != nil {
		return nil, err
	}

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	userModel := models.User{
		ID:       uint(lastInsertId),
		Username: request.Username,
		Password: string(hashedPassword[:]),
		Email:    request.Email,
	}

	return &userModel, nil
}

func Logout(db *sql.DB) error {
	return nil
}
