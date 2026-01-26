package createmodel

type CreateProjectInput struct {
	Name string
}

func NewCreateProjectInput(name string) *CreateProjectInput {
	return &CreateProjectInput{
		Name: name,
	}
}