package grpchandlmapper

import (
	getowneridmodel "projectservice/internal/usecase/models/getownerid"
	projectservicev1 "projectservice/proto/projectservice"
)

func GetOwnerIdRequestToInput(req *projectservicev1.GetOwnerIdRequest) *getowneridmodel.GetOwnerIdInput {
	return getowneridmodel.NewGetOwnerIdInput(req.ProjectId)
}
