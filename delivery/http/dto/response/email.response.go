package response

type EmailVerifyResponse struct {
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
}
