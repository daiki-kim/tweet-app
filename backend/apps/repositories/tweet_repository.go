package repositories

import (
	"errors"
	"log"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"gorm.io/gorm"
)

type ITweetRepository interface {
	CreateTweet(tweet *models.Tweet) (*models.Tweet, error)
	GetTweet(id uint) (*models.Tweet, error)
	GetUserTweets(userId uint) ([]*models.Tweet, error)
	UpdateTweet(updateTweet *models.Tweet) (*models.Tweet, error)
}

type TweetRepository struct {
	DB *gorm.DB
}

func NewTweetRepository(db *gorm.DB) ITweetRepository {
	return &TweetRepository{DB: db}
}

func (r *TweetRepository) CreateTweet(tweet *models.Tweet) (*models.Tweet, error) {
	result := r.DB.Create(tweet)
	if result.Error != nil {
		log.Println("failed to create tweet: ", result.Error)
		return nil, result.Error
	}

	log.Println("created tweet: ", tweet)
	return tweet, nil
}

func (r *TweetRepository) GetTweet(id uint) (*models.Tweet, error) {
	var tweet models.Tweet

	result := r.DB.First(&tweet, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("tweet not found: ", result.Error)
		return nil, errors.New("tweet not found")
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &tweet, nil
}

func (r *TweetRepository) GetUserTweets(userId uint) ([]*models.Tweet, error) {
	var tweets []*models.Tweet

	result := r.DB.Where("user_id = ?", userId).Find(&tweets)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("tweets not found")
	} else if result.Error != nil {
		return nil, result.Error
	}

	return tweets, nil
}

func (r *TweetRepository) UpdateTweet(updateTweet *models.Tweet) (*models.Tweet, error) {
	result := r.DB.Save(updateTweet)
	if result.Error != nil {
		return nil, result.Error
	}

	return updateTweet, nil
}
