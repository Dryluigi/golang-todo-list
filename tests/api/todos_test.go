package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dryluigi/golang-todo-list/routers"
	"github.com/Dryluigi/golang-todo-list/tests/api/helpers"
)

func TestMain(m *testing.M) {
	m.Run()
}
func TestGetTodos(t *testing.T) {
	router := routers.SetupRouter()
	req, err := http.NewRequest("GET", "/todos", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	helpers.AssertStatusCode(t, http.StatusOK, rr.Code)
}

func GetTodos(t *testing.T) {
}
