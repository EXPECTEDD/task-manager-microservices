package deletemodel

type DeleteTaskInput struct {
	TaskId    uint32
	ProjectId uint32
}

func NewDeleteTaskInput(taskId, projectId uint32) *DeleteTaskInput {
	return &DeleteTaskInput{
		TaskId:    taskId,
		ProjectId: projectId,
	}
}
