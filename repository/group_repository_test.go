package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"noda/data/model"
	"noda/data/transfer"
	"noda/failure"
	"regexp"
	"testing"
)

const groupID string = "942d76f4-28b2-44be-8339-232b62c0ef22"

func TestGroupRepository_Save(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`SELECT make_group ($1, $2, $3);`)
		res   string
		err   error
		next  = &transfer.GroupCreation{Name: "name", Description: "desc"}
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, next.Name, next.Description).
			WillReturnRows(sqlmock.
				NewRows([]string{"make_group"}).
				AddRow(groupID))
		res, err = r.Save(userID, next)
		assert.NoError(t, err)
		assert.Equal(t, groupID, res)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, next.Name, next.Description).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.Save(userID, next)
		assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
		assert.Equal(t, "", res)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, next.Name, next.Description).
			WillReturnError(&pq.Error{})
		res, err = r.Save(userID, next)
		assert.Error(t, err)
		assert.Equal(t, "", res)
	})
}

func TestGroupRepository_FetchByID(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`SELECT * FROM fetch_group_by_id ($1, $2);`)
		res   *model.Group
		err   error
		group = &model.Group{
			ID:          uuid.MustParse(groupID),
			OwnerID:     uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			CreatedAt:   nil,
			UpdatedAt:   nil,
		}
		columns = []string{"id", "owner_id", "name", "description", "created_at", "updated_at"}
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnRows(sqlmock.
				NewRows(columns).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt))
		res, err = r.FetchByID(userID, groupID)
		assert.NoError(t, err)
		assert.Equal(t, group, res)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.FetchByID(userID, groupID)
		assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
		assert.Nil(t, res)
	})

	t.Run("group not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with ID"})
		res, err = r.FetchByID(userID, groupID)
		assert.ErrorIs(t, err, failure.ErrGroupNotFound)
		assert.Nil(t, res)
	})

	t.Run("deadline (5s) exceeded", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(errors.New("context deadline exceeded"))
		res, err = r.FetchByID(userID, groupID)
		assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
		assert.Nil(t, res)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(&pq.Error{})
		res, err = r.FetchByID(userID, groupID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

}

func TestGroupRepository_Fetch(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`
		SELECT "group_id" AS "id",
					 "owner_id",
					 "name",
					 "description",
					 "created_at",
					 "updated_at"
			FROM fetch_groups ($1, $2, $3, $4, $5);`)
		res       []*model.Group
		err       error
		page, rpp int64
		needle    = "name"
		sortBy    = "+name"
		group     = model.Group{
			ID:          uuid.New(),
			OwnerID:     uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			CreatedAt:   nil,
			UpdatedAt:   nil,
		}
		columns = []string{"id", "owner_id", "name", "description", "created_at", "updated_at"}
	)

	t.Run("success with 2 records", func(t *testing.T) {
		page, rpp = 1, 2
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortBy).
			WillReturnRows(sqlmock.
				NewRows(columns).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt))
		res, err = r.Fetch(userID, page, rpp, needle, sortBy)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("success with the default number of records (10)", func(t *testing.T) {
		page, rpp = 1, -1000
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortBy).
			WillReturnRows(sqlmock.
				NewRows(columns).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt))
		res, err = r.Fetch(userID, page, rpp, needle, sortBy) /* Should set 'rpp' to 10.  */
		assert.NoError(t, err)
		assert.Len(t, res, 10)
	})

	t.Run("success with custom pagination and RPP", func(t *testing.T) {
		page, rpp = 2, 5
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortBy).
			WillReturnRows(sqlmock.
				NewRows(columns).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt))
		res, err = r.Fetch(userID, page, rpp, needle, sortBy)
		assert.NoError(t, err)
		assert.Len(t, res, 5)
	})

	t.Run("success with searching", func(t *testing.T) {
		page, rpp, needle = 1, 7, "name"
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortBy).
			WillReturnRows(sqlmock.
				NewRows(columns).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt))
		res, err = r.Fetch(userID, page, rpp, needle, sortBy)
		assert.NoError(t, err)
		assert.Len(t, res, 7)
	})

	t.Run("no response/error for weird needle", func(t *testing.T) {
		page, rpp, needle = 1, 5, "aljfkjaksjpiwquramakjsfasjfkjwpoijefj"
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortBy).
			WillReturnRows(sqlmock.NewRows(columns))
		res, err = r.Fetch(userID, page, rpp, needle, sortBy)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res, 0)
	})

	t.Run("user not found", func(t *testing.T) {
		page, rpp = 1, 10
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortBy).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.Fetch(userID, page, rpp, needle, sortBy)
		assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
		assert.Nil(t, res)
	})

	t.Run("deadline (5s) exceeded", func(t *testing.T) {
		page, rpp = 1, 10
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortBy).
			WillReturnError(errors.New("context deadline exceeded"))
		res, err = r.Fetch(userID, page, rpp, needle, sortBy)
		assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
		assert.Nil(t, res)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortBy).
			WillReturnError(&pq.Error{})
		res, err = r.Fetch(userID, page, rpp, needle, sortBy)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("unexpected scanning error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortBy).
			WillReturnRows(sqlmock.
				NewRows([]string{"group_id", "owner_id", "name", "description", "created_at", "updated_at"}).
				AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt))
		res, err = r.Fetch(userID, page, rpp, needle, sortBy)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGroupRepository_Update(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`SELECT update_group ($1, $2, $3, $4);`)
		res   bool
		err   error
		up    = &transfer.GroupUpdate{}
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID, up.Name, up.Description).
			WillReturnRows(sqlmock.
				NewRows([]string{"update_group"}).
				AddRow(true))
		res, err = r.Update(userID, groupID, up)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("did not update and no error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID, up.Name, up.Description).
			WillReturnRows(sqlmock.
				NewRows([]string{"update_group"}).
				AddRow(false))
		res, err = r.Update(userID, groupID, up)
		assert.False(t, res)
		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID, up.Name, up.Description).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.Update(userID, groupID, up)
		assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
		assert.False(t, res)
	})

	t.Run("group not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID, up.Name, up.Description).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with ID"})
		res, err = r.Update(userID, groupID, up)
		assert.ErrorIs(t, err, failure.ErrGroupNotFound)
		assert.False(t, res)
	})

	t.Run("deadline (5s) exceeded", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID, up.Name, up.Description).
			WillReturnError(errors.New("context deadline exceeded"))
		res, err = r.Update(userID, groupID, up)
		assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
		assert.False(t, res)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID, up.Name, up.Description).
			WillReturnError(&pq.Error{})
		res, err = r.Update(userID, groupID, up)
		assert.Error(t, err)
		assert.False(t, res)
	})
}

func TestGroupRepository_Remove(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`SELECT delete_group ($1, $2);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnRows(sqlmock.
				NewRows([]string{"delete_group"}).
				AddRow(true))
		res, err = r.Remove(userID, groupID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("did not delete and no error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnRows(sqlmock.
				NewRows([]string{"delete_group"}).
				AddRow(false))
		res, err = r.Remove(userID, groupID)
		assert.False(t, res)
		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.Remove(userID, groupID)
		assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
		assert.False(t, res)
	})

	t.Run("group not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with ID"})
		res, err = r.Remove(userID, groupID)
		assert.ErrorIs(t, err, failure.ErrGroupNotFound)
		assert.False(t, res)
	})

	t.Run("deadline (5s) exceeded", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(errors.New("context deadline exceeded"))
		res, err = r.Remove(userID, groupID)
		assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
		assert.False(t, res)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(&pq.Error{})
		res, err = r.Remove(userID, groupID)
		assert.Error(t, err)
		assert.False(t, res)
	})
}
