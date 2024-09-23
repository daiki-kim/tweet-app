package mocks

import (
	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/stretchr/testify/mock"
)

type MockFollowerService struct {
	mock.Mock
}

func (m *MockFollowerService) Follow(followerId, followeeId uint) (*models.Follower, error) {
	args := m.Called(followerId, followeeId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Follower), args.Error(1)
}

func (m *MockFollowerService) GetFollower(id uint) (*models.Follower, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Follower), args.Error(1)
}

func (m *MockFollowerService) GetFollows(followerId uint) ([]*models.Follower, error) {
	args := m.Called(followerId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Follower), args.Error(1)
}

func (m *MockFollowerService) GetFollowers(followeeId uint) ([]*models.Follower, error) {
	args := m.Called(followeeId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Follower), args.Error(1)
}

func (m *MockFollowerService) DeleteFollower(id uint, user_id uint) error {
	args := m.Called(id, user_id)
	return args.Error(0)
}
