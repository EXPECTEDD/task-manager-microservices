package handlmapper

import (
	createdto "taskservice/internal/transport/rest/handler/dto/create"
	deletedto "taskservice/internal/transport/rest/handler/dto/delete"
	getalldto "taskservice/internal/transport/rest/handler/dto/getall"
	updatedto "taskservice/internal/transport/rest/handler/dto/update"
	createmodel "taskservice/internal/usecase/models/createtask"
	deletemodel "taskservice/internal/usecase/models/deletetask"
	getallmodel "taskservice/internal/usecase/models/getalltasks"
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

func UpdateRequestToInput(req *updatedto.UpdateRequest, taskId uint32, projectId uint32) *updatemodel.UpdateTaskInput {
	return updatemodel.NewUpdateTaskInput(
		taskId,
		projectId,
		req.NewDescription,
		req.NewDeadline,
	)
}

func UpdateOutputToResponse(out *updatemodel.UpdateTaskOutput) *updatedto.UpdateResponse {
	return &updatedto.UpdateResponse{
		Updated: out.Updated,
	}
}

func DeleteOutputToResponse(out *deletemodel.DeleteTaskOutput) *deletedto.DeleteResponse {
	return &deletedto.DeleteResponse{
		Deleted: out.Deleted,
	}
}

func GetAllOutputToResponse(out *getallmodel.GetAllTasksOutput) *getalldto.GetAllResponse {
	return &getalldto.GetAllResponse{
		Tasks: out.Tasks,
	}
}
