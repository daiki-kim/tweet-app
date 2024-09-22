package repositories

import (
	"log"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"gorm.io/gorm"
)

type ITweetRepository interface {
	CreateTweet(tweet *models.Tweet) error
}

type TweetRepository struct {
	DB *gorm.DB
}

func NewTweetRepository(db *gorm.DB) ITweetRepository {
	return &TweetRepository{DB: db}
}

func (r *TweetRepository) CreateTweet(tweet *models.Tweet) error {
	result := r.DB.Create(tweet)
	if result.Error != nil {
		log.Println("failed to create tweet: ", result.Error)
		return result.Error
	}

	log.Println("created tweet: ", tweet)
	return nil
}
