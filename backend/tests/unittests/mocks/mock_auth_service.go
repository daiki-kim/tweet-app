package mocks

import (
	"github.com/daiki-kim/tweet-app/backend/apps/services"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) SignupUsingOAuth(name, email, dobString string) error {
	args := m.Called(name, email, dobString)
	return args.Error(0)
}

func (m *MockAuthService) Signup(name, email, dobString, password string) error {
	args := m.Called(name, email, dobString, password)
	return args.Error(0)
}

func (m *MockAuthService) LoginUsingOAuth(email string) (*services.LoginResponse, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.LoginResponse), args.Error(1)
}
