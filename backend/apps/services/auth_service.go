package services

import (
	"errors"
	"log"
	"time"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
	"github.com/daiki-kim/tweet-app/backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	SignupUsingOAuth(name, email, dobString string) error
	Signup(name, email, dobString, password string) error
	LoginUsingOAuth(email string) (*LoginResponse, error)
	Login(email, password string) (*LoginResponse, error)
}

type AuthService struct {
	repository repositories.IUserRepository
}

func NewAuthService(repository repositories.IUserRepository) IAuthService {
	return &AuthService{repository: repository}
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// string型をTime型に変換
func Str2time(t string) (time.Time, error) {
	return time.Parse("2006-01-02", t)
}

// OAuth、Normal共通で利用するユーザーモデルを準備
func PrepareBaseUserModel(name, email, dobString string) (*models.User, error) {
	dob, err := Str2time(dobString)
	if err != nil {
		log.Println("failed to convert string to time: ", err)
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
		log.Println("failed to prepare user model: ", err)
		return err
	}

	return s.repository.CreateUser(user)
}

// ユーザー入力情報を使用するNormalのサインアップ
// パスワードが必要
func (s *AuthService) Signup(name, email, dobString, password string) error {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password: ", err)
		return err
	}

	user, err := PrepareBaseUserModel(name, email, dobString)
	if err != nil {
		log.Println("failed to prepare user model: ", err)
		return err
	}

	user.Password = string(hashedPassword)
	return s.repository.CreateUser(user)
}

// OAuthからのログイン
// googleから取得したemailを使用してtokenを発行
func (s *AuthService) LoginUsingOAuth(email string) (*LoginResponse, error) {
	// emailからユーザーモデルを取得
	user, err := s.repository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Claim構造体のポインタを生成して、トークンを発行
	claim := auth.NewClaim(user.Email)
	token, err := claim.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Claim構造体のポインタを生成して、リフレッシュトークンを発行
	refreshTokenClaim := auth.NewClaim(user.Email)
	refreshToken, err := refreshTokenClaim.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	loginResponse := &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
	return loginResponse, nil
}

// Normalログイン
// ユーザーが入力したemailとpasswordを使用してtokenを発行
func (s *AuthService) Login(email, password string) (*LoginResponse, error) {
	// emailからユーザーモデルを取得
	user, err := s.repository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// ハッシュ化されたパスワードと入力されたパスワードを比較
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Claim構造体のポインタを生成してトークンを発行
	claim := auth.NewClaim(user.Email)
	token, err := claim.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Claim構造体のポインタを生成してリフレッシュトークンを発行
	refreshTokenClaim := auth.NewClaim(user.Email)
	refreshToken, err := refreshTokenClaim.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	loginResponse := &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
	return loginResponse, nil
}
