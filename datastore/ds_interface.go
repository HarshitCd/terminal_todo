package datastore

import "terminal_todo/types"

type DataStore interface {
	AddTodo(task string) (types.ToDo, error)
	DeleteTodo(id int64) error
	GetTodos() ([]types.ToDo, error)
}
