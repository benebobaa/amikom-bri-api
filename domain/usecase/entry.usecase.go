package usecase

import (
	"context"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/benebobaa/amikom-bri-api/domain/repository"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"math"
	"strings"
	"time"
)

type EntryUsecase interface {
	FindAllHistoryTransfer(ctx context.Context, requestData *request.SearchPaginationRequest, userID uuid.UUID) (*response.EntryResponses, error)
	DeleteEntry(ctx context.Context, entryID int64, userID uuid.UUID) error
	FindAllFilterDate(ctx context.Context, requestData *request.SearchPaginationRequest, userID uuid.UUID) (*response.EntryResponses, string, error)
}

type entryUsecaseImpl struct {
	DB                *gorm.DB
	GoPdf             *util.PDFGenerator
	EntryRepository   repository.EntryRepository
	AccountRepository repository.AccountRepository
}

func NewEntryUsecase(db *gorm.DB, pdf *util.PDFGenerator, entryRepository repository.EntryRepository, accountRepository repository.AccountRepository) EntryUsecase {
	return &entryUsecaseImpl{
		DB:                db,
		GoPdf:             pdf,
		EntryRepository:   entryRepository,
		AccountRepository: accountRepository,
	}
}

func (e *entryUsecaseImpl) FindAllHistoryTransfer(ctx context.Context, requestData *request.SearchPaginationRequest, userID uuid.UUID) (*response.EntryResponses, error) {

	tx := e.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Find account with user id
	account, _, err := e.AccountRepository.FindByUserID(tx, userID)

	if err != nil {
		log.Printf("Error when find account by user id : %v", err)
		return nil, err
	}

	// Get all data with account id
	entries, total, err := e.EntryRepository.FindAll(tx, requestData, account.ID)

	if err != nil {
		log.Printf("Failed find all entries : %+v", err)
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

	return entity.ToEntryResponses(entries, resultPaging), nil
}

func (e *entryUsecaseImpl) DeleteEntry(ctx context.Context, entryID int64, userID uuid.UUID) error {
	return nil
}

func (e *entryUsecaseImpl) FindAllFilterDate(ctx context.Context, requestData *request.SearchPaginationRequest, userID uuid.UUID) (*response.EntryResponses, string, error) {
	tx := e.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Find account by user id
	account, _, err := e.AccountRepository.FindByUserID(tx, userID)

	if err != nil {
		log.Printf("Error when find account by user id : %v", err)
		return nil, "", err
	}

	// Tolower filter param daily,monthly,yearly
	lowerFilter := strings.ToLower(requestData.Filter)

	// Parsed date
	parsedDate, err := e.parsedDateByFilter(lowerFilter, requestData.Date)
	if err != nil {
		log.Printf("Error when parsed date : %v", err)
		return nil, "", err
	}

	requestData.ParsedDate = parsedDate
	// Get all with filtered date
	log.Println("parsed date", parsedDate)
	entries, total, err := e.EntryRepository.FindAllFilterDate(tx, requestData, account.ID)

	if err != nil {
		log.Printf("Failed find all entries : %+v", err)
		return nil, "", err
	}

	if requestData.ExportPdf {

		// Check entries not empty
		if len(entries) == 0 {
			return nil, "", util.CannotExportEmptyData
		}

		//Calculate total in and out
		totalIn, totalOut := e.CalculateTotalInAndOut(entries)

		//err = e.GoPdf.GeneratePdf(entries, totalIn, totalOut)
		fileName, err := e.GoPdf.GeneratePdf(entries, totalIn, totalOut)
		if err != nil {
			log.Printf("Error when generate transaction_pdf : %v", err)
			return nil, "", err
		}

		return nil, fileName, nil
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
		return nil, "", err
	}

	return entity.ToEntryResponses(entries, resultPaging), "", nil
}

// Parsed date by filter
func (e *entryUsecaseImpl) parsedDateByFilter(filter string, date string) (time.Time, error) {

	switch filter {
	case request.DAILY:
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			log.Printf("Error when parse date daily : %v", err)
			return parsedDate, util.DateFormatNotValid
		}
		return parsedDate, nil
	case request.MONTHLY:
		parsedDate, err := time.Parse("2006-01", date)
		if err != nil {
			log.Printf("Error when parse date monthly : %v", err)
			return parsedDate, util.DateFormatNotValid
		}
		return parsedDate, nil
	case request.YEARLY:
		parsedDate, err := time.Parse("2006", date)
		if err != nil {
			log.Printf("Error when parse date yearly : %v", err)
			return parsedDate, util.DateFormatNotValid
		}
		return parsedDate, nil
	}

	return time.Now(), util.DateFormatNotValid
}

func (e *entryUsecaseImpl) CalculateTotalInAndOut(entries []entity.Entry) (int64, int64) {
	var totalIn, totalOut int64
	for _, entry := range entries {
		if entry.EntryType == entity.EntryIn {
			totalIn += entry.Amount
		} else {
			totalOut += entry.Amount
		}
	}
	return totalIn, totalOut
}
