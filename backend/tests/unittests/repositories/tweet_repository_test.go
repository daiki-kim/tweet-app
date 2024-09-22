// repository unit test by using mock
// TODO: ENUM型も使えるDBでテストを実装する 2024/09/22

package repositories_test

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
	*gorm.DB
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}
