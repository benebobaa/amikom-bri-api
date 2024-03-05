package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type NotificationRepository interface {
	NotificationCreate(tx *gorm.DB, value *entity.Notification) error
	FindAllNotif(tx *gorm.DB, userId uuid.UUID) ([]entity.Notification, error)
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

func (n *notificationRepositoryImpl) FindAllNotif(tx *gorm.DB, userId uuid.UUID) ([]entity.Notification, error) {
	var notif []entity.Notification
	result := tx.Where("user_id = ?", userId).Find(&notif)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find all notification : %v", result.Error))
		return nil, result.Error
	}

	return notif, nil
}
