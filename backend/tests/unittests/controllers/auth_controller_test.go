package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daiki-kim/tweet-app/backend/apps/controllers"
	"github.com/daiki-kim/tweet-app/backend/apps/dtos"
	"github.com/daiki-kim/tweet-app/backend/tests/unittests/mocks"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignupUsingOAuthSuccess(t *testing.T) {
	// モックレポジトリを準備
	mockAuthService := &mocks.MockAuthService{}
	testAuthController := controllers.NewAuthController(mockAuthService)

	// ginエンジンの設定
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// セッションストアを設定
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("my_session", store))

	// OAuthからのユーザーデータを準備
	userData := dtos.OAuthSignupInput{
		Name:  "testuser",
		Email: "test@example.com",
	}
	SessionData, _ := json.Marshal(userData)

	// OAuthからのユーザーデータをセッションに保存
	r.POST("/api/v1/auth/signup/oauth", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("user_data", string(SessionData))
		session.Save()
		testAuthController.SignupUsingOAuth(c)
	})

	// リクエスト作成
	reqBody := []byte(`{"dob": "2020-01-01"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/signup/oauth", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// mockAuthServiceのmockメソッドを準備
	userData.Dob = "2020-01-01"
	mockAuthService.On("SignupUsingOAuth", userData.Name, userData.Email, userData.Dob).Return(nil)

	// テスト実行
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusCreated, w.Code)
	mockAuthService.AssertExpectations(t)
}
