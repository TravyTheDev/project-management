package types

type ForgotPasswordReq struct {
	Email string `json:"email" validate:"required,email"`
}

type NumbersConfirmReq struct {
	Numbers int `json:"numbers"`
}

type PasswordResetReq struct {
	Password        string `json:"password" validate:"required,min=3,max=130"`
	PasswordConfirm string `json:"passwordConfirm"`
}
