package userdomain

type UserDomain struct {
	Id           uint32
	FirstName    string
	MiddleName   string
	LastName     string
	HashPassword string
	Email        string
}

func NewUserDomain(id uint32, firstname, middlename, lastname, hashPassword, email string) *UserDomain {
	return &UserDomain{
		Id:           id,
		FirstName:    firstname,
		MiddleName:   middlename,
		LastName:     lastname,
		HashPassword: hashPassword,
		Email:        email,
	}
}
