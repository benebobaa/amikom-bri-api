package request

import (
	"github.com/benebobaa/amikom-bri-api/domain/entity"
)

type UserRegisterRequest struct {
	Username        string `json:"username"  validate:"required"`
	FullName        string `json:"full_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Pin             string `json:"pin" validate:"required,len=6,number"`
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
