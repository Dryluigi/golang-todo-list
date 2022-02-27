package main

import (
	"log"
	"net/http"

	"github.com/Dryluigi/golang-todo-list/database"
	"github.com/Dryluigi/golang-todo-list/routers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	app := routers.SetupRouter()

	database.Connect()

	log.Fatal(http.ListenAndServe(":8080", app))
}
