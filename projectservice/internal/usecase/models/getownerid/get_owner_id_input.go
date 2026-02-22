package getowneridmodel

type GetOwnerIdInput struct {
	ProjectId uint32
}

func NewGetOwnerIdInput(projectId uint32) *GetOwnerIdInput {
	return &GetOwnerIdInput{
		ProjectId: projectId,
	}
}
