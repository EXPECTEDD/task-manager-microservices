package createdto

type CreateResponse struct {
	TaskId uint32 `json:"task_id" binding:"required"`
}
