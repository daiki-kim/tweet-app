package repositories

import (
	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"gorm.io/gorm"
)

type IFollowerRepository interface {
	CreateFollower(follower *models.Follower) (*models.Follower, error)
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
