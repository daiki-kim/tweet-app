package services

import (
	"log"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
)

type IFollowerService interface {
	Follow(followerId, followeeId uint) (*models.Follower, error)
	GetFollowers(followeeId uint) ([]*models.User, error)
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

func (s *FollowerService) GetFollowers(followeeId uint) ([]*models.User, error) {
	followers, err := s.repository.GetFollowers(followeeId)
	if err != nil {
		return nil, err
	}
	log.Println("[GetFollowers] followers: ", followers)

	var followersUserData []*models.User
	for _, follower := range followers {
		log.Println("[GetFollowers] follower: ", follower)
		followersUserData = append(followersUserData, follower.Follower)
	}

	log.Println("[GetFollowers] followersUserData: ", followersUserData)

	return followersUserData, nil
}
