package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"gorm.io/gorm"
	"log"
)

type EmailRepository interface {
	EmailVerifyCreate(tx *gorm.DB, value *entity.Email) error
	FindEmailBySecretCode(tx *gorm.DB, secretCode string) (*entity.Email, error)
	UpdateEmailVerify(tx *gorm.DB, email *entity.Email) error
}

type emailRepositoryImpl struct {
}

func NewEmailRepository() EmailRepository {
	return &emailRepositoryImpl{}
}

func (e *emailRepositoryImpl) EmailVerifyCreate(tx *gorm.DB, value *entity.Email) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create email verify : %v", result.Error))
		return result.Error
	}

	return nil
}

func (e *emailRepositoryImpl) FindEmailBySecretCode(tx *gorm.DB, secretCode string) (*entity.Email, error) {
	var email entity.Email
	result := tx.Where("secret_code = ?", secretCode).First(&email)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find email by secret code : %v", result.Error))
		return nil, result.Error
	}

	return &email, nil
}

func (e *emailRepositoryImpl) UpdateEmailVerify(tx *gorm.DB, email *entity.Email) error {
	result := tx.Save(email)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when update email verify : %v", result.Error))
		return result.Error
	}

	return nil
}
