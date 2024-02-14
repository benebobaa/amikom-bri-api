package usecase

import (
	"context"
	"errors"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/benebobaa/amikom-bri-api/domain/repository"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"math"
)

type ExpensesPlanUsecase interface {
	ExpensesPlanCreate(ctx context.Context, requestData *request.ExpensesPlanRequest, userId string) (*response.ExpensesPlanResponse, error)
	DeletePlan(ctx context.Context, id int64, userId string) error
	UpdatePlan(ctx context.Context, requestData *request.ExpensesPlanUpdateRequest, id int64, userId string) (*response.ExpensesPlanResponse, error)
	FindAllFilter(ctx context.Context, requestData *request.SearchPaginationRequest, userID string) (*response.ExpensesPlanResponses, error)
	ExpensesPlanExportPdf(ctx context.Context, userId string) (string, error)
}

type expensesPlanUsecaseImpl struct {
	DB                 *gorm.DB
	Validate           *validator.Validate
	GoPdf              *util.PDFGenerator
	ExpensesRepository repository.ExpensesPlanRepository
}

func NewExpensesPlanUsecase(db *gorm.DB, validate *validator.Validate, goPdf *util.PDFGenerator, expensesRepository repository.ExpensesPlanRepository) ExpensesPlanUsecase {
	return &expensesPlanUsecaseImpl{
		DB:                 db,
		Validate:           validate,
		GoPdf:              goPdf,
		ExpensesRepository: expensesRepository,
	}
}

func (e *expensesPlanUsecaseImpl) ExpensesPlanCreate(ctx context.Context, requestData *request.ExpensesPlanRequest, userId string) (*response.ExpensesPlanResponse, error) {
	tx := e.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := e.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	// Convert request data to entity, then create expenses plan
	expensesEntity := requestData.ToEntity(userId)
	err = e.ExpensesRepository.ExpensesPlanCreate(tx, expensesEntity)

	if err != nil {
		log.Printf("Failed create expenses plan : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return expensesEntity.ToExpensesResponse(), nil

}

func (e *expensesPlanUsecaseImpl) DeletePlan(ctx context.Context, id int64, userId string) error {
	tx := e.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := e.Validate.Var(id, "required")
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return err
	}

	expensesPlanEntity := &entity.ExpensesPlan{
		ID:     id,
		UserID: userId,
	}

	err = e.ExpensesRepository.FindByID(tx, expensesPlanEntity)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Expenses plan not found : %+v", err)
			return util.ExpensesPlanNotFound
		}
		log.Printf("Failed find expenses plan by id : %+v", err)
		return err
	}

	err = e.ExpensesRepository.Delete(tx, expensesPlanEntity)

	if err != nil {
		log.Printf("Failed delete expenses plan : %+v", err)
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (e *expensesPlanUsecaseImpl) UpdatePlan(ctx context.Context, requestData *request.ExpensesPlanUpdateRequest, id int64, userId string) (*response.ExpensesPlanResponse, error) {
	tx := e.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := e.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	// Find expenses plan by id
	err = e.ExpensesRepository.FindByID(tx, &entity.ExpensesPlan{ID: id, UserID: userId})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Expenses plan not found : %+v", err)
			return nil, util.ExpensesPlanNotFound
		}
		log.Printf("Failed find expenses plan by id : %+v", err)
		return nil, err
	}

	// Update expenses plan
	expensesPlanEntity := requestData.ToEntity(userId)
	expensesPlanEntity.ID = id

	err = e.ExpensesRepository.Update(tx, expensesPlanEntity)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Expenses plan not found : %+v", err)
			return nil, util.ExpensesPlanNotFound
		}
		log.Printf("Failed update expenses plan by id : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return expensesPlanEntity.ToExpensesResponse(), nil
}

func (e *expensesPlanUsecaseImpl) FindAllFilter(ctx context.Context, requestData *request.SearchPaginationRequest, userID string) (*response.ExpensesPlanResponses, error) {
	tx := e.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Get all users with pagination and search
	expensePlans, total, err := e.ExpensesRepository.FindAllWithFilter(tx, requestData, userID)

	if err != nil {
		log.Printf("Failed find all users : %+v", err)
		return nil, err
	}

	// Calculate the total page pagination
	resultPaging := &response.PageMetaData{
		Page:      requestData.Page,
		Size:      requestData.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(requestData.Size))),
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return entity.ToExpensesResponses(expensePlans, resultPaging), nil
}

func (e *expensesPlanUsecaseImpl) ExpensesPlanExportPdf(ctx context.Context, userId string) (string, error) {
	tx := e.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Get all expenses plan
	expensePlans, err := e.ExpensesRepository.FindAlL(tx, userId)

	if err != nil {
		log.Printf("Failed find all expenses plan : %+v", err)
		return "", err
	}

	// Check if expenses plan empty, return err
	if len(expensePlans) == 0 {
		return "", util.CannotExportEmptyData
	}

	// Generate pdf expenses plan
	filename, err := e.GoPdf.GeneratePdfExpensePlans(expensePlans)

	if err != nil {
		log.Printf("Failed generate pdf expenses plan : %+v", err)
		return "", err
	}

	return filename, nil
}
