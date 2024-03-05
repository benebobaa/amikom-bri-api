package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type ExpensesPlanRepository interface {
	ExpensesPlanCreate(tx *gorm.DB, value *entity.ExpensesPlan) error
	FindByID(tx *gorm.DB, value *entity.ExpensesPlan) error
	Delete(tx *gorm.DB, value *entity.ExpensesPlan) error
	Update(tx *gorm.DB, value *entity.ExpensesPlan) error
	FindAlL(tx *gorm.DB, userId uuid.UUID) ([]entity.ExpensesPlan, error)
	FindAllWithFilter(tx *gorm.DB, request *request.SearchPaginationRequest, userId uuid.UUID) ([]entity.ExpensesPlan, int64, error)
	FilterExpenses(request *request.SearchPaginationRequest) func(tx *gorm.DB) *gorm.DB
}

type expensesPlanRepositoryImpl struct {
}

func NewExpensesPlanRepository() ExpensesPlanRepository {
	return &expensesPlanRepositoryImpl{}
}

func (e *expensesPlanRepositoryImpl) ExpensesPlanCreate(tx *gorm.DB, value *entity.ExpensesPlan) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create expenses plan : %v", result.Error))
		return result.Error
	}

	return nil
}

func (e *expensesPlanRepositoryImpl) FindByID(tx *gorm.DB, value *entity.ExpensesPlan) error {

	result := tx.Where("user_id = ?", value.UserID).First(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find expenses plan by id : %v", result.Error))
		return result.Error
	}

	return nil
}

func (e *expensesPlanRepositoryImpl) Delete(tx *gorm.DB, value *entity.ExpensesPlan) error {

	result := tx.Delete(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when delete expenses plan : %v", result.Error))
		return result.Error
	}

	return nil
}

func (e *expensesPlanRepositoryImpl) Update(tx *gorm.DB, value *entity.ExpensesPlan) error {
	result := tx.Model(value).Where("user_id = ?", value.UserID).Updates(value).First(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when update user : %v", result.Error))
		return result.Error
	}

	return nil
}

func (e *expensesPlanRepositoryImpl) FindAlL(tx *gorm.DB, userId uuid.UUID) ([]entity.ExpensesPlan, error) {
	var expenses []entity.ExpensesPlan

	result := tx.Where("user_id = ?", userId).Find(&expenses)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find expenses plan : %v", result.Error))
		return nil, result.Error
	}

	return expenses, nil
}

func (e *expensesPlanRepositoryImpl) FindAllWithFilter(db *gorm.DB, request *request.SearchPaginationRequest, userId uuid.UUID) ([]entity.ExpensesPlan, int64, error) {
	var expenses []entity.ExpensesPlan

	if err := db.Scopes(e.FilterExpenses(request)).Offset((request.Page-1)*request.Size).Limit(request.Size).Where("user_id", userId).Find(&expenses).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.ExpensesPlan{}).Scopes(e.FilterExpenses(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return expenses, total, nil
}

func (p *expensesPlanRepositoryImpl) FilterExpenses(request *request.SearchPaginationRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if keyword := request.Keyword; keyword != "" {
			keyword = "%" + keyword + "%"
			tx = tx.Where("title LIKE ?", keyword)
		}

		return tx
	}
}
