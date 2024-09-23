package services

import (
	"errors"
	"testing"
	"time"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/daiki-kim/tweet-app/backend/tests/unittests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFollowSuccess(t *testing.T) {
	// モックレポジトリを準備
	mockRepo, testFollowerService := prepareTestFollowerService()

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
	mockRepo, testFollowerService := prepareTestFollowerService()

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

func TestGetFollowersNotFound(t *testing.T) {
	// モックレポジトリを準備
	mockRepo, testFollowerService := prepareTestFollowerService()

	// モックレポジトリを呼び出し
	mockRepo.On("GetFollowers", uint(1)).Return(nil, errors.New("followers not found"))

	followersUserData, err := testFollowerService.GetFollowers(1)

	assert.Error(t, err)
	assert.Nil(t, followersUserData)
	assert.Equal(t, "followers not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetFollowersSuccess(t *testing.T) {
	// モックレポジトリを準備
	mockRepo, testFollowerService := prepareTestFollowerService()

	// フォロワーのユーザーデータを準備
	testuser1 := &models.User{
		Name:     "testuser1",
		Email:    "test1@example.com",
		Password: "testpassword",
		Dob:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	testuser3 := &models.User{
		Name:     "testuser3",
		Email:    "test3@example.com",
		Password: "testpassword",
		Dob:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// フォローモデルを準備
	testFollower1Follows2 := &models.Follower{
		FollowerID: 1,
		FolloweeID: 2,
		Follower:   testuser1,
		Followee:   nil,
	}
	testFollower3Follows2 := &models.Follower{
		FollowerID: 3,
		FolloweeID: 2,
		Follower:   testuser3,
		Followee:   nil,
	}

	// モックレポジトリを呼び出し
	mockRepo.On("GetFollowers", uint(2)).Return([]*models.Follower{testFollower1Follows2, testFollower3Follows2}, nil)

	followersUserData, err := testFollowerService.GetFollowers(2)

	assert.NoError(t, err)
	assert.Equal(t, testuser1, followersUserData[0])
	assert.Equal(t, testuser3, followersUserData[1])
}

func prepareTestFollowerService() (*mocks.MockFollowerRepository, services.IFollowerService) {
	mockRepo := &mocks.MockFollowerRepository{}
	testFollowerService := services.NewFollowerService(mockRepo)
	return mockRepo, testFollowerService
}
