package logmodel

type LoginOutput struct {
	SessionId string
}

func NewLoginOutput(sessionId string) *LoginOutput {
	return &LoginOutput{
		SessionId: sessionId,
	}
}
