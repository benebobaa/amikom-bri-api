package repository

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"gorm.io/gorm"
	"log"
)

type UserRepository interface {
	UserCreate(tx *gorm.DB, value *entity.User) error
	FindUsernameIsExists(tx *gorm.DB, username string) (bool, error)
	FindByEmail(tx *gorm.DB, email string) (*entity.User, error)
	FindByUsernameOrEmail(tx *gorm.DB, value *entity.User) (*entity.User, error)
	UpdateUser(tx *gorm.DB, value *entity.User) error
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

func (u *userRepositoryImpl) FindUsernameIsExists(tx *gorm.DB, username string) (bool, error) {
	var value entity.User

	result := tx.Where("username = ?", username).First(&value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find user with username : %v", result.Error))
		return false, result.Error
	}

	if value.Username == username {
		return true, nil
	}

	return false, nil
}

func (u *userRepositoryImpl) FindByEmail(tx *gorm.DB, email string) (*entity.User, error) {
	var value entity.User

	result := tx.Where("email = ?", value.Email).First(&value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find user with email : %v", result.Error))
		return nil, result.Error
	}

	return &value, nil
}

func (u *userRepositoryImpl) FindByUsernameOrEmail(tx *gorm.DB, value *entity.User) (*entity.User, error) {
	var user entity.User

	result := tx.Where("username = ?", value.Username).Or("email = ?", value.Email).First(&user)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find username or email : %v", result.Error))
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepositoryImpl) UpdateUser(tx *gorm.DB, value *entity.User) error {
	result := tx.Model(value).Updates(value).First(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when update user : %v", result.Error))
		return result.Error
	}

	return nil
}
