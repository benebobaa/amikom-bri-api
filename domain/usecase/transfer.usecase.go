package usecase

import (
	"context"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/benebobaa/amikom-bri-api/domain/repository"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/mail"
	"github.com/benebobaa/amikom-bri-api/util/onesignal"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type TransferUsecase interface {
	TransferMoney(ctx context.Context, requestData *request.TransferRequest, userID string) (*response.TransferResponse, error)
}

type transferUsecaseImpl struct {
	DB                     *gorm.DB
	Validate               *validator.Validate
	Onesignal              *onesignal.OneSignal
	TitanMail              mail.EmailSender
	TransferRepository     repository.TransferRepository
	AccountRepository      repository.AccountRepository
	EntryRepository        repository.EntryRepository
	NotificationRepository repository.NotificationRepository
}

func NewTransferUsecase(db *gorm.DB, validate *validator.Validate, signal *onesignal.OneSignal, titanMail mail.EmailSender, transferRepository repository.TransferRepository,
	accountRepository repository.AccountRepository, entryRepository repository.EntryRepository, notificationRepository repository.NotificationRepository) TransferUsecase {
	return &transferUsecaseImpl{
		DB:                     db,
		Validate:               validate,
		Onesignal:              signal,
		TitanMail:              titanMail,
		TransferRepository:     transferRepository,
		AccountRepository:      accountRepository,
		EntryRepository:        entryRepository,
		NotificationRepository: notificationRepository,
	}
}

func (t *transferUsecaseImpl) TransferMoney(ctx context.Context, requestData *request.TransferRequest, userID string) (*response.TransferResponse, error) {
	tx := t.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := t.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	// Find account by id
	account, err := t.AccountRepository.FindByID(tx, requestData.FromAccountID)
	if err != nil {
		log.Printf("Error when find to account with id : %+v", err)
		return nil, util.AccounDoesNotExist
	}

	// Check if account not belong to user
	if account.UserID != userID {
		log.Printf("Account not belong to user")
		return nil, util.AccountNotBelongToUser
	}

	// Get to account id email
	toAccount, err := t.AccountRepository.FindByID(tx, requestData.ToAccountID)
	if err != nil {
		log.Printf("Error when find to account with id : %+v", err)
		return nil, util.DestinationAccountNotExist
	}

	// Check pin before transaction
	isValid := util.CheckPassword(requestData.Pin, account.User.HashedPin)

	if !isValid {
		log.Printf("Pin not valid")
		return nil, util.InvalidPin
	}

	log.Printf("account user: %+v", account.User)
	if err != nil {
		log.Printf("Error when find account by id : %+v", err)
		return nil, err
	}

	// Check if account balance is sufficient
	if account.Balance < requestData.Amount {
		log.Printf("Insufficient balance")
		return nil, util.InsufficientBalance
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
		return nil, err
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
		return nil, err
	}

	// Updates balance account and send notification, onesignal and email
	if requestData.FromAccountID < requestData.ToAccountID {
		err = t.updateAccountBalance(tx, account.User.Email, toAccount.User.Email, requestData.FromAccountID, -requestData.Amount, requestData.ToAccountID, requestData.Amount)
		if err != nil {
			log.Printf("Error when update balance : %+v", err)
			return nil, err
		}
		// Send push notif onsignal
		err = t.createAndSendNotification(tx, account.User.ID, toAccount.User.ID, requestData.Amount, requestData.Amount)
		if err != nil {
			log.Printf("Error when create notification : %+v", err)
			return nil, err
		}

	} else {
		err = t.updateAccountBalance(tx, account.User.Email, toAccount.User.Email, requestData.ToAccountID, requestData.Amount, requestData.FromAccountID, -requestData.Amount)
		if err != nil {
			log.Printf("Error when update balance : %+v", err)
			return nil, err
		}

		// Send push notif onsignal
		err = t.createAndSendNotification(tx, account.User.ID, toAccount.User.ID, requestData.Amount, requestData.Amount)
		if err != nil {
			log.Printf("Error when create notification : %+v", err)
			return nil, err
		}

	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	// data transfer, account mapping to response
	transferResp := response.TransferResponse{
		Transfer:    transferEntity.ToTransfeInfo(),
		FromAccount: account.ToAccountResponse(),
		ToAccount:   toAccount.ToAccountResponse(),
	}

	return &transferResp, nil
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

func (t *transferUsecaseImpl) createAndSendNotification(tx *gorm.DB, fromId, toId string, amount1, amount2 int64) error {

	if amount1 > amount2 {
		tfInEntity := entity.GetNotificationInEntity(toId, amount1)

		err := t.NotificationRepository.NotificationCreate(tx, tfInEntity)
		if err != nil {
			log.Printf("Error when create tf in notification : %+v", err)
			return err
		}

		tfOutEntity := entity.GetNotificationOutEntity(fromId, amount2)

		err = t.NotificationRepository.NotificationCreate(tx, tfOutEntity)
		if err != nil {
			log.Printf("Error when create tf out notification : %+v", err)
			return err
		}

		go func() {
			t.Onesignal.SendNotification(tfInEntity.Title, tfInEntity.Description)
			t.Onesignal.SendNotification(tfOutEntity.Title, tfOutEntity.Description)
		}()
	} else {
		tfInEntity := entity.GetNotificationInEntity(toId, amount2)

		err := t.NotificationRepository.NotificationCreate(tx, tfInEntity)
		if err != nil {
			log.Printf("Error when create tf in notification : %+v", err)
			return err
		}

		tfOutEntity := entity.GetNotificationOutEntity(fromId, amount1)

		err = t.NotificationRepository.NotificationCreate(tx, tfOutEntity)
		if err != nil {
			log.Printf("Error when create tf out notification : %+v", err)
			return err
		}
		go func() {
			t.Onesignal.SendNotification(tfInEntity.Title, tfInEntity.Description)
			t.Onesignal.SendNotification(tfOutEntity.Title, tfOutEntity.Description)
		}()
	}

	return nil
}
