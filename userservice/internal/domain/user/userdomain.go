package userdomain

type UserDomain struct {
	FirstName string
	MiddleName string
	LastName string
	Password string
	Email string
}

func NewUserDomain(firstname, middlename, lastname, password, email string) *UserDomain {
	return &UserDomain{
		FirstName: firstname,
		MiddleName: middlename,
		LastName: lastname,
		Password: password,
		Email: email,
	}
}
