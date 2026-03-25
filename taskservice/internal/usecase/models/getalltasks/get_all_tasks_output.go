package getallmodel

import taskdomain "taskservice/internal/domain/task"

type GetAllTasksOutput struct {
	Tasks []*taskdomain.TaskDomain
}

func NewGetAllTasksOutput(tasks []*taskdomain.TaskDomain) *GetAllTasksOutput {
	return &GetAllTasksOutput{
		Tasks: tasks,
	}
}
