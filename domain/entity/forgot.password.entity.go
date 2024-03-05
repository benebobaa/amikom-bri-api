package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ForgotPassword struct {
	ID             int            `gorm:"primaryKey;autoIncrement"`
	UserID         uuid.UUID      `gorm:"column:user_id;not null"`
	ResetToken     string         `gorm:"column:reset_token;not null"`
	IsUsed         bool           `gorm:"column:is_used;not null;default:false"`
	RequestTime    time.Time      `gorm:"column:request_timestamp;not null;default:CURRENT_TIMESTAMP"`
	ExpirationTime time.Time      `gorm:"column:expiration_timestamp;not null"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (f *ForgotPassword) TableName() string {
	return "forgot_password"
}
