package logmodel

type LoginInput struct {
	Email    string
	Password string
}

func NewLoginInput(email, password string) *LoginInput {
	return &LoginInput{
		Email:    email,
		Password: password,
	}
}
