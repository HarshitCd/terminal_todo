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
	"time"

	"github.com/joho/godotenv"
)

func humanReadableTime(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return "time error"
	}
	now := time.Now().In(loc)
	duration := now.Sub(t)
	readableTime := ""

	switch {
	case duration < time.Minute:
		secs := int(duration.Seconds())
		if secs <= 1 {
			readableTime = "just now"
		} else {
			readableTime = fmt.Sprintf("%d seconds ago", secs)
		}
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes <= 1 {
			readableTime = "a minute ago"
		} else {
			readableTime = fmt.Sprintf("%d minutes ago", minutes)
		}
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours <= 1 {
			readableTime = "an hour ago"
		} else {
			readableTime = fmt.Sprintf("%d hours ago", hours)
		}
	case duration < 48*time.Hour:
		readableTime = "yesterday"
	case duration < 7*24*time.Hour:
		readableTime = fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	case duration == 7*24*time.Hour:
		readableTime = "a week ago"
	case duration < 30*24*time.Hour:
		readableTime = fmt.Sprintf("%d weeks ago", int(duration.Hours()/(24*7)))
	default:
		readableTime = t.Format("02 Jan 2006")
	}

	return readableTime
}

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
	headerFormate := "%-5s %-*s %-12s %-19s"
	formatString := "%-5d %-*s %-12s %-19s"
	taskWidth := 25

	for _, todo := range todos {
		if taskWidth-5 < len(todo.Task) {
			taskWidth = len(todo.Task) + 5
		}
	}

	clearScreen()
	header := fmt.Sprintf(headerFormate, "ID", taskWidth, "ToDo", "Status", "Created")

	fmt.Println(coloredString(header, 176, 50, 255))
	fmt.Println()

	for _, todo := range todos {
		fmt.Printf(formatString+"\n", todo.Id, taskWidth, todo.Task, todo.Status, humanReadableTime(todo.CreateTime))
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
		if mds != nil {
			mds.Close()
		}

		log.Error(err.Error())
		os.Exit(1)
	}
	defer mds.Close()

	ts := service.InitializeTodoService(mds)

	args := os.Args
	if len(args) < 2 {
		log.Error("Need to pass in an operaion (add/get/set/delete)")
		os.Exit(1)
	}

	op := args[1]
	op = strings.ToUpper(op)

	var status string
	if len(args) == 2 && op == "DELETE" {
		status = ""
	} else if len(args) == 3 {
		status = args[2]
		status = strings.ToLower(status)
	} else {
		log.Error("Need to pass in a status (add/get/set)")
		os.Exit(1)
	}
	fmt.Println(status, args)

	switch op {
	case "GET":
		display(ts.GetTodos(status))
	case "ADD":
		display(ts.GetTodos("all"))
		ts.AddTodo(status)
		display(ts.GetTodos("all"))

	case "SET":
		display(ts.GetTodos("all"))
		ts.SetTodos(status)
		display(ts.GetTodos("all"))
	case "DELETE":
		display(ts.GetTodos("all"))
		ts.DeleteTodo()
		display(ts.GetTodos("all"))
	default:
		log.Error("Entered invalid operation (add/get/delete)")
		os.Exit(1)
	}

}
