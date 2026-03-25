package getmodel

import taskdomain "taskservice/internal/domain/task"

type GetTaskOutput struct {
	Task *taskdomain.TaskDomain
}

func NewGetTaskOutput(task *taskdomain.TaskDomain) *GetTaskOutput {
	return &GetTaskOutput{
		Task: task,
	}
}
