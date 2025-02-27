package contract

import (
	"time"
)

type Task struct {
	Title       string
	Description int
	Status      string
}

type GetTask struct {
	Id          int
	Title       string
	Description int
	Status      string
	Created_at  time.Time
	Updated_at  time.Time
}

type PutDelTask struct {
	Id          int
	Title       string
	Description int
	Status      string
}
