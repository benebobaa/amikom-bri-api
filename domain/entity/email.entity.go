package entity

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Email struct {
	ID         int            `gorm:"column:id;primaryKey;autoIncrement;not null" `
	UserID     uuid.UUID      `gorm:"column:user_id;not null"`
	Email      string         `gorm:"column:email;not null"`
	SecretCode string         `gorm:"column:secret_code;not null"`
	IsUsed     bool           `gorm:"column:is_used;not null;default:false"`
	ExpiredAt  time.Time      `gorm:"column:expired_at"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at"`
	DeletedAT  gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (e *Email) TableName() string {
	return "verify_emails"
}

func (e *Email) ToEmailVerifyResponse() *response.EmailVerifyResponse {
	return &response.EmailVerifyResponse{
		Email:           e.Email,
		IsEmailVerified: e.IsUsed,
	}
}
