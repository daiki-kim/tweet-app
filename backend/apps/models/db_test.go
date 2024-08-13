package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModels(t *testing.T) {
	var testUser User
	gotModels := GetModels()
	assert.Equal(t, gotModels[0], &testUser)
}

func TestNewDatabaseFactoryWithSQLite(t *testing.T) {
	db, err := NewDatabaseFactory(InstanceSqLite)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
