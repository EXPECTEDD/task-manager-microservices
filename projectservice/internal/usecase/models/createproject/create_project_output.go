package createmodel

type CreateProjectOutput struct {
	IsCreated bool
}

func NewCreateProjectOutput(isCreated bool) *CreateProjectOutput {
	return &CreateProjectOutput{
		IsCreated: isCreated,
	}
}