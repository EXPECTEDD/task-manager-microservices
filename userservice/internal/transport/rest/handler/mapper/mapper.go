package handlmapper

import (
	logindto "userservice/internal/transport/rest/handler/dto/login"
	regdto "userservice/internal/transport/rest/handler/dto/registration"
	logmodel "userservice/internal/usecase/models/login"
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
		UserId: ro.UserId,
	}
}

func LogRequestToInput(l *logindto.LoginRequest) *logmodel.LoginInput {
	return logmodel.NewLoginInput(
		l.Email,
		l.Password,
	)
}

func LogOutputToResponse(lo *logmodel.LoginOutput) *logindto.LoginResponse {
	return &logindto.LoginResponse{
		FirstName:  lo.FirstName,
		MiddleName: lo.MiddleName,
		LastName:   lo.LastName,
	}
}
