package repository

import (
	"log"

	"github.com/YungBenn/tech-shop-microservices/internal/auth/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindAllUsers() ([]entity.User, error)
	FindUserByID(id string) (entity.User, error)
	FindUserByEmail(email string) (entity.User, error)
	SaveUser(user entity.User) (*entity.User, error)
	DeleteUserByID(user entity.User) (*entity.User, error)
}

type AuthRepositoryImpl struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewAuthRepository(db *gorm.DB, log *logrus.Logger) AuthRepository {
	return &AuthRepositoryImpl{db, log}
}

func (u *AuthRepositoryImpl) FindAllUsers() ([]entity.User, error) {
	var users []entity.User

	err := u.db.Model([]entity.User{}).Find(&users).Error
	return users, err
}

func (u *AuthRepositoryImpl) FindUserByID(id string) (entity.User, error) {
	var user = entity.User{ID: id}

	result := u.db.Model([]entity.User{}).First(&user)
	if result.RowsAffected == 0 {
		log.Println("Error")
	}

	return user, nil
}

func (u *AuthRepositoryImpl) FindUserByEmail(email string) (entity.User, error) {
	var user entity.User

	result := u.db.Where(&entity.User{Email: email}).First(&user)
	if result.RowsAffected == 0 {
		u.log.Error("User not found")
	}

	return user, nil
}

func (u *AuthRepositoryImpl) SaveUser(user entity.User) (*entity.User, error) {
	result := u.db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *AuthRepositoryImpl) DeleteUserByID(user entity.User) (*entity.User, error) {
	result := u.db.Delete(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
