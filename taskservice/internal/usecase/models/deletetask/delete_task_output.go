package deletemodel

type DeleteTaskOutput struct {
	Deleted bool
}

func NewDeleteTaskOutput(deleted bool) *DeleteTaskOutput {
	return &DeleteTaskOutput{
		Deleted: deleted,
	}
}
