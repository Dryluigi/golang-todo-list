package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Dryluigi/golang-todo-list/apiHelper/response"
	"github.com/Dryluigi/golang-todo-list/database"
	"github.com/Dryluigi/golang-todo-list/requests"
	services "github.com/Dryluigi/golang-todo-list/services/api"
	"github.com/Dryluigi/golang-todo-list/services/auth"
	"github.com/gorilla/mux"
)

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	todos, err := services.GetTodosByUserId(database.DB, userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	response.BuildSuccessResponse(w, http.StatusOK, "Todo Fetched", true, todos)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		response.BuildErrorResponse(w, http.StatusBadRequest, "Invalid id. Id must be an integer")
		return
	}

	todo, err := services.GetTodoById(database.DB, userId, uint(id))

	if err == services.ErrEntityNotFound {
		response.BuildErrorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	if err == services.ErrUnauthorizedAccess {
		response.BuildErrorResponse(w, http.StatusForbidden, "Forbidden")
		return
	}

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	response.BuildSuccessResponse(w, http.StatusOK, "Todo Fetched", true, todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	todo := requests.TodoCreate{}
	err = json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	serialized, err := services.SaveTodo(database.DB, userId, &todo)

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	response.BuildSuccessResponse(w, http.StatusCreated, "Todo Created", true, serialized)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		response.BuildErrorResponse(w, http.StatusBadRequest, "Invalid id. Id must be an integer")
		return
	}

	todoUpdate := requests.TodoUpdate{}

	err = json.NewDecoder(r.Body).Decode(&todoUpdate)

	if err != nil {
		response.BuildErrorResponse(w, 400, "Invalid data")
		return
	}

	serializer, err := services.UpdateTodo(database.DB, userId, uint(id), todoUpdate)

	if err == services.ErrEntityNotFound {
		response.BuildErrorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	if err == services.ErrUnauthorizedAccess {
		response.BuildErrorResponse(w, http.StatusForbidden, "Forbidden")
		return
	}

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	response.BuildSuccessResponse(w, 200, "Todo Updated", true, serializer)

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		response.BuildErrorResponse(w, http.StatusBadRequest, "Invalid id. Id must be an integer")
		return
	}

	err = services.DeleteTodo(database.DB, userId, uint(id))

	if err == services.ErrEntityNotFound {
		response.BuildErrorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	if err == services.ErrUnauthorizedAccess {
		response.BuildErrorResponse(w, http.StatusForbidden, "Forbidden")
		return
	}

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	response.BuildSuccessResponse(w, http.StatusOK, "Delete Success", true, nil)
}

func ToggleTodoCompletion(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)

	if err != nil {
		response.BuildServerErrorResponse(w, err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		response.BuildErrorResponse(w, http.StatusBadRequest, "Invalid id. Id must be an integer")
		return
	}

	todo, err := services.ToggleTodoCompletion(database.DB, userId, uint(id))

	if err != nil {
		if err == services.ErrEntityNotFound {
			response.BuildErrorResponse(w, http.StatusNotFound, "Not Found")
			return
		} else if err == services.ErrUnauthorizedAccess {
			response.BuildErrorResponse(w, http.StatusForbidden, "Forbidden")
			return
		} else {
			response.BuildServerErrorResponse(w, err)
			return
		}
	}

	response.BuildSuccessResponse(w, http.StatusOK, "Toco completion changed", true, todo)
}

func getUserIdFromRequest(r *http.Request) (uint, error) {
	userIdContext := r.Context().Value(auth.UserIdContextKey)
	userId, ok := userIdContext.(uint)

	if !ok {
		return 0, errors.New("invalid user id")
	}

	return userId, nil
}
