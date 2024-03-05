package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID           string         `gorm:"column:id;primaryKey;"`
	UserID       uuid.UUID      `gorm:"column:user_id;not null"`
	RefreshToken string         `gorm:"column:refresh_token;not null"`
	UserAgent    string         `gorm:"column:user_agent;not null"`
	ClientIP     string         `gorm:"column:client_ip;not null"`
	IsBlocked    bool           `gorm:"column:is_blocked;not null"`
	ExpiredAt    time.Time      `gorm:"column:expired_at;not null"`
	CreatedAt    time.Time      `gorm:"column:created_at;not null"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (s *Session) TableName() string {
	return "sessions"
}
