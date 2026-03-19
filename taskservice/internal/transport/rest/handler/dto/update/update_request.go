package updatedto

import "time"

type UpdateRequest struct {
	NewDescription *string    `json:"new_description"`
	NewDeadline    *time.Time `json:"new_deadline"`
}
