package getallmodel

import taskdomain "taskservice/internal/domain/task"

type GetALlTasksOutput struct {
	Tasks []*taskdomain.TaskDomain
}

func NewGetAllTasksOutput(tasks []*taskdomain.TaskDomain) *GetALlTasksOutput {
	return &GetALlTasksOutput{
		Tasks: tasks,
	}
}
