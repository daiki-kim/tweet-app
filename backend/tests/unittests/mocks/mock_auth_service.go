package mocks

import "github.com/stretchr/testify/mock"

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
