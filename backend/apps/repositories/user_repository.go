package repositories

import (
	"errors"
	"log"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		log.Println("failed to create user: ", result.Error)
		return result.Error
	}

	return nil
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	result := r.db.First(user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("user not found: ", result.Error)
		return nil, errors.New("user not found")
	}

	if result.Error != nil {
		log.Println("failed to find user: ", result.Error)
		return nil, result.Error
	}

	return user, nil
}
