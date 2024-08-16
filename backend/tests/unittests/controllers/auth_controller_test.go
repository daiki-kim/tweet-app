package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daiki-kim/tweet-app/backend/apps/controllers"
	"github.com/daiki-kim/tweet-app/backend/apps/dtos"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
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
	r.POST("/api/v1/auth/signup/oauth", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Set("user_data", string(SessionData))
		session.Save()
		testAuthController.SignupUsingOAuth(ctx)
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

func TestLoginUsingOAuthSuccess(t *testing.T) {
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
	userData := dtos.OAuthLoginInput{
		Email: "test@example.com",
	}
	SessionData, _ := json.Marshal(userData)

	// OAuthからのユーザーデータをセッションに保存
	r.GET("/api/v1/auth/login/oauth", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Set("user_data", string(SessionData))
		session.Save()
		testAuthController.LoginUsingOAuth(ctx)
	})

	// login responseを準備
	loginResponse := &services.LoginResponse{
		Token:        "test_token",
		RefreshToken: "test_refresh_token",
	}

	// mockAuthServiceのmockメソッドを準備
	mockAuthService.On("LoginUsingOAuth", userData.Email).Return(loginResponse, nil)

	// テスト実行
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/login/oauth", nil)
	r.ServeHTTP(w, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, w.Code)
	mockAuthService.AssertExpectations(t)
}

func TestLoginUsingOAuthNotUserFound(t *testing.T) {
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
	userData := dtos.OAuthLoginInput{
		Email: "test@example.com",
	}
	SessionData, _ := json.Marshal(userData)

	// OAuthからのユーザーデータをセッションに保存
	r.GET("/api/v1/auth/login/oauth", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Set("user_data", string(SessionData))
		session.Save()
		testAuthController.LoginUsingOAuth(ctx)
	})

	// mockAuthServiceのmockメソッドを準備
	mockAuthService.On("LoginUsingOAuth", userData.Email).Return(nil, errors.New("user not found"))

	// テスト実行
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/login/oauth", nil)
	r.ServeHTTP(w, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockAuthService.AssertExpectations(t)
}

func TestLoginSuccess(t *testing.T) {
	// モックレポジトリを準備
	mockAuthService := &mocks.MockAuthService{}
	testAuthController := controllers.NewAuthController(mockAuthService)

	// ginエンジンの設定
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// セッションストアを設定
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("my_session", store))

	// テストユーザーデータを準備
	userData := dtos.LoginInput{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Login APIを準備
	r.POST("/api/v1/auth/login", testAuthController.Login)

	// リクエスト作成
	reqBody := []byte(`{"email": "test@example.com", "password": "testpassword"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// login responseを準備
	loginResponse := &services.LoginResponse{
		Token:        "test_token",
		RefreshToken: "test_refresh_token",
	}

	// mockAuthServiceのmockメソッドを準備
	mockAuthService.On("Login", userData.Email, userData.Password).Return(loginResponse, nil)

	// テスト実行
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, w.Code)
	mockAuthService.AssertExpectations(t)
}

func TestLoginUserNotFound(t *testing.T) {
	// モックレポジトリを準備
	mockAuthService := &mocks.MockAuthService{}
	testAuthController := controllers.NewAuthController(mockAuthService)

	// ginエンジンの設定
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// セッションストアを設定
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("my_session", store))

	// テストユーザーデータを準備
	wrongUserData := dtos.LoginInput{
		Email:    "wronguser@example.com",
		Password: "testpassword",
	}

	// Login APIを準備
	r.POST("/api/v1/auth/login", testAuthController.Login)

	// リクエスト作成
	reqBody := []byte(`{"email": "wronguser@example.com", "password": "testpassword"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// mockAuthServiceのmockメソッドを準備
	mockAuthService.On("Login", wrongUserData.Email, wrongUserData.Password).Return(nil, errors.New("user not found"))

	// テスト実行
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockAuthService.AssertExpectations(t)
}

func TestLoginUserInvalidPassword(t *testing.T) {
	// モックレポジトリを準備
	mockAuthService := &mocks.MockAuthService{}
	testAuthController := controllers.NewAuthController(mockAuthService)

	// ginエンジンの設定
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// セッションストアを設定
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("my_session", store))

	// テストユーザーデータを準備
	wrongUserData := dtos.LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Login APIを準備
	r.POST("/api/v1/auth/login", testAuthController.Login)

	// リクエスト作成
	reqBody := []byte(`{"email": "test@example.com", "password": "wrongpassword"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// mockAuthServiceのmockメソッドを準備
	mockAuthService.On("Login", wrongUserData.Email, wrongUserData.Password).Return(nil, errors.New("invalid password"))

	// テスト実行
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuthService.AssertExpectations(t)
}
