package datastore

import (
	"terminal_todo/types"
	"time"
)

type SimpleDS struct {
	ToDos []types.ToDo
}

func InitializeSimpleDS() *SimpleDS {
	todos := []types.ToDo{
		{Id: 1, Task: "Buy groceries", CreateTime: time.Now().Add(-48 * time.Hour)},
		{Id: 2, Task: "Finish report", CreateTime: time.Now().Add(-24 * time.Hour)},
		{Id: 3, Task: "Call Alice", CreateTime: time.Now().Add(-12 * time.Hour)},
		{Id: 4, Task: "Book flight tickets", CreateTime: time.Now().Add(-6 * time.Hour)},
		{Id: 5, Task: "Read Go documentation", CreateTime: time.Now()},
	}

	return &SimpleDS{
		ToDos: todos,
	}
}

func (sds *SimpleDS) AddTodo(task string) (types.ToDo, error) {
	todo := types.ToDo{
		Id: int64(len(sds.ToDos) + 1), Task: task, CreateTime: time.Now(),
	}

	sds.ToDos = append(sds.ToDos, todo)
	return todo, nil
}

func (sds *SimpleDS) DeleteTodo(id int64) error {
	for i, todo := range sds.ToDos {
		if todo.Id == id {
			sds.ToDos = append(sds.ToDos[:i], sds.ToDos[i+1:]...)
		}
	}

	return nil
}

func (sds *SimpleDS) GetTodos() ([]types.ToDo, error) {
	return sds.ToDos, nil
}
