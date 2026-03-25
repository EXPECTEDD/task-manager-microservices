package getdto

import taskdomain "taskservice/internal/domain/task"

type GetResponse struct {
	Task *taskdomain.TaskDomain `json:"task"`
}
