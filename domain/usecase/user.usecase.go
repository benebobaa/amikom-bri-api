package usecase

import (
	"context"
	"errors"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/response"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/benebobaa/amikom-bri-api/domain/repository"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/token"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"math"
	"time"
)

type UserUsecase interface {
	RegisterNewUser(ctx context.Context, requestData *request.UserRegisterRequest) (*response.UserResponse, error)
	VerifyUserEmail(ctx context.Context, secretCode string) (*response.EmailVerifyResponse, error)
	SoftDeleteUser(ctx context.Context, requestUsername, payloadUsername string) error
	ProfileUser(ctx context.Context, userID string) (*response.UserProfileResponse, error)
	GetAllUsers(ctx context.Context, requestData *request.SearchPaginationRequest) (*response.UserResponses, error)
}

type userUsecaseImpl struct {
	DB                *gorm.DB
	Validate          *validator.Validate
	ViperConfig       util.Config
	TokenMaker        token.Maker
	UserRepository    repository.UserRepository
	EmailRepository   repository.EmailRepository
	AccountRepository repository.AccountRepository
}

func NewUserUsecase(db *gorm.DB, validate *validator.Validate, userRepository repository.UserRepository,
	emailRepository repository.EmailRepository, accountRepository repository.AccountRepository) UserUsecase {
	return &userUsecaseImpl{
		DB:                db,
		Validate:          validate,
		UserRepository:    userRepository,
		EmailRepository:   emailRepository,
		AccountRepository: accountRepository,
	}

}

func (u *userUsecaseImpl) RegisterNewUser(ctx context.Context, requestData *request.UserRegisterRequest) (*response.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	_, exists, _ := u.UserRepository.FindUsernameIsExists(tx, requestData.Username)

	if exists {
		return nil, util.UsernameAlreadyExists
	}

	resultUser, _ := u.UserRepository.FindByEmail(tx, requestData.Email)

	log.Printf("resultUser : %+v", resultUser)
	if resultUser != nil && resultUser.IsEmailVerified {
		return nil, util.EmailAlreadyExists
	}

	hashedPassword, _ := util.HashPassword(requestData.Password)
	requestData.Password = hashedPassword

	requestUserEntity := requestData.ToEntity()

	err = u.UserRepository.UserCreate(tx, requestUserEntity)

	if err != nil {
		log.Printf("Failed create user : %+v", err)
		return nil, err
	}

	requestEmail := request.EmailRequest{
		UserID:     requestUserEntity.ID,
		Username:   requestUserEntity.Username,
		Email:      requestUserEntity.Email,
		SecretCode: util.RandomCombineIntAndString() + util.RandomCombineIntAndString(),
		ExpiredAt:  time.Now().Add(time.Minute * 5),
	}

	err = u.EmailRepository.EmailVerifyCreate(tx, requestEmail.ToEntity())

	if err != nil {
		log.Printf("Failed create email verify : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Printf("Failed commit db transaction : %+v", err)
		return nil, err
	}

	return requestUserEntity.ToUserResponse(), nil
}

func (u *userUsecaseImpl) VerifyUserEmail(ctx context.Context, secretCode string) (*response.EmailVerifyResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	log.Printf("secretCode : %+v", secretCode)
	err := u.Validate.Var(secretCode, "required")
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return nil, err
	}

	resultEmail, err := u.EmailRepository.FindEmailBySecretCode(tx, secretCode)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Email not found : %+v", err)
			return nil, util.EmailVerifyCodeNotValid
		}
		log.Printf("Failed find email by secret code : %+v", err)
		return nil, err
	}

	if resultEmail.IsUsed {
		return nil, util.EmailVerifyAlreadyUsed
	}

	if resultEmail.ExpiredAt.Before(time.Now()) {
		log.Printf("Email verify expired")
		return nil, util.EmailVerifyExpired
	}

	resultEmail.IsUsed = true
	err = u.EmailRepository.UpdateEmailVerify(tx, resultEmail)
	if err != nil {
		log.Printf("Failed update email verify : %+v", err)
		return nil, err
	}

	userRequestEntity := &entity.User{ID: resultEmail.UserID, IsEmailVerified: true}

	_, isExists, _ := u.AccountRepository.FindByUserID(tx, resultEmail.UserID)

	if isExists {
		log.Printf("Account already exists")
		return nil, util.AccountAlreadyExists
	}

	err = u.UserRepository.UpdateUser(tx, userRequestEntity)

	if err != nil {
		log.Printf("Failed update user verified email : %+v", err)
		return nil, err
	}

	u.AccountRepository.AccountCreate(tx, &entity.Account{UserID: resultEmail.UserID})

	err = tx.Commit().Error
	if err != nil {
		log.Printf("Failed commit db transaction : %+v", err)
		return nil, err
	}

	return resultEmail.ToEmailVerifyResponse(), nil
}

func (u *userUsecaseImpl) SoftDeleteUser(ctx context.Context, requestUsername, payloadUsername string) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Var(requestUsername, "required")
	if err != nil {
		log.Printf("Invalid request body : %+v", err)
		return err
	}

	if requestUsername != payloadUsername {
		return util.UnauthorizedDeleteUser
	}

	resultUser, _, err := u.UserRepository.FindUsernameIsExists(tx, requestUsername)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Username not found : %+v", err)
			return util.UsernameNotFound
		}
		log.Printf("Failed find username : %+v", err)
		return err
	}

	err = u.UserRepository.DeleteUser(tx, resultUser)

	if err != nil {
		log.Printf("Failed delete user : %+v", err)
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Printf("Failed commit db transaction : %+v", err)
		return err
	}

	return nil
}

func (u *userUsecaseImpl) ProfileUser(ctx context.Context, userID string) (*response.UserProfileResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Var(userID, "required")
	if err != nil {
		log.Printf("Invalid request : %+v", err)
		return nil, err
	}

	resultUser, err := u.UserRepository.FindByUserID(tx, userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found : %+v", err)
			return nil, util.UserNotFound
		}
		log.Printf("Failed find user : %+v", err)
		return nil, err
	}

	return resultUser.ToUserProfileResponse(), nil
}

func (u *userUsecaseImpl) GetAllUsers(ctx context.Context, requestData *request.SearchPaginationRequest) (*response.UserResponses, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	users, total, err := u.UserRepository.FindAllUsers(tx, requestData)

	if err != nil {
		log.Printf("Failed find all users : %+v", err)
		return nil, err
	}

	resultPaging := &response.PostPageMetaData{
		Page:      requestData.Page,
		Size:      requestData.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(requestData.Size))),
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return entity.ToUserResponses(users, resultPaging), nil
}
