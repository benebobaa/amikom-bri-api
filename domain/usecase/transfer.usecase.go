package usecase

import (
	"context"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/benebobaa/amikom-bri-api/domain/repository"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/mail"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type TransferUsecase interface {
	TransferMoney(ctx context.Context, requestData *request.TransferRequest, userID string) error
}

type transferUsecaseImpl struct {
	DB                 *gorm.DB
	Validate           *validator.Validate
	TitanMail          mail.EmailSender
	TransferRepository repository.TransferRepository
	AccountRepository  repository.AccountRepository
	EntryRepository    repository.EntryRepository
}

func NewTransferUsecase(db *gorm.DB, validate *validator.Validate, titanMail mail.EmailSender, transferRepository repository.TransferRepository,
	accountRepository repository.AccountRepository, entryRepository repository.EntryRepository) TransferUsecase {
	return &transferUsecaseImpl{
		DB:                 db,
		Validate:           validate,
		TitanMail:          titanMail,
		TransferRepository: transferRepository,
		AccountRepository:  accountRepository,
		EntryRepository:    entryRepository,
	}
}

func (t *transferUsecaseImpl) TransferMoney(ctx context.Context, requestData *request.TransferRequest, userID string) error {
	tx := t.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := t.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return err
	}

	// Find account by id
	account, err := t.AccountRepository.FindByID(tx, requestData.FromAccountID)
	if err != nil {
		log.Printf("Error when find to account with id : %+v", err)
		return util.AccounDoesNotExist
	}

	// Check if account not belong to user
	if account.UserID != userID {
		log.Printf("Account not belong to user")
		return util.AccountNotBelongToUser
	}

	// Get to account id email
	toAccount, err := t.AccountRepository.FindByID(tx, requestData.ToAccountID)
	if err != nil {
		log.Printf("Error when find to account with id : %+v", err)
		return util.DestinationAccountNotExist
	}

	// Check pin before transaction
	isValid := util.CheckPassword(requestData.Pin, account.User.HashedPin)

	if !isValid {
		log.Printf("Pin not valid")
		return util.InvalidPin
	}

	log.Printf("account user: %+v", account.User)
	if err != nil {
		log.Printf("Error when find account by id : %+v", err)
		return err
	}

	// Check if account balance is sufficient
	if account.Balance < requestData.Amount {
		log.Printf("Insufficient balance")
		return util.InsufficientBalance
	}

	// Create transfer
	transferEntity := requestData.ToEntity()
	err = t.TransferRepository.TransferCreate(tx, transferEntity)

	if err != nil {
		log.Printf("Error when create transfer : %+v", err)
	}

	// Create entry out
	entryEntityOut := &entity.Entry{
		AccountID: requestData.FromAccountID,
		Amount:    -requestData.Amount,
		EntryType: entity.EntryOut,
	}

	err = t.EntryRepository.EntryCreate(tx, entryEntityOut)

	if err != nil {
		log.Printf("Error when create entry out : %+v", err)
		return err
	}

	// Create entry in
	entryEntityIn := &entity.Entry{
		AccountID: requestData.ToAccountID,
		Amount:    requestData.Amount,
		EntryType: entity.EntryIn,
	}

	err = t.EntryRepository.EntryCreate(tx, entryEntityIn)

	if err != nil {
		log.Printf("Error when create entry in : %+v", err)
		return err
	}

	// Updates balance account
	if requestData.FromAccountID < requestData.ToAccountID {
		err = t.updateAccountBalance(tx, account.User.Email, toAccount.User.Email, requestData.FromAccountID, -requestData.Amount, requestData.ToAccountID, requestData.Amount)
		if err != nil {
			log.Printf("Error when update balance : %+v", err)
			return err
		}
	} else {
		err = t.updateAccountBalance(tx, account.User.Email, toAccount.User.Email, requestData.ToAccountID, requestData.Amount, requestData.FromAccountID, -requestData.Amount)
		if err != nil {
			log.Printf("Error when update balance : %+v", err)
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (t *transferUsecaseImpl) updateAccountBalance(
	tx *gorm.DB,
	fromUser string,
	toUser string,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64) error {

	// Send notification emails to both sender and receiver
	if amount1 > amount2 {
		go func() {
			subject, content, toEmail := mail.GetReceiverParamTransferNotification(toUser, fromUser, amount1)
			err := t.TitanMail.SendEmail(subject, content, toEmail, []string{}, []string{}, []string{})
			if err != nil {
				log.Printf("Failed send email : %+v", err)
			}
		}()

		go func() {
			subject, content, toEmail := mail.GetSenderParamTransferNotification(fromUser, amount1)
			err := t.TitanMail.SendEmail(subject, content, toEmail, []string{}, []string{}, []string{})
			if err != nil {
				log.Printf("Failed send email : %+v", err)
			}
		}()
	} else {
		go func() {
			subject, content, toEmail := mail.GetReceiverParamTransferNotification(toUser, fromUser, amount2)
			err := t.TitanMail.SendEmail(subject, content, toEmail, []string{}, []string{}, []string{})
			if err != nil {
				log.Printf("Failed send email : %+v", err)
			}
		}()

		go func() {
			subject, content, toEmail := mail.GetSenderParamTransferNotification(fromUser, amount2)
			err := t.TitanMail.SendEmail(subject, content, toEmail, []string{}, []string{}, []string{})
			if err != nil {
				log.Printf("Failed send email : %+v", err)
			}
		}()
	}

	err := t.AccountRepository.AddAccountBalance(tx, accountID1, amount1)
	if err != nil {
		log.Printf("Error when add account1 balance : %+v", err)
		return err
	}

	err = t.AccountRepository.AddAccountBalance(tx, accountID2, amount2)
	if err != nil {
		log.Printf("Error when add account2 balance : %+v", err)
		return err
	}
	return nil
}
