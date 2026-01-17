package handlmapper

import (
	regdto "userservice/internal/transport/rest/handler/dto/registration"
	regmodel "userservice/internal/usecase/models/registration"
)

func RegRequestToInput(r *regdto.RegistrationRequest) *regmodel.RegInput {
	return regmodel.NewRegInput(
		r.FirstName,
		r.MiddleName,
		r.LastName,
		r.Password,
		r.Email,
	)
}

func RegOutputToResponse(ro *regmodel.RegOutput) *regdto.RegistrationResponse {
	return &regdto.RegistrationResponse{
		IsRegistered: ro.IsRegistered,
	}
}
