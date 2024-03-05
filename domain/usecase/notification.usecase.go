package usecase

import (
	"context"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/benebobaa/amikom-bri-api/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationUsecase interface {
	GetAllNotifHistory(ctx context.Context, userId uuid.UUID) ([]entity.Notification, error)
}

type notificationUsecaseImpl struct {
	DB                     *gorm.DB
	NotificationRepository repository.NotificationRepository
}

func NewNotificationUsecase(db *gorm.DB, notificationRepository repository.NotificationRepository) NotificationUsecase {
	return &notificationUsecaseImpl{
		DB:                     db,
		NotificationRepository: notificationRepository,
	}
}

func (n *notificationUsecaseImpl) GetAllNotifHistory(ctx context.Context, userId uuid.UUID) ([]entity.Notification, error) {
	tx := n.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	notif, err := n.NotificationRepository.FindAllNotif(tx, userId)
	if err != nil {
		return nil, err
	}

	return notif, nil
}
