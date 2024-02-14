package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"gorm.io/gorm"
	"log"
)

type EntryRepository interface {
	EntryCreate(tx *gorm.DB, value *entity.Entry) error
	FindByID(tx *gorm.DB, id int64) (*entity.Entry, error)
	DeleteEntry(tx *gorm.DB, value *entity.Entry) error
	FindAll(db *gorm.DB, request *request.SearchPaginationRequest, accountID int64) ([]entity.Entry, int64, error)
	FindAllFilterDate(db *gorm.DB, request *request.SearchPaginationRequest, accountID int64) ([]entity.Entry, int64, error)
	FilterEntries(request *request.SearchPaginationRequest) func(tx *gorm.DB) *gorm.DB
}

type entryRepositoryImpl struct {
}

func NewEntryRepository() EntryRepository {
	return &entryRepositoryImpl{}
}

func (e *entryRepositoryImpl) EntryCreate(tx *gorm.DB, value *entity.Entry) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create entry : %v", result.Error))
		return result.Error
	}

	return nil

}

func (e *entryRepositoryImpl) FindAll(db *gorm.DB, request *request.SearchPaginationRequest, accountID int64) ([]entity.Entry, int64, error) {
	var entries []entity.Entry

	err := db.Scopes(e.FilterEntries(request)).Offset((request.Page-1)*request.Size).Limit(request.Size).Where("account_id", accountID).Find(&entries).Error

	if err != nil {
		return nil, 0, err
	}

	var total int64 = 0

	err = db.Model(&entity.Entry{}).Scopes(e.FilterEntries(request)).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (e *entryRepositoryImpl) FindByID(tx *gorm.DB, id int64) (*entity.Entry, error) {
	var entry entity.Entry

	result := tx.Where("id = ?", id).First(&entry)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find entry by id : %v", result.Error))
		return nil, result.Error
	}

	return &entry, nil
}

func (e *entryRepositoryImpl) DeleteEntry(tx *gorm.DB, value *entity.Entry) error {
	result := tx.Delete(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when delete entry : %v", result.Error))
		return result.Error
	}

	return nil
}

func (e *entryRepositoryImpl) FindAllFilterDate(db *gorm.DB, request *request.SearchPaginationRequest, accountID int64) ([]entity.Entry, int64, error) {
	var entries []entity.Entry

	err := db.Scopes(e.FilterEntries(request)).Preload("Account").Preload("Account.User").Offset((request.Page-1)*request.Size).Limit(request.Size).Where("account_id", accountID).Find(&entries).Error

	if err != nil {
		return nil, 0, err
	}

	var total int64 = 0

	err = db.Model(&entity.Entry{}).Scopes(e.FilterEntries(request)).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (e *entryRepositoryImpl) FilterEntries(requestParam *request.SearchPaginationRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if keyword := requestParam.Keyword; keyword != "" {
			keyword = "%" + keyword + "%"
			tx = tx.Where("entry_type LIKE ?", keyword)
		}

		if filter := requestParam.Filter; filter != "" {

			switch filter {

			case request.DAILY:
				tx = tx.Where("DATE(date) = ?", requestParam.Date)
			case request.MONTHLY:
				tx = tx.Where("EXTRACT(YEAR FROM date) = ? AND EXTRACT(MONTH FROM date) = ?", requestParam.ParsedDate.Year(), requestParam.ParsedDate.Month())
			case request.YEARLY:

				log.Println("date year", requestParam.ParsedDate)
				tx = tx.Where("EXTRACT(YEAR FROM date) = ?", requestParam.ParsedDate.Year())
			}
		}

		return tx
	}
}
