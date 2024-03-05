package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type ForgotPasswordRepository interface {
	ForgotPasswordCreate(tx *gorm.DB, value *entity.ForgotPassword) error
	FindByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.ForgotPassword, error)
	UpdateForgotPassword(tx *gorm.DB, value *entity.ForgotPassword) error
}

type forgotPasswordRepositoryImpl struct {
}

func NewForgotPasswordRepository() ForgotPasswordRepository {
	return &forgotPasswordRepositoryImpl{}
}

func (f *forgotPasswordRepositoryImpl) ForgotPasswordCreate(tx *gorm.DB, value *entity.ForgotPassword) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create forgot password : %v", result.Error))
		return result.Error
	}

	return nil
}

func (f *forgotPasswordRepositoryImpl) FindByUserID(tx *gorm.DB, userID uuid.UUID) (*entity.ForgotPassword, error) {
	var forgotPassword entity.ForgotPassword
	result := tx.Where("user_id = ?", userID).Order("id DESC").First(&forgotPassword)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find by user id : %v", result.Error))
		return nil, result.Error
	}

	return &forgotPassword, nil
}

func (f *forgotPasswordRepositoryImpl) UpdateForgotPassword(tx *gorm.DB, value *entity.ForgotPassword) error {
	result := tx.Model(value).Where("user_id = ?", value.UserID).Updates(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when update forgot password : %v", result.Error))
		return result.Error
	}

	return nil
}
