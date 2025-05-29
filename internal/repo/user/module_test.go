package user

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"twitter-api/pkg/db"
)

func TestRepo_GetPostByID(t *testing.T) {
	t.Run("error on get comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).Return(pgx.ErrNoRows)

		comment, err := pool.GetByID(context.Background(), 1)
		assert.NotNil(t, err)
		assert.Nil(t, comment)
	})
	t.Run("success on get comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).Return(nil)

		comment, err := pool.GetByID(context.Background(), 0)
		assert.Nil(t, err)
		assert.NotNil(t, comment)
	})
	t.Run("error on get comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		comment, err := pool.GetByID(context.Background(), 0)
		assert.NotNil(t, err)
		assert.Nil(t, comment)
	})
}
