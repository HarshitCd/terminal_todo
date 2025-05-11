package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"terminal_todo/datastore"
	"terminal_todo/types"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type TodoService struct {
	Ds datastore.DataStore
}

func InitializeTodoService(ds datastore.DataStore) *TodoService {
	l, err := InitializeLogger()
	if err != nil {
		fmt.Println("error while setting up logger")
		os.Exit(1)
	}
	log = l

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
		log.Error("error has taken place in add service ", err)
		os.Exit(1)
	}

	return todo
}

func (ts *TodoService) DeleteTodo() {
	var id int64
	fmt.Print("Enter Todo Id: ")
	fmt.Scan(&id)

	err := ts.Ds.DeleteTodo(id)
	if err != nil {
		log.Error("error has taken place in delete service ", err)
		os.Exit(1)
	}
}

func (ts *TodoService) GetTodos() []types.ToDo {
	todos, err := ts.Ds.GetTodos()
	if err != nil {
		log.Error("error has taken place in get service ", err)
		os.Exit(1)
	}

	return todos
}
