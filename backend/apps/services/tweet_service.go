package services

import (
	"errors"

	"github.com/daiki-kim/tweet-app/backend/apps/dtos"
	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/apps/repositories"
)

type ITweetService interface {
	CreateTweet(userId uint, tweetTypeString string, content string) (*models.Tweet, error)
	GetTweet(id uint) (*models.Tweet, error)
	GetUserTweets(userId uint) ([]*models.Tweet, error)
	UpdateTweet(id, userId uint, inputTweet *dtos.UpdateTweetInput) (*models.Tweet, error)
	DeleteTweet(id, userId uint) error
}

type TweetService struct {
	repository repositories.ITweetRepository
}

func NewTweetService(repository repositories.ITweetRepository) ITweetService {
	return &TweetService{repository: repository}
}

func (s *TweetService) CreateTweet(userId uint, tweetTypeString string, content string) (*models.Tweet, error) {
	var (
		tweetType models.TweetType
		err       error
	)

	// stringで受け取ったtweetTypeStringをenumに変換
	tweetType, err = models.Str2TweetType(tweetTypeString)
	if err != nil {
		return nil, err
	}

	tweet := &models.Tweet{
		UserID:  userId,
		Type:    tweetType,
		Content: content,
	}

	return s.repository.CreateTweet(tweet)
}

func (s *TweetService) GetTweet(id uint) (*models.Tweet, error) {
	return s.repository.GetTweet(id)
}

func (s *TweetService) GetUserTweets(userId uint) ([]*models.Tweet, error) {
	return s.repository.GetUserTweets(userId)
}

func (s *TweetService) UpdateTweet(id, userId uint, inputTweet *dtos.UpdateTweetInput) (*models.Tweet, error) {
	updatedTweet, err := s.repository.GetTweet(id)
	if err != nil {
		return nil, err
	}

	if updatedTweet.UserID != userId {
		return nil, errors.New("this tweet is not yours")
	}

	if inputTweet.Type != "" {
		tweetType, err := models.Str2TweetType(inputTweet.Type)
		if err != nil {
			return nil, err
		}
		updatedTweet.Type = tweetType
	}

	if inputTweet.Content != "" {
		updatedTweet.Content = inputTweet.Content
	}

	return s.repository.UpdateTweet(updatedTweet)
}

func (s *TweetService) DeleteTweet(id, userId uint) error {
	targetTweet, err := s.repository.GetTweet(id)
	if err != nil {
		return err
	}

	if targetTweet.UserID != userId {
		return errors.New("this tweet is not yours")
	}

	return s.repository.DeleteTweet(id)
}
