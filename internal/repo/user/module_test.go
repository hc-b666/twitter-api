package user

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"twitter-api/pkg/db"
	"twitter-api/pkg/logger"
)

func TestNewRepo(t *testing.T) {
	mockPool := new(db.MockPool)
	mockLogger := &logger.Logger{}

	repo, err := NewRepo(mockPool, mockLogger)

	require.NoError(t, err)
	require.NotNil(t, repo)
}

func TestRepo_GetByID(t *testing.T) {
	t.Run("error on get user by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(pgx.ErrNoRows)

		user, err := pool.GetByID(context.Background(), 1)
		assert.NotNil(t, err)
		assert.Nil(t, user)
	})
	t.Run("success on get user by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(nil)

		user, err := pool.GetByID(context.Background(), 0)
		assert.Nil(t, err)
		assert.NotNil(t, user)
	})
	t.Run("error on get user by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(errors.New("error"))

		user, err := pool.GetByID(context.Background(), 0)
		assert.NotNil(t, err)
		assert.Nil(t, user)
	})
}

func TestRepo_GetAll_Success(t *testing.T) {

	ctx := context.Background()
	mockPool := new(db.MockPool)
	mockRows := new(db.MockRows)
	mockLogger := &logger.Logger{}
	var rows pgx.Rows = mockRows
	mockPool.On("Query", mock.Anything, mock.Anything).Return(rows, nil)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(nil)
	mockRows.On("Close").Return(nil)

	r, _ := NewRepo(mockPool, mockLogger)

	adminU, err := r.GetAll(ctx)

	require.NoError(t, err)
	require.Len(t, adminU, 2)

	mockPool.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestRepo_GetAll_Empty(t *testing.T) {
	ctx := context.Background()
	mockPool := new(db.MockPool)
	mockRows := new(db.MockRows)
	mockLogger := &logger.Logger{}
	mockPool.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(nil)
	mockRows.On("Close").Return(nil)

	r, _ := NewRepo(mockPool, mockLogger)

	adminU, err := r.GetAll(ctx)

	require.NoError(t, err)
	require.Len(t, adminU, 0)

	mockPool.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestRepo_GetByEmail(t *testing.T) {
	t.Run("error on get user by email", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).Return(pgx.ErrNoRows)

		user, err := pool.GetByEmail(context.Background(), "")
		assert.NotNil(t, err)
		assert.Nil(t, user)
	})
	t.Run("success on get user by email", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).Return(nil)

		user, err := pool.GetByEmail(context.Background(), "")
		assert.Nil(t, err)
		assert.NotNil(t, user)
	})
	t.Run("error on get user by email", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		user, err := pool.GetByEmail(context.Background(), "")
		assert.NotNil(t, err)
		assert.Nil(t, user)
	})
}

func TestRepo_Create(t *testing.T) {
	userDTO := &RegisterUserDTO{Email: "email", Password: "password"}
	t.Run("error on create user", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(pgx.ErrNoRows)

		user, err := pool.Create(context.Background(), userDTO)
		assert.NotNil(t, err)
		assert.Zero(t, user)
	})
	t.Run("success on create user", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(nil)

		user, err := pool.Create(context.Background(), userDTO)
		assert.Nil(t, err)
		assert.NotNil(t, user)
	})
	t.Run("error on create user ", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(errors.New("error"))

		user, err := pool.Create(context.Background(), userDTO)
		assert.NotNil(t, err)
		assert.Zero(t, user)
	})
}
