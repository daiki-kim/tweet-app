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

type FollowerTestSuite struct {
	tests.DBSQLiteSuite
	originalDB *gorm.DB
}

func TestFollowerTestSuite(t *testing.T) {
	suite.Run(t, new(FollowerTestSuite))
}

func (suite *FollowerTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	if models.DB == nil {
		log.Fatal("models.DB is nil")
	}
	suite.originalDB = models.DB
}

func (suite *FollowerTestSuite) AfterTest(suiteName, testName string) {
	models.DB = suite.originalDB
}

func (suite *FollowerTestSuite) TestFollowerRepository() {
	testuser1 := &models.User{
		Name:     "testuser1",
		Email:    "test1@example.com",
		Password: "testpassword",
		Dob:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	testuser2 := &models.User{
		Name:     "testuser2",
		Email:    "test2@example.com",
		Password: "testpassword",
		Dob:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	testFollower := &models.Follower{
		FollowerID: 1,
		FolloweeID: 2,
	}

	testUserRepository := repositories.NewUserRepository(models.DB)
	testFollowerRepository := repositories.NewFollowerRepository(models.DB)

	err := testUserRepository.CreateUser(testuser1)
	suite.Nil(err)
	err = testUserRepository.CreateUser(testuser2)
	suite.Nil(err)

	returnFollower, err := testFollowerRepository.CreateFollower(testFollower)
	suite.Nil(err)
	suite.Equal(uint(1), returnFollower.FollowerID)
	suite.Equal(uint(2), returnFollower.FolloweeID)
}
