package services

import (
	"errors"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
)

type IFollowerService interface {
	Follow(followerId, followeeId uint) (*models.Follower, error)
	GetFollower(id uint) (*models.Follower, error)
	GetFollows(followerId uint) ([]*models.Follower, error)
	GetFollowers(followeeId uint) ([]*models.Follower, error)
	DeleteFollower(id uint, user_id uint) error
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

func (s *FollowerService) GetFollower(id uint) (*models.Follower, error) {
	return s.repository.GetFollower(id)
}

func (s *FollowerService) GetFollows(followerId uint) ([]*models.Follower, error) {
	return s.repository.GetFollowees(followerId)
}

func (s *FollowerService) GetFollowers(followeeId uint) ([]*models.Follower, error) {
	return s.repository.GetFollowers(followeeId)
}

func (s *FollowerService) DeleteFollower(id uint, user_id uint) error {
	follower, err := s.repository.GetFollower(id)
	if err != nil {
		return err
	}

	if follower.FollowerID != user_id {
		return errors.New("you don't have permission to delete this follower")
	}

	return s.repository.DeleteFollower(id)
}
