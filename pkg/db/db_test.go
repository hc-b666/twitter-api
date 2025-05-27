package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDB(t *testing.T) {
	t.Run("check invalid dsn", func(t *testing.T) {
		_, err := NewDB("invalid-dsn")
		assert.Error(t, err)
	})
	t.Run("check valid dsn", func(t *testing.T) {
		pool, err := NewDB("postgres://postgres:postgres@localhost:5432/postgres")
		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})
}
