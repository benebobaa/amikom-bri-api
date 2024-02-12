package usecase

import (
	"context"
	"errors"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"github.com/benebobaa/amikom-bri-api/domain/repository"
	"github.com/benebobaa/amikom-bri-api/util"
	token2 "github.com/benebobaa/amikom-bri-api/util/token"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type LoginUseCase interface {
	LoginUser(ctx context.Context, requestData *request.LoginRequest, userAgent, clientIP string) (*response.LoginResponse, error)
}

type loginUseCaseImpl struct {
	DB                *gorm.DB
	Validate          *validator.Validate
	ViperConfig       util.Config
	TokenMaker        token2.Maker
	UserRepository    repository.UserRepository
	SessionRepository repository.SessionRepository
}

func NewLoginUseCase(db *gorm.DB, validate *validator.Validate, tokenMaker token2.Maker,
	viperConfig util.Config, userRepository repository.UserRepository, sessionRepository repository.SessionRepository) LoginUseCase {
	return &loginUseCaseImpl{
		DB:                db,
		Validate:          validate,
		ViperConfig:       viperConfig,
		TokenMaker:        tokenMaker,
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
	}
}

func (l *loginUseCaseImpl) LoginUser(ctx context.Context, requestData *request.LoginRequest, userAgent, clientIP string) (*response.LoginResponse, error) {
	tx := l.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := l.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	requestUserEntity := requestData.ToUserEntity()

	resultUser, err := l.UserRepository.FindByUsernameOrEmail(tx, requestUserEntity)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Username or email not found : %+v", err)
			return nil, util.UsernameOrEmailNotFound
		}
		log.Printf("Failed find username or email : %+v", err)
		return nil, err
	}

	if !resultUser.IsEmailVerified {
		log.Printf("Email not verified")
		return nil, util.EmailNotVerified
	}

	isValid := util.CheckPassword(requestData.Password, resultUser.HashedPassword)
	if !isValid {
		log.Printf("Password not valid")
		return nil, util.InvalidPassword
	}

	accessToken, accessPayload, err := l.TokenMaker.CreateToken(token2.UserPayload{
		Username: resultUser.Username,
	}, l.ViperConfig.TokenAccessSymetricKey, l.ViperConfig.TokenAccessDuration)

	if err != nil {
		log.Printf("Failed create access token : %+v", err)
		return nil, err
	}

	refreshToken, refreshPayload, err := l.TokenMaker.CreateToken(token2.UserPayload{
		Username: resultUser.Username,
	}, l.ViperConfig.TokenRefreshSymetricKey, l.ViperConfig.RefreshTokenDuration)

	if err != nil {
		log.Printf("Failed create refresh token : %+v", err)
		return nil, err
	}

	sessionRequest := &request.SessionRequest{
		ID:           refreshPayload.ID.String(),
		Username:     resultUser.Username,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	}

	sessionEntity := sessionRequest.ToEntity()

	err = l.SessionRepository.SessionCreate(tx, sessionEntity)

	if err != nil {
		log.Printf("Failed create session user : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	sessionResponse := &response.SessionsResponse{
		SessionsID:            sessionEntity.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}
	return resultUser.ToLoginResponseWithToken(sessionResponse), nil
}
