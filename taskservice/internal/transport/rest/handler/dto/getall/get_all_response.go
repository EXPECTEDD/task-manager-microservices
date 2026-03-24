package getalldto

import taskdomain "taskservice/internal/domain/task"

type GetAllResponse struct {
	Tasks []*taskdomain.TaskDomain `json:"tasks"`
}
