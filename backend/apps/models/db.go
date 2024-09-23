package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/daiki-kim/tweet-app/backend/configs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	InstanceSQLite int = iota
	InstanceMySQL
)

var (
	DB                            *gorm.DB
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Tweet{},
		&Follower{},
	}
}

// create new database
func NewDatabaseFactory(instance int) (db *gorm.DB, err error) {
	switch instance {
	case InstanceMySQL:
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
			"tweet_app",
			"tweet_password",
			"localhost",
			3306,
			"tweet_app_database",
		)
		// dsn := fmt.Sprintf(
		// 	"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		// 	configs.Config.DBUser,
		// 	configs.Config.DBPassword,
		// 	configs.Config.DBHost,
		// 	configs.Config.DBPort,
		// 	configs.Config.DBName,
		// )
		log.Println(dsn)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	case InstanceSQLite:
		db, err = gorm.Open(sqlite.Open(configs.Config.DBName), &gorm.Config{})

	default:
		return nil, errInvalidSQLDatabaseInstance
	}

	return db, err
}

// initialize database
func SetDatabase(instance int) (err error) {
	db, err := NewDatabaseFactory(instance)
	if err != nil {
		return err
	}

	DB = db
	return err
}
