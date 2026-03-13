package handlmapper

import (
	createdto "taskservice/internal/transport/rest/handler/dto/create"
	createmodel "taskservice/internal/usecase/models/createtask"
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
