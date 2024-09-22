package services

import (
	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
)

type ITweetService interface {
	CreateTweet(userId uint, tweetTypeString string, content string) error
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
