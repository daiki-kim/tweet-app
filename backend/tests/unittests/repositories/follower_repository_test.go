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
	// prepare test user data
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
	testuser3 := &models.User{
		Name:     "testuser3",
		Email:    "test3@example.com",
		Password: "testpassword",
		Dob:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// prepare test follower data
	testFollower1Follows2 := &models.Follower{
		FollowerID: 1,
		FolloweeID: 2,
	}
	testFollower3Follows2 := &models.Follower{
		FollowerID: 3,
		FolloweeID: 2,
	}
	testFollower1Follows3 := &models.Follower{
		FollowerID: 1,
		FolloweeID: 3,
	}

	// prepare test repository
	testUserRepository := repositories.NewUserRepository(models.DB)
	testFollowerRepository := repositories.NewFollowerRepository(models.DB)

	// create users
	err := testUserRepository.CreateUser(testuser1)
	suite.Nil(err)
	err = testUserRepository.CreateUser(testuser2)
	suite.Nil(err)
	err = testUserRepository.CreateUser(testuser3)
	suite.Nil(err)

	// create followers
	follower1Follows2, err := testFollowerRepository.CreateFollower(testFollower1Follows2)
	suite.Nil(err)
	suite.Equal(uint(1), follower1Follows2.FollowerID)
	suite.Equal(uint(2), follower1Follows2.FolloweeID)

	follower3Follows2, err := testFollowerRepository.CreateFollower(testFollower3Follows2)
	suite.Nil(err)
	suite.Equal(uint(3), follower3Follows2.FollowerID)
	suite.Equal(uint(2), follower3Follows2.FolloweeID)

	follower1Follows3, err := testFollowerRepository.CreateFollower(testFollower1Follows3)
	suite.Nil(err)
	suite.Equal(uint(1), follower1Follows3.FollowerID)
	suite.Equal(uint(3), follower1Follows3.FolloweeID)

	// get follower
	follower, err := testFollowerRepository.GetFollower(1)
	suite.Nil(err)
	suite.Equal(uint(1), follower.FollowerID)
	suite.Equal(uint(2), follower.FolloweeID)

	// get followees
	// get follower datas user1 follows
	followees, err := testFollowerRepository.GetFollowees(1)
	suite.Nil(err)
	suite.Equal(2, len(followees))
	suite.Equal(testuser2.Name, followees[0].Followee.Name)
	suite.Equal(testuser2.Email, followees[0].Followee.Email)
	suite.Equal(testuser3.Name, followees[1].Followee.Name)
	suite.Equal(testuser3.Email, followees[1].Followee.Email)

	// get followers
	// get follower datas user2 is followed
	followers, err := testFollowerRepository.GetFollowers(2)
	suite.Nil(err)
	suite.Equal(2, len(followers))
	suite.Equal(testuser1.Name, followers[0].Follower.Name)
	suite.Equal(testuser1.Email, followers[0].Follower.Email)
	suite.Equal(testuser3.Name, followers[1].Follower.Name)
	suite.Equal(testuser3.Email, followers[1].Follower.Email)

	// delete followers
	err = testFollowerRepository.DeleteFollower(1) // delete follower1
	suite.Nil(err)
}
