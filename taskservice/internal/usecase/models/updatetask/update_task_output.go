package updatemodel

type UpdateTaskOutput struct {
	Updated bool
}

func NewUpdateTaskOutput(updated bool) *UpdateTaskOutput {
	return &UpdateTaskOutput{
		Updated: updated,
	}
}
