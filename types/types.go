package types

import "time"

type ToDo struct {
	Id         int64
	Task       string
	Status     string
	CreateTime time.Time
}
