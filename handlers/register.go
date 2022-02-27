package handlers

import (
	"github.com/Dryluigi/golang-todo-list/controllers/api"
	"github.com/Dryluigi/golang-todo-list/middlewares"
	"github.com/gorilla/mux"
)

func RegisterHandler(router *mux.Router) {
	router.HandleFunc("/", api.HelloWorld).Methods("GET")

	router.HandleFunc("/auth/register", api.AuthRegister).Methods("POST")
	router.HandleFunc("/auth/login", api.AuthLogin).Methods("POST")

	todoRouter := router.PathPrefix("/todos").Subrouter()
	todoRouter.Use(middlewares.RequireAuth)
	todoRouter.HandleFunc("", api.GetAllTodos).Methods("GET")
	todoRouter.HandleFunc("", api.CreateTodo).Methods("POST")
	todoRouter.HandleFunc("/toggleCompletion/{id}", api.ToggleTodoCompletion).Methods("PATCH")
	todoRouter.HandleFunc("/{id}", api.GetTodo).Methods("GET")
	todoRouter.HandleFunc("/{id}", api.DeleteTodo).Methods("DELETE")
	todoRouter.HandleFunc("/{id}", api.UpdateTodo).Methods("PUT")
}
