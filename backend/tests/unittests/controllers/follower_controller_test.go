package controllers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daiki-kim/tweet-app/backend/apps/controllers"
	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/tests/unittests/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFollowSuccess(t *testing.T) {
	// モックサービスを準備
	mockFollowerService, testFollowerController := prepareTestController()

	// ginエンジンの設定
	r := setupTestRouter()

	// Follow APIを準備
	r.POST("/api/v1/follower", func(c *gin.Context) {
		// テストのために context に followerId を設定
		c.Set("user_id", "1")
		testFollowerController.Follow(c)
	})

	// リクエスト作成
	reqBody := []byte(`{"followee_id": 2}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/follower", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// レスポンスを準備
	w := httptest.NewRecorder()

	// Follow responseを準備
	followerResponse := &models.Follower{
		ID:         1,
		FollowerID: 1,
		FolloweeID: 2,
		Follower:   nil,
		Followee:   nil,
	}

	// モックサービスを準備
	mockFollowerService.On("Follow", uint(1), uint(2)).Return(followerResponse, nil)

	// follower responseを準備
	followerResponseJson := `{
		"id": 1,
		"follower_id": 1,
		"followee_id": 2,
		"follower": null,
		"followee": null
	}`

	// リクエスト実行
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, followerResponseJson, w.Body.String())
	mockFollowerService.AssertExpectations(t)
}

func TestGetFollowersSuccess(t *testing.T) {
	// モックサービスを準備
	mockFollowerService, testFollowerController := prepareTestController()

	// ginエンジンの設定
	r := setupTestRouter()

	// GetFollowers APIを準備
	r.GET("/api/v1/follower/:followee_id", testFollowerController.GetFollowers)

	// リクエスト作成
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/follower/3", nil)

	// レスポンスを準備
	w := httptest.NewRecorder()

	// Follower responseを準備
	followerResponse := []*models.Follower{
		{
			ID:         1,
			FollowerID: 1,
			FolloweeID: 3,
			Follower:   nil,
			Followee:   nil,
		},
		{
			ID:         2,
			FollowerID: 2,
			FolloweeID: 3,
			Follower:   nil,
			Followee:   nil,
		},
	}

	// モックサービスを準備
	mockFollowerService.On("GetFollowers", uint(3)).Return(followerResponse, nil)

	// follower responseを準備
	followerResponseJson := `[
		{
			"id": 1,
			"follower_id": 1,
			"followee_id": 3,
			"follower": null,
			"followee": null
		},
		{
			"id": 2,
			"follower_id": 2,
			"followee_id": 3,
			"follower": null,
			"followee": null
		}
	]`

	// リクエスト実行
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, followerResponseJson, w.Body.String())
	mockFollowerService.AssertExpectations(t)
}

func TestDeleteFollower(t *testing.T) {
	// モックサービスを準備
	mockFollowerService, testFollowerController := prepareTestController()

	// ginエンジンの設定
	r := setupTestRouter()

	// DeleteFollower APIを準備
	r.DELETE("/api/v1/follower/:id", func(c *gin.Context) {
		// テストのために context に followerId を設定
		c.Set("user_id", "1")
		testFollowerController.DeleteFollower(c)
	})

	// リクエスト作成
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/follower/1", nil)

	// レスポンスを準備
	w := httptest.NewRecorder()

	// モックサービスを準備
	mockFollowerService.On("DeleteFollower", uint(1), uint(1)).Return(nil)

	// リクエスト実行
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockFollowerService.AssertExpectations(t)
}

func prepareTestController() (*mocks.MockFollowerService, controllers.IFollowerController) {
	mockFollowerService := &mocks.MockFollowerService{}
	testFollowerController := controllers.NewFollowerController(mockFollowerService)

	return mockFollowerService, testFollowerController
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}
