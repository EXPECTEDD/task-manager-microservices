package getallmodel

type GetAllTasksInput struct {
	ProjectId uint32
}

func NewGetAllTasksInput(projectId uint32) *GetAllTasksInput {
	return &GetAllTasksInput{
		ProjectId: projectId,
	}
}
