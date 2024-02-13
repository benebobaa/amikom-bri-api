package entity

import (
	"database/sql"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              string         `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username        string         `gorm:"column:username;"`
	Email           string         `gorm:"column:email"`
	FullName        string         `gorm:"column:full_name"`
	HashedPassword  string         `gorm:"column:hashed_password"`
	IsEmailVerified bool           `gorm:"column:is_email_verified"`
	Account         Account        `gorm:"foreignKey:UserID;references:ID"`
	HashedPin       sql.NullString `gorm:"column:hashed_pin"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
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
		Account:         u.Account.ToAccountResponse(),
		CreatedAt:       u.CreatedAt,
	}
}

func (u *User) ToLoginResponseWithToken(sessionResp *response.SessionsResponse) *response.LoginResponse {
	return &response.LoginResponse{
		Token: sessionResp,
		User:  u.ToUserResponse(),
	}
}

func (u *User) ToUserProfileResponse() *response.UserProfileResponse {
	return &response.UserProfileResponse{
		Username:        u.Username,
		Email:           u.Email,
		FullName:        u.FullName,
		IsEmailVerified: u.IsEmailVerified,
		Account:         u.Account.ToAccountResponse(),
	}
}

func ToUserResponses(users []User, pagingMetadata *response.PageMetaData) *response.UserResponses {
	var userResponses []response.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *user.ToUserResponse())
	}
	return &response.UserResponses{
		Users:  userResponses,
		Paging: pagingMetadata,
	}
}
