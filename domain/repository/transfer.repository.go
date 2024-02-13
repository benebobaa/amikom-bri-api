package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"gorm.io/gorm"
	"log"
)

type TransferRepository interface {
	TransferCreate(tx *gorm.DB, value *entity.Transfer) error
	FindByID(tx *gorm.DB, id int64) (*entity.Transfer, error)
}

type transferRepositoryImpl struct {
}

func NewTransferRepository() TransferRepository {
	return &transferRepositoryImpl{}
}

func (t *transferRepositoryImpl) TransferCreate(tx *gorm.DB, value *entity.Transfer) error {

	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create transfer : %v", result.Error))
		return result.Error
	}

	return nil
}

func (t *transferRepositoryImpl) FindByID(tx *gorm.DB, id int64) (*entity.Transfer, error) {
	var transfer entity.Transfer

	result := tx.Where("id = ?", id).First(&transfer)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find transfer by id : %v", result.Error))
		return nil, result.Error
	}

	return &transfer, nil
}
