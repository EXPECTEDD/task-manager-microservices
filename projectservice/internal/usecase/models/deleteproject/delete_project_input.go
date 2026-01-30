package deletemodel

type DeleteProjectInput struct {
	OwnerId uint32
	Name    string
}

func NewDeleteProjectInput(ownerId uint32, name string) *DeleteProjectInput {
	return &DeleteProjectInput{
		OwnerId: ownerId,
		Name:    name,
	}
}
