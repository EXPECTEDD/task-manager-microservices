package getowneridmodel

type GetOwnerIdOutput struct {
	OwnerId uint32
}

func NewGetOwnerIdOutput(ownerId uint32) *GetOwnerIdOutput {
	return &GetOwnerIdOutput{
		OwnerId: ownerId,
	}
}
