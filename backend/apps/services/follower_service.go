package services

import (
	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
)

type IFollowerService interface {
	Follow(followerId, followeeId uint) (*models.Follower, error)
}

type FollowerService struct {
	repository repositories.IFollowerRepository
}

func NewFollowerService(repository repositories.IFollowerRepository) IFollowerService {
	return &FollowerService{repository: repository}
}

func (s *FollowerService) Follow(followerId, followeeId uint) (*models.Follower, error) {
	follower := &models.Follower{
		FollowerID: followerId,
		FolloweeID: followeeId,
	}

	return s.repository.CreateFollower(follower)
}
