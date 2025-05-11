package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"terminal_todo/datastore"
	"terminal_todo/types"
)

type TodoService struct {
	Ds datastore.DataStore
}

func InitializeTodoService(ds datastore.DataStore) *TodoService {
	return &TodoService{
		Ds: ds,
	}
}

func (ts *TodoService) AddTodo() types.ToDo {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Todo: ")
	input, _ := reader.ReadString('\n')
	task := strings.TrimSpace(input)

	todo, err := ts.Ds.AddTodo(task)
	if err != nil {
		fmt.Println("error has taken place in add service", err)
	}

	return todo
}

func (ts *TodoService) DeleteTodo() {
	var id int64
	fmt.Print("Enter Todo Id: ")
	fmt.Scan(&id)

	err := ts.Ds.DeleteTodo(id)
	if err != nil {
		fmt.Println("error has taken place in delete service", err)
	}
}

func (ts *TodoService) GetTodos() []types.ToDo {
	todos, err := ts.Ds.GetTodos()
	if err != nil {
		fmt.Println("error has taken place in get service", err)
	}

	return todos
}
