package tests

import (
	"os"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/daiki-kim/tweet-app/backend/configs"
	"github.com/stretchr/testify/suite"
)

const testDBName = "test.db"

type DBSQLiteSuite struct {
	suite.Suite
}

// tweetモデルはenum型がSQLiteで対応していないため省略
func getTestModels() []interface{} {
	return []interface{}{
		&models.User{},
	}
}

// sqliteのテストスイートをセットアップ
func (suite *DBSQLiteSuite) SetupSuite() {
	configs.Config.DBName = testDBName
	err := models.SetDatabase(models.InstanceSQLite)
	suite.Assert().Nil(err)

	for _, model := range getTestModels() {
		err := models.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

// sqliteのテストスイートをクリーンアップ
func (suite *DBSQLiteSuite) TearDownSuite() {
	err := os.Remove(configs.Config.DBName)
	suite.Assert().Nil(err)
}
