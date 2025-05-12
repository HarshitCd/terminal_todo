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
	dsnTemplate := "%s:%s@tcp(%s)/%s?parseTime=true&loc=Asia%%2FKolkata"
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

func (mds *MySqlDS) AddTodo(task string, status string) (types.ToDo, error) {
	query := "INSERT INTO todos (task, status) VALUES (?, ?)"
	stmt, err := mds.Db.Prepare(query)
	if err != nil {
		return types.ToDo{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(task, status)
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
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (mds *MySqlDS) SetTodos(status string, id int64) error {
	query := "UPDATE todos SET status=? WHERE id=?"
	stmt, err := mds.Db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(status, id)
	if err != nil {
		return err
	}

	return nil
}

func (mds *MySqlDS) GetTodos(status string) ([]types.ToDo, error) {
	todos := []types.ToDo{}
	var rows *sql.Rows
	if status == "all" {
		query := "SELECT id, task, status, create_time FROM todos ORDER BY create_time DESC"
		rs, err := mds.Db.Query(query)
		if err != nil {
			return todos, err
		}

		rows = rs
	} else {
		query := "SELECT id, task, status, create_time FROM todos WHERE status=? ORDER BY create_time DESC"
		stmt, err := mds.Db.Prepare(query)
		if err != nil {
			return todos, nil
		}
		defer stmt.Close()

		rs, err := stmt.Query(status)
		if err != nil {
			return todos, err
		}

		rows = rs
	}
	defer rows.Close()

	for rows.Next() {
		var todo types.ToDo
		err := rows.Scan(&todo.Id, &todo.Task, &todo.Status, &todo.CreateTime)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}
