// repository unit test by using sqlite

package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// define tweet type
type TweetType string

// define the enum of tweet type
const (
	Text  TweetType = "text"
	Image TweetType = "image"
	Video TweetType = "video"
)

// convert the string to the enum to be used in the database
func Str2TweetType(tweetTypeString string) (TweetType, error) {
	switch tweetTypeString {
	case "text":
		return Text, nil
	case "image":
		return Image, nil
	case "video":
		return Video, nil
	default:
		return "", errors.New("invalid tweet type")
	}
}

// convert the enum from the database to the custom type
func (t *TweetType) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*t = TweetType(v) // convert string to enum
	case []byte:
		*t = TweetType(string(v)) // convert []byte to string before convert to enum
	default:
		return errors.New("type assertion to string failed")
	}

	return nil
	// str, ok := value.(string)
	// if !ok {
	// 	return errors.New("type assertion to string failed")
	// }

	// *t = TweetType(str)
	// return nil
}

// convert the custom type to the enum
func (t TweetType) Value() (driver.Value, error) {
	switch t {
	case Text, Image, Video:
		return string(t), nil
	default:
		return nil, fmt.Errorf("unsupported tweet type: %s", t)
	}
}

type Tweet struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Type      TweetType `gorm:"type:enum('text', 'image', 'video');not null" json:"type"`
	Content   string    `gorm:"type:varchar(140);not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// relations
	// User情報をTweetと一緒に取得したい場合はPreload("User")を使用する
	User *User `gorm:"foreignKey:UserID;references:ID" json:"user"`
}
