package handlmapper

import (
	createdto "projectservice/internal/transport/rest/handler/dto/create"
	deletedto "projectservice/internal/transport/rest/handler/dto/delete"
	getalldto "projectservice/internal/transport/rest/handler/dto/getall"
	createmodel "projectservice/internal/usecase/models/createproject"
	deletemodel "projectservice/internal/usecase/models/deleteproject"
	getallmodel "projectservice/internal/usecase/models/getallprojects"
)

func CreateRequestToInput(cr *createdto.CreateRequest, userId uint32) *createmodel.CreateProjectInput {
	return createmodel.NewCreateProjectInput(userId, cr.Name)
}

func CreateOutputToResponse(co *createmodel.CreateProjectOutput) *createdto.CreateResponse {
	return &createdto.CreateResponse{
		ProjectId: co.ProjectId,
	}
}

func DeleteRequestToInput(dr *deletedto.DeleteRequest, userId uint32) *deletemodel.DeleteProjectInput {
	return deletemodel.NewDeleteProjectInput(userId, dr.ProjectId)
}

func DeleteOutputToResponse(do *deletemodel.DeleteProjectOutput) *deletedto.DeleteResponse {
	return &deletedto.DeleteResponse{
		IsDeleted: do.IsDeleted,
	}
}

func GetAllOutputToResponse(gao *getallmodel.GetAllProjectsOutput) *getalldto.GetAllResponse {
	return &getalldto.GetAllResponse{
		Projects: gao.Projects,
	}
}
