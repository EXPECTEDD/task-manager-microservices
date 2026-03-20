package handlmapper

import (
	createdto "taskservice/internal/transport/rest/handler/dto/create"
	updatedto "taskservice/internal/transport/rest/handler/dto/update"
	createmodel "taskservice/internal/usecase/models/createtask"
	updatemodel "taskservice/internal/usecase/models/updatetask"
)

func CreateRequestToInput(req *createdto.CreateRequest, projectId uint32) *createmodel.CreateTaskInput {
	return createmodel.NewCreateInput(
		projectId,
		req.Description,
		req.Deadline,
	)
}

func CreateOutputToResponse(out *createmodel.CreateTaskOutput) *createdto.CreateResponse {
	return &createdto.CreateResponse{
		TaskId: out.TaskId,
	}
}

func UpdateRequestToInput(req *updatedto.UpdateRequest, taskId uint32) *updatemodel.UpdateTaskInput {
	return updatemodel.NewUpdateTaskInput(
		taskId,
		req.NewDescription,
		req.NewDeadline,
	)
}

func UpdateOutputToResponse(out *updatemodel.UpdateTaskOutput) *updatedto.UpdateResponse {
	return &updatedto.UpdateResponse{
		Updated: out.Updated,
	}
}
