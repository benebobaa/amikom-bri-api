package request

import "github.com/benebobaa/amikom-bri-api/domain/entity"

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,max=32"`
}

func (u *LoginRequest) ToUserEntity() *entity.User {
	return &entity.User{
		Username: u.UsernameOrEmail,
		Email:    u.UsernameOrEmail,
	}
}
