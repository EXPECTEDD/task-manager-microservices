package handlmapper

import (
	regdto "userservice/internal/transport/rest/handler/dto/registration"
	regmodel "userservice/internal/usecase/models/registration"
)

func RegRequestToInput(r *regdto.RegistrationRequest) (*regmodel.RegInput, error) {
	ri, err := regmodel.NewRegInput(
		r.FirstName,
		r.MiddleName,
		r.LastName,
		r.Password,
		r.Email,
	)
	if err != nil {
		return nil, err
	}
	return ri, nil
}

func RegOutputToResponse(ro *regmodel.RegOutput) *regdto.RegistrationResponse {
	return &regdto.RegistrationResponse{
		IsRegistered: ro.IsRegistered,
	}
}
