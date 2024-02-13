package response

import "time"

type UserResponse struct {
	Username        string           `json:"username"`
	FullName        string           `json:"full_name"`
	Email           string           `json:"email"`
	IsEmailVerified bool             `json:"is_email_verified"`
	Account         *AccountResponse `json:"account"`
	CreatedAt       time.Time        `json:"created_at"`
}

type UserProfileResponse struct {
	Username        string           `json:"username"`
	FullName        string           `json:"full_name"`
	Email           string           `json:"email"`
	IsEmailVerified bool             `json:"is_email_verified"`
	Account         *AccountResponse `json:"account"`
}

type UserResponses struct {
	Users  []UserResponse `json:"users"`
	Paging *PageMetaData  `json:"paging"`
}
