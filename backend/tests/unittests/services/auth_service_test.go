package services_test

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/daiki-kim/tweet-app/backend/tests/unittests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestPrepareBaseUserModel(t *testing.T) {
	// ユーザーモデルを準備
	name := "testuser"
	email := "test@example.com"
	dobString := "2020-01-01"
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	user, err := services.PrepareBaseUserModel(name, email, dobString)

	assert.NoError(t, err)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, dob, user.Dob)
}

func TestSignupUsingOAuthSuccess(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockUserRepository{}
	testAuthService := services.NewAuthService(mockRepo)

	// ユーザーモデルを準備
	name := "testuser"
	email := "test@example.com"
	dobString := "2020-01-01"
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedUser := &models.User{
		Name:  name,
		Email: email,
		Dob:   dob,
	}

	// SignupUsingOAuthで使用するmockメソッドを準備
	mockRepo.On("CreateUser", mock.MatchedBy(func(user *models.User) bool {
		return user.Name == expectedUser.Name && user.Email == expectedUser.Email && user.Dob == expectedUser.Dob
	})).Return(nil)

	// サインアップ
	err := testAuthService.SignupUsingOAuth(name, email, dobString)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSignupUsingOAuthErrorByNonEmail(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	testAuthService := services.NewAuthService(mockRepo)

	// emailが入力されていないユーザーモデルを準備
	name := "testuser"
	email := ""
	dobString := "2020-01-01"
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	wrongUser := &models.User{
		Name:  name,
		Email: email,
		Dob:   dob,
	}

	// SignupUsingOAuthで使用するmockメソッドを準備
	mockRepo.On("CreateUser", wrongUser).Return(errors.New("email is required"))

	// サインアップ
	err := testAuthService.SignupUsingOAuth(name, email, dobString)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSignupSuccess(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockUserRepository{}
	testAuthService := services.NewAuthService(mockRepo)

	// ユーザーモデルを準備
	name := "testuser"
	email := "test@example.com"
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPasswordString := string(hashedPassword)
	dobString := "2020-01-01"
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedUser := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPasswordString,
		Dob:      dob,
	}

	log.Println(hashedPassword)

	// SignupUsingOAuthで使用するmockメソッドを準備
	mockRepo.On("CreateUser", mock.MatchedBy(func(user *models.User) bool {
		return user.Name == expectedUser.Name &&
			user.Email == expectedUser.Email &&
			user.Dob == expectedUser.Dob &&
			bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
	})).Return(nil)

	// サインアップ
	err := testAuthService.Signup(name, email, dobString, password)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLoginUsingOAuthSuccess(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockUserRepository{}
	testAuthService := services.NewAuthService(mockRepo)

	// ユーザーモデルを準備
	name := "testuser"
	email := "test@example.com"
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedUser := &models.User{
		Name:  name,
		Email: email,
		Dob:   dob,
	}

	// FindUserByEmailで使用するmockメソッドを準備
	mockRepo.On("FindUserByEmail", email).Return(expectedUser, nil)

	// ログイン
	loginResponse, err := testAuthService.LoginUsingOAuth(expectedUser.Email)

	assert.NoError(t, err)
	assert.NotNil(t, loginResponse)
	mockRepo.AssertExpectations(t)
}

func TestLoginUsingOAuthUserNotFound(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockUserRepository{}
	testAuthService := services.NewAuthService(mockRepo)

	// ユーザーが存在しないemailを準備
	notExistEmail := "test@example.com"

	// FindUserByEmailで使用するmockメソッドを準備
	mockRepo.On("FindUserByEmail", notExistEmail).Return(nil, errors.New("user not found"))

	// ログイン
	loginResponse, err := testAuthService.LoginUsingOAuth(notExistEmail)

	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, loginResponse)
	mockRepo.AssertExpectations(t)
}

func TestLoginSuccess(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockUserRepository{}
	testAuthService := services.NewAuthService(mockRepo)

	// ユーザーモデルを準備
	name := "testuser"
	email := "test@example.com"
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPasswordString := string(hashedPassword)
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedUser := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPasswordString,
		Dob:      dob,
	}

	// FindUserByEmailで使用するmockメソッドを準備
	mockRepo.On("FindUserByEmail", mock.MatchedBy(func(email string) bool {
		return email == expectedUser.Email
	})).Return(expectedUser, nil)

	// ログイン
	loginResponse, err := testAuthService.Login(email, password)

	assert.NoError(t, err)
	assert.NotNil(t, loginResponse)
	mockRepo.AssertExpectations(t)
}

func TestLoginNotUserFound(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockUserRepository{}
	testAuthService := services.NewAuthService(mockRepo)

	// ユーザーが存在しないemailとpasswordを準備
	notExistEmail := "test@example.com"
	notExistPassword := "testpassword"

	// FindUserByEmailで使用するmockメソッドを準備
	mockRepo.On("FindUserByEmail", notExistEmail).Return(nil, errors.New("user not found"))

	// ログイン
	loginResponse, err := testAuthService.Login(notExistEmail, notExistPassword)

	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, loginResponse)
	mockRepo.AssertExpectations(t)
}

func TestLoginInvalidPassword(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockUserRepository{}
	testAuthService := services.NewAuthService(mockRepo)

	// ユーザーモデルを準備
	name := "testuser"
	email := "test@example.com"
	wrongPassword := "wrongpassword"
	correctPassword := "correctpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	hashedPasswordString := string(hashedPassword)
	dob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedUser := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPasswordString,
		Dob:      dob,
	}

	// FindUserByEmailで使用するmockメソッドを準備
	mockRepo.On("FindUserByEmail", mock.MatchedBy(func(email string) bool {
		return email == expectedUser.Email
	})).Return(expectedUser, nil)

	// 不正なパスワードでログイン
	loginResponse, err := testAuthService.Login(email, wrongPassword)

	assert.Equal(t, "invalid password", err.Error())
	assert.Nil(t, loginResponse)
	mockRepo.AssertExpectations(t)
}
