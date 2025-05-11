package types

import "time"

type ToDo struct {
	Id         int64
	Task       string
	CreateTime time.Time
}
