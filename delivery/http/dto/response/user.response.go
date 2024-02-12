package response

import "time"

type UserResponse struct {
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at"`
}
