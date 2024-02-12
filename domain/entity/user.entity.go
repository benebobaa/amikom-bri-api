package entity

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"time"
)

type User struct {
	Username        string    `gorm:"column:username;primaryKey"`
	Email           string    `gorm:"column:email"`
	FullName        string    `gorm:"column:full_name"`
	HashedPassword  string    `gorm:"column:hashed_password"`
	IsEmailVerified bool      `gorm:"column:is_email_verified"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
	DeletedAt       time.Time `gorm:"column:deleted_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) ToUserResponse() *response.UserResponse {
	return &response.UserResponse{
		Username:        u.Username,
		Email:           u.Email,
		FullName:        u.FullName,
		IsEmailVerified: u.IsEmailVerified,
		CreatedAt:       u.CreatedAt,
	}
}

func (u *User) ToLoginResponseWithToken(sessionResp *response.SessionsResponse) *response.LoginResponse {
	return &response.LoginResponse{
		Token: sessionResp,
		User:  u.ToUserResponse(),
	}
}
