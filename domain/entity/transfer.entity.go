package entity

import (
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
