package datastore

import (
	"database/sql"
	"fmt"
	"terminal_todo/types"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// CREATE TABLE IF NOT EXISTS todos (id INT AUTO_INCREMENT PRIMARY KEY, task VARCHAR(255) NOT NULL, create_time DATETIME DEFAULT CURRENT_TIMESTAMP );
type MySqlDS struct {
	Dsn string
	Db  *sql.DB
}

func InitializeMySqlDS(user string, password string, host string, database string) (*MySqlDS, error) {
	dsnTemplate := "%s:%s@tcp(%s)/%s?parseTime=true"
	dsn := fmt.Sprintf(dsnTemplate, user, password, host, database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	mds := &MySqlDS{
		Dsn: dsn,
		Db:  db,
	}

	err = db.Ping()
	if err != nil {
		return mds, err
	}

	return mds, nil
}

func (mds *MySqlDS) Close() error {
	return mds.Db.Close()
}

func (mds *MySqlDS) AddTodo(task string) (types.ToDo, error) {
	query := "INSERT INTO todos (task) VALUES (?)"
	stmt, err := mds.Db.Prepare(query)
	if err != nil {
		return types.ToDo{}, err
	}

	res, err := stmt.Exec(task)
	if err != nil {
		return types.ToDo{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return types.ToDo{}, err
	}

	return types.ToDo{
		Id:         id,
		Task:       task,
		CreateTime: time.Now(),
	}, nil
}

func (mds *MySqlDS) DeleteTodo(id int64) error {
	query := "DELETE FROM todos where id=?"
	stmt, err := mds.Db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (mds *MySqlDS) GetTodos() ([]types.ToDo, error) {
	todos := []types.ToDo{}
	query := "SELECT id, task, create_time FROM todos"
	rows, err := mds.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo types.ToDo
		err := rows.Scan(&todo.Id, &todo.Task, &todo.CreateTime)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}
