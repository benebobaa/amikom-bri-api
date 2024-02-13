package usecase

import (
	"context"
	"errors"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/benebobaa/amikom-bri-api/domain/repository"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/mail"
	"github.com/benebobaa/amikom-bri-api/util/token"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type ForgotPasswordUsecase interface {
	ForgotPasswordRequest(ctx context.Context, requestData *request.ForgotPasswordRequest, baseUrl string) error
	ResetPasswordRequest(ctx context.Context, requestData *request.ResetPasswordRequest, secretToken string) error
	CheckResetTokenIsUsed(ctx context.Context, secretToken string) error
}

type forgotPasswordImpl struct {
	DB                       *gorm.DB
	Validate                 *validator.Validate
	ViperConfig              util.Config
	TokenMaker               token.Maker
	TitanMail                mail.EmailSender
	UserRepository           repository.UserRepository
	ForgotPasswordRepository repository.ForgotPasswordRepository
}

func NewForgotPasswordUsecase(db *gorm.DB, validate *validator.Validate, viperConfig util.Config,
	tokenMaker token.Maker, titanMail mail.EmailSender, userRepository repository.UserRepository, forgotPasswordRepository repository.ForgotPasswordRepository) ForgotPasswordUsecase {
	return &forgotPasswordImpl{
		DB:                       db,
		Validate:                 validate,
		ViperConfig:              viperConfig,
		TokenMaker:               tokenMaker,
		TitanMail:                titanMail,
		UserRepository:           userRepository,
		ForgotPasswordRepository: forgotPasswordRepository,
	}
}

func (f *forgotPasswordImpl) ForgotPasswordRequest(ctx context.Context, requestData *request.ForgotPasswordRequest, baseUrl string) error {
	tx := f.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := f.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return err
	}

	// Get user data by email
	emailUser, err := f.UserRepository.FindByEmailVerified(tx, requestData.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found with email : %+v", err)
			return util.EmailNotExists
		}
		log.Printf("Failed find user with email : %+v", err)
		return err
	}

	if emailUser == nil {
		return util.UsernameAndEmailNotMatch
	}

	// Compare username from the findByEmail and the username from the request
	if emailUser.Username != requestData.Username {
		return util.UsernameAndEmailNotMatch
	}

	// Compare the email from the findByUsername and the email from the request
	if emailUser.Email != requestData.Email {
		return util.UsernameAndEmailNotMatch
	}

	// Create Reset Token
	resetToken, _, err := f.TokenMaker.CreateToken(token.UserPayload{
		UserID:   emailUser.ID,
		Username: emailUser.Username,
	}, f.ViperConfig.SecretKeyResetPassword, f.ViperConfig.SecretKeyDuration)

	if err != nil {
		log.Printf("Failed create reset token : %+v", err)
		return err
	}

	requestFpEntity := &entity.ForgotPassword{
		UserID:     emailUser.ID,
		ResetToken: resetToken,
		IsUsed:     false,
	}

	err = f.ForgotPasswordRepository.ForgotPasswordCreate(tx, requestFpEntity)

	if err != nil {
		log.Printf("Failed create forgot password : %+v", err)
		return err
	}

	go func() {
		resetLink := baseUrl + "/users/reset-password?secret=" + resetToken
		subject, content, toEmail := mail.GetSenderParamResetPassword(emailUser.Email, resetLink)
		err := f.TitanMail.SendEmail(subject, content, toEmail, []string{}, []string{}, []string{})
		if err != nil {
			log.Printf("Failed send email : %+v", err)
		}
	}()

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (f *forgotPasswordImpl) ResetPasswordRequest(ctx context.Context, requestData *request.ResetPasswordRequest, secretToken string) error {
	tx := f.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := f.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil
	}

	// Verify reset token
	payload, err := f.TokenMaker.VerifyToken(secretToken, f.ViperConfig.SecretKeyResetPassword)

	if err != nil {
		log.Printf("Failed verify token : %+v", err)
		return err
	}

	userForgotPass, err := f.ForgotPasswordRepository.FindByUserID(tx, payload.UserID)

	if err != nil {
		log.Printf("Failed find by user id : %+v", err)
		return err
	}

	// Check reset token equals on db
	if userForgotPass.ResetToken != secretToken {
		log.Printf("Invalid reset token")
		return util.InvalidResetToken
	}

	// Check user id payload equals on db
	if userForgotPass.UserID != payload.UserID {
		log.Printf("Invalid user id")
		return util.InvalidResetToken
	}

	// Check if reset token already used
	if userForgotPass.IsUsed {
		log.Printf("Reset token already used")
		return util.ResetTokenAlreadyUsed
	}

	// Hash new password
	password, err := util.HashPassword(requestData.NewPassword)
	if err != nil {
		log.Printf("Failed hash password : %+v", err)
		return err
	}

	// Update new password user
	err = f.UserRepository.UpdateUser(tx, &entity.User{ID: userForgotPass.UserID, HashedPassword: password})

	if err != nil {
		log.Printf("Failed update user : %+v", err)
		return err
	}

	// Update reset token to be used
	err = f.ForgotPasswordRepository.UpdateForgotPassword(tx, &entity.ForgotPassword{ID: userForgotPass.ID, UserID: userForgotPass.UserID, IsUsed: true})
	if err != nil {
		log.Printf("Failed update forgot password : %+v", err)
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (f *forgotPasswordImpl) CheckResetTokenIsUsed(ctx context.Context, secretToken string) error {
	tx := f.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	payload, err := f.TokenMaker.VerifyToken(secretToken, f.ViperConfig.SecretKeyResetPassword)

	if err != nil {
		log.Printf("Failed verify token : %+v", err)
		return err
	}

	userForgotPass, err := f.ForgotPasswordRepository.FindByUserID(tx, payload.UserID)

	if err != nil {
		log.Printf("Failed find by user id : %+v", err)
		return err
	}

	if userForgotPass.IsUsed {
		log.Printf("Reset token already used")
		return util.ResetTokenAlreadyUsed
	}

	return nil
}
