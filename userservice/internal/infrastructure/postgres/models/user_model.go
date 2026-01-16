package posmodels

import "database/sql"

type UserPosModel struct {
	Id           uint32         `db:"id"`
	FirstName    string         `db:"first_name"`
	MiddleName   sql.NullString `db:"middle_name"`
	LastName     string         `db:"last_name"`
	HashPassword string         `db:"hash_password"`
	Email        string         `db:"email"`
}

func NewUserPosModel(id uint32, firstName, middleName, lastName, hashPassword, email string) *UserPosModel {
	midName := sql.NullString{
		String: middleName,
		Valid:  middleName != "",
	}

	return &UserPosModel{
		Id:           id,
		FirstName:    firstName,
		MiddleName:   midName,
		LastName:     lastName,
		HashPassword: hashPassword,
		Email:        email,
	}
}
