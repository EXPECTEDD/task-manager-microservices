package updatemodel

import "time"

type UpdateTaskInput struct {
	TaskId         uint32
	ProjectId      uint32
	NewDescription *string
	NewDeadline    *time.Time
}

func NewUpdateTaskInput(taskId uint32, projectId uint32, newDescription *string, newDeadline *time.Time) *UpdateTaskInput {
	return &UpdateTaskInput{
		TaskId:         taskId,
		ProjectId:      projectId,
		NewDescription: newDescription,
		NewDeadline:    newDeadline,
	}
}
