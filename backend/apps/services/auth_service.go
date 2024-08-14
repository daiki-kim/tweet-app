package services

import (
	"time"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	SignupUsingOAuth(name, email, dobString string) error
	Signup(name, email, password, dobString string) error
}

type AuthService struct {
	repository repositories.IUserRepository
}

func NewAuthService(repository repositories.IUserRepository) IAuthService {
	return &AuthService{repository: repository}
}

// string型をTime型に変換
func Str2time(t string) (time.Time, error) {
	return time.Parse("2006-01-02", t)
}

// OAuth、Normal共通で利用するユーザーモデルを準備
func PrepareBaseUserModel(name, email, dobString string) (*models.User, error) {
	dob, err := Str2time(dobString)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:  name,
		Email: email,
		Dob:   dob,
	}

	return user, nil
}

// OAuthでGoogleからのユーザーデータを使用してサインアップ
func (s *AuthService) SignupUsingOAuth(name, email, dobString string) error {
	user, err := PrepareBaseUserModel(name, email, dobString)
	if err != nil {
		return err
	}

	return s.repository.CreateUser(user)
}

// ユーザー入力情報を使用するNormalのサインアップ
// パスワードが必要
func (s *AuthService) Signup(name, email, password, dobString string) error {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := PrepareBaseUserModel(name, email, dobString)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repository.CreateUser(user)
}
