package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"terminal_todo/datastore"
	"terminal_todo/service"
	"terminal_todo/types"

	"github.com/joho/godotenv"
)

func coloredString(s string, r, g, b int) string {
	coloredStringTemplate := "\033[38;2;%d;%d;%dm%s\033[0m"
	return fmt.Sprintf(coloredStringTemplate, r, g, b, s)
}

func clearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func display(todos []types.ToDo) {
	headerFormate := "%-5s %-35s %-19s"
	formatString := "%-5d %-35s %-19s"

	clearScreen()
	header := fmt.Sprintf(headerFormate, "ID", "ToDo", "Created")

	fmt.Println(coloredString(header, 176, 50, 255))
	fmt.Println()

	for i := len(todos) - 1; i >= 0; i-- {
		todo := todos[i]
		fmt.Printf(formatString+"\n", todo.Id, todo.Task, todo.CreateTime.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()
}

var envPath string

func main() {
	// sds := datastore.InitializeSimpleDS()

	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println("error while loading the .env file")
		os.Exit(1)
	}

	log, err := service.InitializeLogger()
	if err != nil {
		fmt.Println(fmt.Sprintf("error while setting up logger, %v", err))
		os.Exit(1)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbDatabase := os.Getenv("DB_DATABASE")

	mds, err := datastore.InitializeMySqlDS(dbUser, dbPassword, dbHost, dbDatabase)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer mds.Db.Close()
	ts := service.InitializeTodoService(mds)

	args := os.Args
	if len(args) < 2 {
		log.Error("Need to pass in an operaion (add/get/delete)")
		os.Exit(1)
	}

	op := args[1]
	op = strings.ToUpper(op)

	switch op {
	case "GET":
		display(ts.GetTodos())
	case "ADD":
		display(ts.GetTodos())
		ts.AddTodo()
		display(ts.GetTodos())
	case "DELETE":
		display(ts.GetTodos())
		ts.DeleteTodo()
		display(ts.GetTodos())
	default:
		log.Error("Entered invalid operation (add/get/delete)")
		os.Exit(1)
	}

}
