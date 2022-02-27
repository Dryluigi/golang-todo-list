package routers

import (
	"github.com/Dryluigi/golang-todo-list/handlers"
	"github.com/Dryluigi/golang-todo-list/middlewares"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.JsonResponse)
	router.Use(middlewares.LogRequest)

	handlers.RegisterHandler(router)

	router.Use(middlewares.EnableCors)
	router.Use(middlewares.LogResponse)

	return router
}
