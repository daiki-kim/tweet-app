package repositories

import (
	"errors"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"gorm.io/gorm"
)

type IFollowerRepository interface {
	CreateFollower(follower *models.Follower) (*models.Follower, error)
	GetFollower(id uint) (*models.Follower, error)
	GetFollowers(followeeId uint) ([]*models.Follower, error)
	DeleteFollower(id uint) error
}

type FollowerRepository struct {
	DB *gorm.DB
}

func NewFollowerRepository(db *gorm.DB) IFollowerRepository {
	return &FollowerRepository{DB: db}
}

func (r *FollowerRepository) CreateFollower(follower *models.Follower) (*models.Follower, error) {
	result := r.DB.Create(follower)
	if result.Error != nil {
		return nil, result.Error
	}

	return follower, nil
}

func (r *FollowerRepository) GetFollower(id uint) (*models.Follower, error) {
	var follower models.Follower
	result := r.DB.First(&follower, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("follower not found")
	}

	return &follower, nil
}

func (r *FollowerRepository) GetFollowers(followeeId uint) ([]*models.Follower, error) {
	var followers []*models.Follower
	result := r.DB.Preload("Follower").Where("followee_id = ?", followeeId).Find(&followers)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("followers not found")
	} else if result.Error != nil {
		return nil, result.Error
	}

	return followers, nil
}

func (r *FollowerRepository) DeleteFollower(id uint) error {
	result := r.DB.Delete(&models.Follower{}, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("follower not found")
	} else if result.Error != nil {
		return result.Error
	}

	return nil
}
