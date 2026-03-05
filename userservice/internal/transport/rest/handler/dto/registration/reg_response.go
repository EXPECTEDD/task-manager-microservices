package regdto

type RegistrationResponse struct {
	UserId uint32 `json:"user_id" binding:"required"`
}
