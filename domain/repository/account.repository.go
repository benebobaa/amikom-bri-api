package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"gorm.io/gorm"
	"log"
)

type AccountRepository interface {
	AccountCreate(tx *gorm.DB, value *entity.Account) error
	FindByUserID(tx *gorm.DB, userID string) (*entity.Account, bool, error)
	FindByID(tx *gorm.DB, id int64) (*entity.Account, error)
	AddAccountBalance(tx *gorm.DB, accountID int64, amount int64) error
}

type accountRepositoryImpl struct {
}

func NewAccountRepository() AccountRepository {
	return &accountRepositoryImpl{}
}

func (a *accountRepositoryImpl) AccountCreate(tx *gorm.DB, value *entity.Account) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create account : %v", result.Error))
		return result.Error
	}

	return nil
}

func (a *accountRepositoryImpl) FindByUserID(tx *gorm.DB, userID string) (*entity.Account, bool, error) {
	var account entity.Account

	result := tx.Where("user_id = ?", userID).First(&account)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find account by user id : %v", result.Error))
		return nil, false, result.Error
	}

	return &account, true, nil
}

func (a *accountRepositoryImpl) FindByID(tx *gorm.DB, id int64) (*entity.Account, error) {
	var account entity.Account

	result := tx.Where("id = ?", id).Preload("User").First(&account)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find account by id : %v", result.Error))
		return nil, result.Error
	}

	return &account, nil
}

func (a *accountRepositoryImpl) AddAccountBalance(tx *gorm.DB, accountID int64, amount int64) error {
	result := tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, accountID)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when add account balance : %v", result.Error))
		return result.Error
	}

	return nil
}
