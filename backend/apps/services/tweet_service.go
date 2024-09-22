package services

import (
	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
)

type ITweetService interface {
	CreateTweet(userId uint, tweetTypeString string, content string) error
	GetTweet(id uint) (*models.Tweet, error)
	GetUserTweets(userId uint) ([]*models.Tweet, error)
}

type TweetService struct {
	repository repositories.ITweetRepository
}

func NewTweetService(repository repositories.ITweetRepository) ITweetService {
	return &TweetService{repository: repository}
}

func (s *TweetService) CreateTweet(userId uint, tweetTypeString string, content string) error {
	var (
		tweetType models.TweetType
		err       error
	)

	// stringで受け取ったtweetTypeStringをenumに変換
	tweetType, err = models.Str2TweetType(tweetTypeString)
	if err != nil {
		return err
	}

	tweet := &models.Tweet{
		UserID:  userId,
		Type:    tweetType,
		Content: content,
	}

	err = s.repository.CreateTweet(tweet)
	if err != nil {
		return err
	}

	return nil
}

func (s *TweetService) GetTweet(id uint) (*models.Tweet, error) {
	return s.repository.GetTweet(id)
}

func (s *TweetService) GetUserTweets(userId uint) ([]*models.Tweet, error) {
	return s.repository.GetUserTweets(userId)
}
