package regmodel

type RegOutput struct {
	UserId uint32
}

func NewRegOutput(userId uint32) *RegOutput {
	return &RegOutput{
		UserId: userId,
	}
}
