package api

import (
	"database/sql"
	"errors"

	"github.com/Dryluigi/golang-todo-list/models"
	"github.com/Dryluigi/golang-todo-list/requests"
)

type TodoSerializer struct {
	ID          uint   `json:"id"`
	Todo        string `json:"todo"`
	Description string `json:"description"`
	IsComplete  bool   `json:"is_complete"`
}

var ErrEntityNotFound error = errors.New("not found")
var ErrUnauthorizedAccess error = errors.New("unauthorized")

func serializeTodoModels(todos []models.Todo) []TodoSerializer {
	serialized := make([]TodoSerializer, 0)

	for _, todo := range todos {
		serialized = append(serialized, serializeTodoModel(&todo))
	}

	return serialized
}

func serializeTodoModel(todo *models.Todo) TodoSerializer {
	return TodoSerializer{
		ID:          todo.ID,
		Todo:        todo.Todo,
		Description: todo.Description,
		IsComplete:  todo.IsComplete == 1,
	}
}

func GetAllTodo(db *sql.DB) ([]TodoSerializer, error) {
	todos := make([]models.Todo, 0)

	rows, err := db.Query("SELECT * FROM tbl_todos")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		todo := models.Todo{}

		if err := rows.Scan(&todo.ID, &todo.Todo, &todo.Description, &todo.IsComplete, &todo.UserId); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	serializeds := serializeTodoModels(todos)

	return serializeds, nil
}

func SaveTodo(db *sql.DB, userId uint, request *requests.TodoCreate) (*TodoSerializer, error) {
	result, err := db.Exec(
		"INSERT INTO tbl_todos (todo, description, is_complete, user_id) VALUES (?, ?, ?, ?)",
		request.Todo,
		request.Description,
		func() int {
			if request.IsComplete {
				return 1
			} else {
				return 0
			}
		}(),
		userId,
	)

	if err != nil {
		return nil, err
	}

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	serialized := TodoSerializer{
		ID:          uint(lastInsertId),
		Todo:        request.Todo,
		Description: request.Description,
		IsComplete:  request.IsComplete,
	}

	return &serialized, nil
}

func GetTodoById(db *sql.DB, userId, todoId uint) (*TodoSerializer, error) {
	if !verifyTodoExists(db, todoId) {
		return nil, ErrEntityNotFound
	}

	todoModel := models.Todo{}

	row := db.QueryRow("SELECT * FROM tbl_todos WHERE id = ?", todoId)

	err := row.Scan(&todoModel.ID, &todoModel.Todo, &todoModel.Description, &todoModel.IsComplete, &todoModel.UserId)

	if err == sql.ErrNoRows {
		return nil, ErrEntityNotFound
	}

	if todoModel.UserId != userId {
		return nil, ErrUnauthorizedAccess
	}

	if err != nil {
		return nil, err
	}

	serialized := serializeTodoModel(&todoModel)

	return &serialized, nil
}

func UpdateTodo(db *sql.DB, userId, todoId uint, request requests.TodoUpdate) (*TodoSerializer, error) {
	if !verifyTodoExists(db, todoId) {
		return nil, ErrEntityNotFound
	}

	if !verifyTodoCreator(db, todoId, userId) {
		return nil, ErrUnauthorizedAccess
	}

	_, err := db.Exec(
		"UPDATE tbl_todos SET todo = ?, description = ?, is_complete = ? WHERE id = ?",
		request.Todo,
		request.Description,
		func() int {
			if request.IsComplete {
				return 1
			} else {
				return 0
			}
		}(),
		todoId,
	)

	if err != nil {
		return nil, err
	}

	serializer := TodoSerializer{
		ID:          todoId,
		Todo:        request.Todo,
		Description: request.Description,
		IsComplete:  request.IsComplete,
	}

	return &serializer, nil
}

func DeleteTodo(db *sql.DB, userId, todoId uint) error {
	if !verifyTodoExists(db, todoId) {
		return ErrEntityNotFound
	}

	if !verifyTodoCreator(db, todoId, userId) {
		return ErrUnauthorizedAccess
	}

	_, err := db.Exec("DELETE FROM tbl_todos WHERE id = ?", todoId)

	if err != nil {
		return err
	}

	return nil
}

func GetTodosByUserId(db *sql.DB, userId uint) ([]TodoSerializer, error) {
	todos := make([]models.Todo, 0)

	rows, err := db.Query("SELECT * FROM tbl_todos WHERE user_id = ?", userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		todo := models.Todo{}

		if err := rows.Scan(&todo.ID, &todo.Todo, &todo.Description, &todo.IsComplete, &todo.UserId); err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	serializeds := serializeTodoModels(todos)

	return serializeds, nil
}

func ToggleTodoCompletion(db *sql.DB, userId uint, todoId uint) (*TodoSerializer, error) {
	if !verifyTodoExists(db, todoId) {
		return nil, ErrEntityNotFound
	}

	if !verifyTodoCreator(db, todoId, userId) {
		return nil, ErrUnauthorizedAccess
	}

	_, err := db.Exec("UPDATE tbl_todos SET is_complete = (CASE WHEN is_complete = '1' THEN '0' ELSE '1' END) WHERE id = ?", todoId)

	if err != nil {
		return nil, err
	}

	todo, err := GetTodoById(db, userId, todoId)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func verifyTodoCreator(db *sql.DB, todoId, userId uint) bool {
	var tempId uint
	err := db.QueryRow("SELECT id FROM tbl_todos WHERE id = ? AND user_id = ?", todoId, userId).Scan(&tempId)

	return err == nil
}

func verifyTodoExists(db *sql.DB, todoId uint) bool {
	var tempId uint
	err := db.QueryRow("SELECT id FROM tbl_todos WHERE id = ?", todoId).Scan(&tempId)

	return err != sql.ErrNoRows
}
