package logmodel

type LoginInput struct {
	Email        string
	HashPassword string
}

func NewLoginInput(email, hashPassword string) *LoginInput {
	return &LoginInput{
		Email:        email,
		HashPassword: hashPassword,
	}
}
