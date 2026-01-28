package createdto

type CreateRequest struct {
	Name string `json:"name" binding:"required"`
}
