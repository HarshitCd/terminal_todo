package datastore

import "terminal_todo/types"

type DataStore interface {
	AddTodo(task string, status string) (types.ToDo, error)
	DeleteTodo(id int64) error
	GetTodos(status string) ([]types.ToDo, error)
	SetTodos(status string, id int64) error
}
