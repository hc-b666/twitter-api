package comment

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"twitter-api/pkg/db"
)

func TestNewRepo(t *testing.T) {
	mockDB := new(db.MockPool)
	r, err := NewRepo(mockDB)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestRepo_Update(t *testing.T) {
	t.Run("error on create comment", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("can't create comment")).Once()

		err := pool.Update(context.Background(), 1, "")
		assert.NotNil(t, err)
	})

	t.Run("success on create task", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("UPDATE 1"), nil).Once()

		err := pool.Update(context.Background(), 1, "this is comment")
		assert.Nil(t, err)
	})
}

func TestRepo_GetCommentByID(t *testing.T) {
	t.Run("error on get comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Return(pgx.ErrNoRows)

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
			mock.Anything, mock.Anything).Return(nil)

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
			mock.Anything, mock.Anything).Return(errors.New("error"))

		comment, err := pool.GetByID(context.Background(), 0)
		assert.NotNil(t, err)
		assert.Nil(t, comment)
	})
}

func TestRepo_SoftDelete(t *testing.T) {
	t.Run("error on delete comment", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("can't delete comment")).Once()

		err := pool.SoftDelete(context.Background(), 1)
		assert.NotNil(t, err)
	})

	t.Run("success on delete comment", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("UPDATE 1"), nil).Once()

		err := pool.SoftDelete(context.Background(), 1)
		assert.Nil(t, err)
	})
}

func TestRepo_HardDelete(t *testing.T) {
	t.Run("error on delete comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(pgx.ErrNoRows)

		err := pool.HardDelete(context.Background(), 1)
		assert.NotNil(t, err)

	})
	t.Run("success on delete comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(nil)

		err := pool.HardDelete(context.Background(), 0)
		assert.Nil(t, err)

	})
	t.Run("error on delete comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(errors.New("error"))

		err := pool.HardDelete(context.Background(), 0)
		assert.NotNil(t, err)

	})
}
