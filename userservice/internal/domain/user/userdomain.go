package userdomain

type UserDomain struct {
	FirstName    string
	MiddleName   string
	LastName     string
	HashPassword string
	Email        string
}

func NewUserDomain(firstname, middlename, lastname, hashPassword, email string) *UserDomain {
	return &UserDomain{
		FirstName:    firstname,
		MiddleName:   middlename,
		LastName:     lastname,
		HashPassword: hashPassword,
		Email:        email,
	}
}
