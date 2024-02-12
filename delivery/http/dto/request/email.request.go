package request

import (
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"time"
)

type EmailRequest struct {
	UserID     string    `json:"user_id" validate:"required"`
	Username   string    `json:"username" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	SecretCode string    `json:"secret_code" validate:"required"`
	ExpiredAt  time.Time `json:"expired_at" validate:"required"`
}

func (e *EmailRequest) ToEntity() *entity.Email {
	return &entity.Email{
		UserID:     e.UserID,
		Email:      e.Email,
		SecretCode: e.SecretCode,
		ExpiredAt:  e.ExpiredAt,
	}
}
