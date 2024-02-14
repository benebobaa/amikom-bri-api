package entity

import (
	"fmt"
	"time"
)

var (
	NotificationCategoryIn  = "Transfer In"
	NotificationCategoryOut = "Transfer Out"
	NotificationTitle       = "Transfer Notification"
	NotificationDescIn      = "You have received a transfer with amount of "
	NotificationDescOut     = "You have sent a transfer with amount of "
)

type Notification struct {
	ID          int64     `gorm:"column:id"`
	UserID      string    `gorm:"column:user_id"`
	Title       string    `gorm:"title"`
	Description string    `gorm:"description"`
	Category    string    `gorm:"category"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at"`
}

func (n *Notification) TableName() string {
	return "notifications"
}

func GetNotificationInEntity(userId string, amount int64) *Notification {
	return &Notification{
		UserID:      userId,
		Title:       NotificationTitle,
		Category:    NotificationCategoryIn,
		Description: fmt.Sprintf("%s %d", NotificationDescIn, amount),
	}
}

func GetNotificationOutEntity(userId string, amount int64) *Notification {
	return &Notification{
		UserID:      userId,
		Title:       NotificationTitle,
		Category:    NotificationCategoryOut,
		Description: fmt.Sprintf("%s %d", NotificationDescOut, amount),
	}
}
