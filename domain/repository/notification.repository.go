package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"gorm.io/gorm"
	"log"
)

type NotificationRepository interface {
	NotificationCreate(tx *gorm.DB, value *entity.Notification) error
}

type notificationRepositoryImpl struct {
}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepositoryImpl{}
}

func (n *notificationRepositoryImpl) NotificationCreate(tx *gorm.DB, value *entity.Notification) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create notification : %v", result.Error))
		return result.Error
	}

	return nil

}
