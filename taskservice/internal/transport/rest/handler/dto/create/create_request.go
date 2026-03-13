package createdto

import "time"

type CreateRequest struct {
	Description string    `json:"description" binding:"required"`
	Deadline    time.Time `json:"deadline"`
}
