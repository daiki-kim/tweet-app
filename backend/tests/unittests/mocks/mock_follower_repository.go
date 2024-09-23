package mocks

import (
	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/stretchr/testify/mock"
)

type MockFollowerRepository struct {
	mock.Mock
}

func (m *MockFollowerRepository) CreateFollower(follower *models.Follower) (*models.Follower, error) {
	args := m.Called(follower)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Follower), args.Error(1)
}

func (m *MockFollowerRepository) GetFollowers(followeeId uint) ([]*models.Follower, error) {
	args := m.Called(followeeId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Follower), args.Error(1)
}
