package request

type ForgotPasswordRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	NewPassword     string `json:"new_password" validate:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword,min=8,max=32"`
}
