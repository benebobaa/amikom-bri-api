package entity

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID        int64          `gorm:"column:id;primaryKey;"`
	UserID    uuid.UUID      `gorm:"column:user_id;not null"`
	Balance   int64          `gorm:"column:balance;not null"`
	User      *User          `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (a *Account) TableName() string {
	return "accounts"
}

func (a *Account) ToAccountResponseLogin(username string) *response.AccountResponse {
	return &response.AccountResponse{
		ID:        a.ID,
		Username:  username,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
	}
}

func (a *Account) ToAccountResponse() *response.AccountResponse {
	return &response.AccountResponse{
		ID:        a.ID,
		Username:  a.User.Username,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
	}
}
