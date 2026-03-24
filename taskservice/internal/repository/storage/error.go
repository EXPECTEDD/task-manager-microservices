package storage

import "errors"

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrTasksNotFound = errors.New("tasks not found")
)
