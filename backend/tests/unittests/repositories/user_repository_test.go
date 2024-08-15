package repositories_test

import (
	"log"
	"testing"
	"time"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
	"github.com/daiki-kim/tweet-app/backend/tests"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserTestSuite struct {
	tests.DBSQLiteSuite
	originalDB *gorm.DB
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	if models.DB == nil {
		log.Fatal("models.DB is nil")
	}
	suite.originalDB = models.DB
}

func (suite *UserTestSuite) AfterTest(suiteName, testName string) {
	models.DB = suite.originalDB
}

func (suite *UserTestSuite) TestUserRepository() {
	testDob := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	user := &models.User{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		// Dob:      testDob,
	}

	testUserRepository := repositories.NewUserRepository(models.DB)
	err := testUserRepository.CreateUser(user)

	suite.Nil(err)

	user, err = testUserRepository.FindUserByEmail("test@example.com")
	log.Println(user)

	suite.Nil(err)
	suite.Equal("testuser", user.Name)
	suite.Equal("test@example.com", user.Email)
	suite.Equal("testpassword", user.Password)
	suite.Equal(testDob, user.Dob)
}
