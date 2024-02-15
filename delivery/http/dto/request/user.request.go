package request

import (
	"github.com/benebobaa/amikom-bri-api/domain/entity"
)

type UserRegisterRequest struct {
	Username        string `json:"username"  validate:"required,nowhitespace"`
	FullName        string `json:"full_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password,min=8,max=32"`
	Pin             string `json:"pin" validate:"required,len=6,number"`
}

type UserUpdateRequest struct {
	Username string `json:"username"  validate:"required,nowhitespace"`
	FullName string `json:"full_name" validate:"required"`
}

func (u *UserUpdateRequest) ToEntity() *entity.User {
	return &entity.User{
		Username: u.Username,
		FullName: u.FullName,
	}
}

func (u *UserRegisterRequest) ToEntity() *entity.User {
	return &entity.User{
		Username:       u.Username,
		Email:          u.Email,
		FullName:       u.FullName,
		HashedPassword: u.Password,
		HashedPin:      u.Pin,
	}
}
