package posmapper

import (
	userdomain "userservice/internal/domain/user"
	posmodels "userservice/internal/infrastructure/postgres/models"
)

func ModelToDomain(um *posmodels.UserPosModel) *userdomain.UserDomain {
	return userdomain.NewUserDomain(
		um.Id,
		um.FirstName,
		um.MiddleName.String,
		um.LastName,
		um.HashPassword,
		um.Email,
	)
}

func DomainToModel(ud *userdomain.UserDomain) *posmodels.UserPosModel {
	return posmodels.NewUserPosModel(
		0,
		ud.FirstName,
		ud.MiddleName,
		ud.LastName,
		ud.HashPassword,
		ud.Email,
	)
}
