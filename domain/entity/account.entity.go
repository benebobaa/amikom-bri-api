package entity

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID        int            `gorm:"column:id;primaryKey;"`
	UserID    string         `gorm:"column:user_id;not null"`
	Balance   int64          `gorm:"column:balance;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (a *Account) TableName() string {
	return "accounts"
}

func (a *Account) ToAccountResponse() *response.AccountResponse {
	return &response.AccountResponse{
		ID:        a.ID,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
	}
}