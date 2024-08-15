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
	testUser := &models.User{
		Name:  "testuser",
		Email: "test@example.com",
		Dob:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// FindUserByEmailで使用するmockメソッドを準備
	mockRepo.On("FindUserByEmail", testUser.Email).Return(testUser, nil)

	// ログイン
	loginResponse, err := testAuthService.LoginUsingOAuth(testUser.Email)
	log.Println(loginResponse)

	assert.NoError(t, err)
	assert.NotNil(t, loginResponse)
	mockRepo.AssertExpectations(t)
}

func TestLoginUsingOAuthUserNotFound(t *testing.T) {
	// モックレポジトリを準備
	mockRepo := &mocks.MockUserRepository{}
	testAuthService := services.NewAuthService(mockRepo)

	// 存在しないemailを使用するユーザーモデルを準備
	notExistUser := &models.User{
		Email: "test@example.com",
	}

	// FindUserByEmailで使用するmockメソッドを準備
	mockRepo.On("FindUserByEmail", notExistUser.Email).Return(nil, errors.New("user not found"))

	// ログイン
	loginResponse, err := testAuthService.LoginUsingOAuth(notExistUser.Email)

	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, loginResponse)
	mockRepo.AssertExpectations(t)
}
