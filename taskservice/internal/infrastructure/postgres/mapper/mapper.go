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

func TaskModelToDomain(tm *posmodels.TaskPosModel) *taskdomain.TaskDomain {
	return &taskdomain.TaskDomain{
		Id:          tm.Id,
		ProjectId:   tm.ProjectId,
		Description: tm.Description,
		Deadline:    tm.Deadline.Time,
	}
}

func TaskModelsToDomains(tm []*posmodels.TaskPosModel) []*taskdomain.TaskDomain {
	tasks := []*taskdomain.TaskDomain{}
	for _, t := range tm {
		tasks = append(tasks, TaskModelToDomain(t))
	}
	return tasks
}
