package posmapper

import (
	userdomain "userservice/internal/domain/user"
	posmodels "userservice/internal/infrastructure/postgres/models"
)

func ModelToDomain(um *posmodels.UserPosModel) *userdomain.UserDomain {
	return userdomain.NewUserDomain(
		um.FirstName,
		um.MiddleName.String,
		um.LastName,
		um.HashPassword,
		um.Email,
	)
}
