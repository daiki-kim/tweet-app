package models_test

import (
	"testing"

	"github.com/daiki-kim/tweet-app/backend/apps/models"
	"github.com/stretchr/testify/assert"
)

func TestGetModels(t *testing.T) {
	var testUser models.User
	gotModels := models.GetModels()
	assert.Equal(t, gotModels[0], &testUser)
}

func TestNewDatabaseFactoryWithSQLite(t *testing.T) {
	db, err := models.NewDatabaseFactory(models.InstanceSQLite)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
