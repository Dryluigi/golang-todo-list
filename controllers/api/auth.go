package api

import (
	"encoding/json"
	"net/http"

	"github.com/Dryluigi/golang-todo-list/apiHelper/response"
	"github.com/Dryluigi/golang-todo-list/database"
	"github.com/Dryluigi/golang-todo-list/requests"
	"github.com/Dryluigi/golang-todo-list/services/api"
	"github.com/Dryluigi/golang-todo-list/services/auth"
)

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	request := requests.AuthLogin{}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		response.BuildErrorResponse(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	token, err := auth.Login(database.DB, &request)

	if err != nil && err == auth.ErrInvalidCredential {
		response.BuildSuccessResponse(w, http.StatusOK, "Invalid Credential", false, nil)
	}

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	response.BuildSuccessResponse(w, http.StatusOK, "Login Success", true, token)
}

func AuthRegister(w http.ResponseWriter, r *http.Request) {
	request := requests.AuthRegister{}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		response.BuildErrorResponse(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	if api.VerifyUserEmailExist(database.DB, request.Email) {
		response.BuildSuccessResponse(w, http.StatusOK, "Email Already Exist", false, nil)
		return
	}

	userModel, err := auth.Register(database.DB, &request)

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	response.BuildSuccessResponse(w, http.StatusOK, "Account Created", true, userModel)
}
