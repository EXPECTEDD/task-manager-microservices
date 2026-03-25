package getmodel

type GetTaskInput struct {
	TaskId    uint32
	ProjectId uint32
}

func NewGetTaskInput(taskId uint32, projectId uint32) *GetTaskInput {
	return &GetTaskInput{
		TaskId:    taskId,
		ProjectId: projectId,
	}
}
