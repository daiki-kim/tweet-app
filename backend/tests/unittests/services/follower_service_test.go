package services

import (
	"errors"
	"testing"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/daiki-kim/tweet-app/backend/tests/unittests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFollowSuccess(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockFollowerRepository{}
	testFollowerService := services.NewFollowerService(mockRepo)

	// フォローモデルを準備
	followerId := uint(1)
	followeeId := uint(2)

	expectedFollower := &models.Follower{
		FollowerID: followerId,
		FolloweeID: followeeId,
	}

	// モックレポジトリを呼び出し
	mockRepo.On("CreateFollower", expectedFollower).Return(expectedFollower, nil)

	follower, err := testFollowerService.Follow(followerId, followeeId)

	assert.NoError(t, err)
	assert.Equal(t, followerId, follower.FollowerID)
	mockRepo.AssertExpectations(t)
}

func TestFollowFail(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockFollowerRepository{}
	testFollowerService := services.NewFollowerService(mockRepo)

	// フォローモデルを準備
	followerId := uint(1)
	followeeId := uint(0)

	wrongFollower := &models.Follower{
		FollowerID: followerId,
		FolloweeID: followeeId,
	}

	// モックレポジトリを呼び出し
	mockRepo.On("CreateFollower", wrongFollower).Return(nil, errors.New("followee_id is no uint"))

	follower, err := testFollowerService.Follow(followerId, followeeId)

	assert.Error(t, err)
	assert.NotEqual(t, wrongFollower, follower)
	mockRepo.AssertExpectations(t)
}
