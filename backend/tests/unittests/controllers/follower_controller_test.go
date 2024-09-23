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
	mockFollowerService := &mocks.MockFollowerService{}
	testFollowerController := controllers.NewFollowerController(mockFollowerService)

	// ginエンジンの設定
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Follow APIを準備
	r.POST("/api/v1/follow", func(c *gin.Context) {
		// テストのために context に followerId を設定
		c.Set("user_id", "1")
		testFollowerController.Follow(c)
	})

	// リクエスト作成
	reqBody := []byte(`{"followee_id": 2}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/follow", bytes.NewBuffer(reqBody))
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
