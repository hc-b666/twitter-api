package post

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"twitter-api/pkg/db"
)

func TestNewRepo(t *testing.T) {
	mockDB := new(db.MockPool)
	r, err := NewRepo(mockDB)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}

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

		err := pool.HardDelete(context.Background(), 1)
		assert.NotNil(t, err)

	})

	t.Run("error on delete comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)

		err := pool.HardDelete(context.Background(), 0)
		assert.NotNil(t, err)

	})
}

func TestRepo_GetAllPosts_Success(t *testing.T) {

	ctx := context.Background()
	mockPool := new(db.MockPool)
	mockRows := new(db.MockRows)
	var rows pgx.Rows = mockRows
	mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything).Return(rows, nil)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(nil)
	mockRows.On("Close").Return(nil)

	r, _ := NewRepo(mockPool)

	comments, err := r.GetAll(ctx, 1, 2)

	require.NoError(t, err)
	require.Len(t, comments, 2)

	mockPool.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestRepo_GetAllPosts_Empty(t *testing.T) {
	ctx := context.Background()
	mockPool := new(db.MockPool)
	mockRows := new(db.MockRows)

	mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything).Return(mockRows, nil)

	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(nil)
	mockRows.On("Close").Return(nil)

	r, _ := NewRepo(mockPool)

	comments, err := r.GetAll(ctx, 1, 2)

	require.NoError(t, err)
	require.Len(t, comments, 0)

	mockPool.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestRepo_Create(t *testing.T) {
	t.Run("error on create post", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(pgx.ErrNoRows)

		post, err := pool.Create(context.Background(), 1, "ff", "ff")
		assert.NotNil(t, err)
		assert.Zero(t, post)
	})
	t.Run("success on create comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(nil)

		post, err := pool.Create(context.Background(), 1, "bb", "dd")

		assert.Nil(t, err)
		assert.NotNil(t, post)
	})
	t.Run("error on create comment by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(errors.New("error"))

		post, err := pool.Create(context.Background(), 1, "dd", "sa")

		assert.NotNil(t, err)
		assert.Zero(t, post)
	})
}

func TestRepo_IsAuthor(t *testing.T) {
	t.Run("error on author of post", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(pgx.ErrNoRows)

		post, err := pool.IsAuthor(context.Background(), 1, 9)
		assert.NotNil(t, err)
		assert.Zero(t, post)
	})
	t.Run("success on author of post ", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(nil)

		comment, err := pool.IsAuthor(context.Background(), 1, 9)

		assert.Nil(t, err)
		assert.NotNil(t, comment)
	})
	t.Run("error on author of post", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything).Return(errors.New("error"))

		comment, err := pool.IsAuthor(context.Background(), 1, 9)

		assert.NotNil(t, err)
		assert.Zero(t, comment)
	})
}

func TestRepo_GetByUserID_Success(t *testing.T) {

	ctx := context.Background()
	mockPool := new(db.MockPool)
	mockRows := new(db.MockRows)
	var rows pgx.Rows = mockRows
	mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(rows, nil)
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(nil)
	mockRows.On("Close").Return(nil)

	r, _ := NewRepo(mockPool)

	comments, err := r.GetByUserID(ctx, 1)

	require.NoError(t, err)
	require.Len(t, comments, 2)

	mockPool.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestRepo_GetByUserID_Empty(t *testing.T) {
	ctx := context.Background()
	mockPool := new(db.MockPool)
	mockRows := new(db.MockRows)

	mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)

	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(nil)
	mockRows.On("Close").Return(nil)

	r, _ := NewRepo(mockPool)

	comments, err := r.GetByUserID(ctx, 1)

	require.NoError(t, err)
	require.Len(t, comments, 0)

	mockPool.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestRepo_UpdateContent(t *testing.T) {
	t.Run("error on create comment", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("can't create comment")).Once()

		err := pool.UpdateContent(context.Background(), 1, "")
		assert.NotNil(t, err)
	})

	t.Run("success on create task", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("UPDATE 1"), nil).Once()

		err := pool.UpdateContent(context.Background(), 1, "this is comment")
		assert.Nil(t, err)
	})
}

func TestRepo_UpdateFileURL(t *testing.T) {
	t.Run("error on update fileURL", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("can't update fileURL")).Once()

		err := pool.UpdateFileURL(context.Background(), 1, "")
		assert.NotNil(t, err)
	})

	t.Run("success on update fileURL", func(t *testing.T) {
		mockPool := new(db.MockPool)
		defer mockPool.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("UPDATE 1"), nil)

		err := pool.UpdateFileURL(context.Background(), 1, "pp.png")
		assert.Nil(t, err)
	})
}

func TestRepo_Update(t *testing.T) {
	t.Run("error on update post by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).Return(pgx.ErrNoRows)

		post, err := pool.Update(context.Background(), 1, "", "")
		assert.NotNil(t, err)
		assert.Nil(t, post)
	})
	t.Run("success on update post by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).Return(nil)

		comment, err := pool.Update(context.Background(), 0, "da", "sa")
		assert.Nil(t, err)
		assert.NotNil(t, comment)
	})
	t.Run("error on update post by id", func(t *testing.T) {
		mockPool := new(db.MockPool)
		mockRow := new(db.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repo{db: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		comment, err := pool.Update(context.Background(), 0, "mycontent", "ooo.png")
		assert.NotNil(t, err)
		assert.Nil(t, comment)
	})
}
