package entity

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"gorm.io/gorm"
	"time"
)

type Transfer struct {
	ID            int64          `gorm:"column:id;primaryKey"`
	FromAccountID int64          `gorm:"column:from_account_id"`
	ToAccountID   int64          `gorm:"column:to_account_id"`
	Amount        int64          `gorm:"column:amount"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (t *Transfer) TableName() string {
	return "transfers"
}

func (t *Transfer) ToTransfeInfo() *response.TranferInfo {
	return &response.TranferInfo{
		Amount:        t.Amount,
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		CreatedAt:     t.CreatedAt,
	}
}
