package api_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Dryluigi/golang-todo-list/database"
	"github.com/Dryluigi/golang-todo-list/models"
	"github.com/Dryluigi/golang-todo-list/requests"
	"github.com/Dryluigi/golang-todo-list/services/api"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var testDB *sql.DB
var testUserId uint = 999999999
var testTodoId uint = 999999999

func TestMain(m *testing.M) {
	godotenv.Load("../../.env")
	database.Connect()

	testDB = database.DB

	m.Run()
}

func TestSuccessGetAllTodo(t *testing.T) {
	todos, err := api.GetAllTodo(testDB)

	if err != nil {
		t.Errorf("there is error: " + err.Error())
	}

	todosType := reflect.TypeOf(todos)

	if todosType.Kind() != reflect.Slice {
		t.Errorf("return type should be slice, %s given", todosType.Kind().String())
	}
}

func TestSuccessGetTodoById(t *testing.T) {
	withDummyTodo(func() {
		todo, err := api.GetTodoById(testDB, testUserId, testTodoId)

		if err != nil {
			t.Errorf("there is error: %s", err.Error())
		}

		if todo.ID != testTodoId {
			t.Errorf("expected todo id %d, got %d instead", testTodoId, todo.ID)
		}
	})
}

func TestForbiddenGetTodoById(t *testing.T) {
	withDummyTodo(func() {
		_, err := api.GetTodoById(testDB, uint(41321312), testTodoId)

		testUnauthorizedError(err, t)
	})
}

func TestNotFoundGetTodoById(t *testing.T) {
	withDummyTodo(func() {
		_, err := api.GetTodoById(testDB, testUserId, uint(131312312))

		testNotFoundError(err, t)
	})
}

func TestSuccessCreateTodo(t *testing.T) {
	dummy := requests.TodoCreate{
		Todo:        "Testing Super",
		Description: "Super Testing",
		IsComplete:  false,
	}
	serialized, err := api.SaveTodo(testDB, uint(1), &dummy)

	if err != nil {
		t.Errorf("there is error: " + err.Error())
	}

	if serialized == nil {
		t.Errorf("the service is supposed to return serialized todo, nil given")
	}
}

func TestSuccessUpdateTodo(t *testing.T) {
	withDummyTodo(func() {
		updateTodo := requests.TodoUpdate{
			Todo:        "This is after update",
			Description: "This is after update hehe",
			IsComplete:  true,
		}

		todo, err := api.UpdateTodo(testDB, testUserId, testTodoId, updateTodo)

		if err != nil {
			t.Errorf("there is error: %s", err.Error())
		}

		if todo.ID != testTodoId {
			t.Errorf("expected todo id %d, %d given", testTodoId, todo.ID)
		}

		if todo.Todo != updateTodo.Todo {
			t.Errorf("expected todo content updated to %s, %s given", updateTodo.Todo, todo.Todo)
		}

		if todo.Description != updateTodo.Description {
			t.Errorf("expected todo description updated to %s, %s given", updateTodo.Description, todo.Description)
		}

		if todo.IsComplete != updateTodo.IsComplete {
			t.Errorf("expected todo completion updated to %t, %t given", updateTodo.IsComplete, todo.IsComplete)
		}
	})
}

func TestForbiddenUpdateTodo(t *testing.T) {
	withDummyTodo(func() {
		_, err := api.UpdateTodo(testDB, uint(1321231231), testTodoId, requests.TodoUpdate{})

		testUnauthorizedError(err, t)
	})
}

func TestNotFoundUpdateTodo(t *testing.T) {
	withDummyTodo(func() {
		_, err := api.UpdateTodo(testDB, testUserId, uint(123131221), requests.TodoUpdate{})

		testNotFoundError(err, t)
	})
}

func TestSuccessDeleteTodo(t *testing.T) {
	withDummyTodo(func() {
		err := api.DeleteTodo(testDB, testUserId, testTodoId)

		if err != nil {
			t.Errorf("expected error is nil, %s given", err.Error())
		}
	})
}

func TestForbiddenDeleteTodo(t *testing.T) {
	withDummyTodo(func() {
		err := api.DeleteTodo(testDB, uint(1321231231), testTodoId)

		testUnauthorizedError(err, t)
	})
}

func TestNotFoundDeleteTodo(t *testing.T) {
	withDummyTodo(func() {
		err := api.DeleteTodo(testDB, testUserId, uint(123131221))

		testNotFoundError(err, t)
	})
}

func TestSuccessToggleTodoCompletion(t *testing.T) {
	withDummyTodo(func() {
		todo, err := api.ToggleTodoCompletion(testDB, testUserId, testTodoId)

		if err != nil {
			t.Errorf("there is error: %s", err.Error())
		}

		if !todo.IsComplete {
			t.Errorf("expected todo completion to be %t, %t given", true, todo.IsComplete)
		}

		todo, err = api.ToggleTodoCompletion(testDB, testUserId, testTodoId)

		if err != nil {
			t.Errorf("there is error: %s", err.Error())
		}

		if todo.IsComplete {
			t.Errorf("expected todo completion to be %t, %t given", false, todo.IsComplete)
		}
	})
}

func TestForbiddenToggleCompletionTodo(t *testing.T) {
	withDummyTodo(func() {
		_, err := api.ToggleTodoCompletion(testDB, uint(123123123), testTodoId)

		testUnauthorizedError(err, t)
	})
}

func TestNotFoundToggleCompletionTodo(t *testing.T) {
	withDummyTodo(func() {
		_, err := api.ToggleTodoCompletion(testDB, testUserId, uint(123131221))

		testNotFoundError(err, t)
	})
}

func insertTestUser(db *sql.DB) (*models.User, error) {
	testPassword, err := bcrypt.GenerateFromPassword([]byte("123"), 10)

	if err != nil {
		return nil, err
	}

	userModel := models.User{
		ID:       testUserId,
		Email:    "test1312@gmail.com",
		Password: string(testPassword[:]),
		Username: "Yanto yanti",
	}
	_, err = db.Exec(
		"INSERT INTO tbl_users (id, email, password, username) VALUES (?, ?, ?, ?)",
		userModel.ID,
		userModel.Email,
		userModel.Password,
		userModel.Username,
	)

	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

func withDummyTodo(actualTest func()) {
	insertTestUser(testDB)
	insertTestTodo(testDB)
	defer cleanTestTodo(testDB)
	defer cleanTestUser(testDB)

	actualTest()
}

func cleanTestUser(db *sql.DB) {
	db.Exec("DELETE FROM tbl_users WHERE id = ?", testUserId)
}

func insertTestTodo(db *sql.DB) (*models.Todo, error) {
	todoModel := models.Todo{
		ID:          testTodoId,
		Todo:        "Testing super duper",
		Description: "Hehe",
		IsComplete:  0,
		UserId:      testUserId,
	}
	_, err := db.Exec(
		"INSERT INTO tbl_todos (id, todo, description, is_complete, user_id) VALUES (?, ?, ?, ?, ?)",
		todoModel.ID,
		todoModel.Todo,
		todoModel.Description,
		todoModel.IsComplete,
		todoModel.UserId,
	)

	if err != nil {
		return nil, err
	}

	return &todoModel, err
}

func cleanTestTodo(db *sql.DB) {
	db.Exec("DELETE FROM tbl_todos WHERE id = ?", testTodoId)
}

func testNotFoundError(err error, t *testing.T) {
	if err == nil {
		t.Errorf("expected error not nil, nil given")
	}

	if err != api.ErrEntityNotFound {
		t.Errorf("expected error to be %s, %s given", api.ErrEntityNotFound.Error(), err.Error())
	}
}

func testUnauthorizedError(err error, t *testing.T) {
	if err == nil {
		t.Errorf("expected error not nil, nil given")
	}

	if err != api.ErrUnauthorizedAccess {
		t.Errorf("expected error to be %s, %s given", api.ErrUnauthorizedAccess.Error(), err.Error())
	}
}
