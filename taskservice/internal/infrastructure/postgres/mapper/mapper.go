package posmapper

import (
	taskdomain "taskservice/internal/domain/task"
	posmodels "taskservice/internal/infrastructure/postgres/models"
)

func TaskDomainToModel(td *taskdomain.TaskDomain) *posmodels.TaskPosModel {
	return posmodels.NewTaskPosModel(
		td.Id,
		td.ProjectId,
		td.Description,
		td.Deadline,
	)
}

func TaskModelsToDomains(tm []*posmodels.TaskPosModel) []*taskdomain.TaskDomain {
	tasks := []*taskdomain.TaskDomain{}
	for _, t := range tm {
		tasks = append(tasks, taskdomain.RestoreTaskDomain(
			t.Id,
			t.ProjectId,
			t.Description,
			t.Deadline.Time,
		))
	}
	return tasks
}
