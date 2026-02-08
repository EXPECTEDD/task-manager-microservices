package deletemodel

type DeleteProjectInput struct {
	ProjectId uint32
}

func NewDeleteProjectInput(projectId uint32) *DeleteProjectInput {
	return &DeleteProjectInput{
		ProjectId: projectId,
	}
}
