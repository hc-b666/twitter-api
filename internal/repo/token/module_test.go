package token

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"twitter-api/pkg/db"
	"twitter-api/pkg/types"
)

func TestNewRepo(t *testing.T) {
	mockDB := new(db.MockPool)
	r, err := NewRepo(mockDB)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestRepo_GetByToken(t *testing.T) {
	t.Run("error on get token", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(pgx.ErrNoRows)

		token, err := pool.GetByToken(context.Background(), "dsada2313aad")
		assert.NotNil(t, err)
		assert.Nil(t, token)
	})
	t.Run("success on get token", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(nil)

		comment, err := pool.GetByToken(context.Background(), "mmkadojad")
		assert.Nil(t, err)
		assert.NotNil(t, comment)
	})
	t.Run("error on get token", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(errors.New("error"))

		comment, err := pool.GetByToken(context.Background(), "dhsfhskd")
		assert.NotNil(t, err)
		assert.Nil(t, comment)
	})
}

func TestRepo_Create(t *testing.T) {
	t.Run("error on create token", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("can't create token")).Once()

		accToken, refreshToken, err := pool.Create(context.Background(), 1, types.Regular)
		assert.NotNil(t, err)
		assert.Empty(t, accToken)
		assert.Empty(t, refreshToken)
	})

	t.Run("success on create task", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("INSERT 1"), nil).Once()

		accToken, refreshToken, err := pool.Create(context.Background(), 1, types.Regular)
		assert.Nil(t, err)
		assert.NotNil(t, accToken)
		assert.NotNil(t, refreshToken)
	})
}
