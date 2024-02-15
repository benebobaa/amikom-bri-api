package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/delivery/http/dto/request"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"gorm.io/gorm"
	"log"
)

type UserRepository interface {
	UserCreate(tx *gorm.DB, value *entity.User) error
	FindUsernameIsExists(tx *gorm.DB, username string) (*entity.User, bool, error)
	FindByEmailVerified(tx *gorm.DB, email string) (*entity.User, error)
	FindByUsernameOrEmail(tx *gorm.DB, value *entity.User) (*entity.User, error)
	UpdateUser(tx *gorm.DB, value *entity.User) error
	DeleteUser(tx *gorm.DB, value *entity.User) error
	FindByUserID(tx *gorm.DB, userID string) (*entity.User, error)
	FindAllUsers(db *gorm.DB, request *request.SearchPaginationRequest) ([]entity.User, int64, error)
	FilterUser(request *request.SearchPaginationRequest) func(tx *gorm.DB) *gorm.DB
}

type userRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{}
}

func (u *userRepositoryImpl) UserCreate(tx *gorm.DB, value *entity.User) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create user : %v", result.Error))
		return result.Error
	}

	return nil
}

func (u *userRepositoryImpl) FindUsernameIsExists(tx *gorm.DB, username string) (*entity.User, bool, error) {
	var value entity.User

	result := tx.Where("username = ?", username).Preload("Account").First(&value)

	log.Println("result", &value)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find user with username : %v", result.Error))
		return nil, false, result.Error
	}

	if value.Username == username {
		return &value, true, nil
	}

	return &value, false, nil
}

func (u *userRepositoryImpl) FindByEmailVerified(tx *gorm.DB, email string) (*entity.User, error) {
	var value entity.User

	result := tx.Where("email = ?", email).Where("is_email_verified = true").First(&value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find user with email : %v", result.Error))
		return nil, result.Error
	}

	return &value, nil
}

func (u *userRepositoryImpl) FindByUsernameOrEmail(tx *gorm.DB, value *entity.User) (*entity.User, error) {
	var user entity.User

	result := tx.Preload("Account").Where("username = ?", value.Username).Or("email = ?", value.Email).First(&user)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find username or email : %v", result.Error))
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepositoryImpl) UpdateUser(tx *gorm.DB, value *entity.User) error {
	result := tx.Model(value).Where("id = ?", value.ID).Updates(value).Preload("Account").First(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when update user : %v", result.Error))
		return result.Error
	}

	return nil
}

func (u *userRepositoryImpl) DeleteUser(tx *gorm.DB, value *entity.User) error {
	result := tx.Delete(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when delete user : %v", result.Error))
		return result.Error
	}

	return nil
}

func (u *userRepositoryImpl) FindByUserID(tx *gorm.DB, userID string) (*entity.User, error) {
	var user entity.User

	result := tx.Preload("Account").Where("id = ?", userID).First(&user)

	log.Println("USER WITH ACCOUNT", user)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find user by user id : %v", result.Error))
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepositoryImpl) FindAllUsers(db *gorm.DB, request *request.SearchPaginationRequest) ([]entity.User, int64, error) {
	var users []entity.User

	if err := db.Scopes(u.FilterUser(request)).Preload("Account").Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.User{}).Scopes(u.FilterUser(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (p *userRepositoryImpl) FilterUser(request *request.SearchPaginationRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if keyword := request.Keyword; keyword != "" {
			keyword = "%" + keyword + "%"
			tx = tx.Where("full_name LIKE ?", keyword)
		}

		return tx
	}
}
