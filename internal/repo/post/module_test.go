package post

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"twitter-api/pkg/db"
)

func TestNewRepo(t *testing.T) {
	mockDB := new(db.MockPool)
	r, err := NewRepo(mockDB)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}
