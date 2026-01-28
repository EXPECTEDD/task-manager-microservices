package createdto

type CreateResponse struct {
	IsCreated bool `json:"is_created" binding:"required"`
}
