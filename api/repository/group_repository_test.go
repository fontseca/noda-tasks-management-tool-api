package repository

import (
	"errors"
	"noda/api/data/model"
	"noda/api/data/transfer"
	"noda/failure"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const groupID string = "942d76f4-28b2-44be-8339-232b62c0ef22"

func TestInsertGroup(t *testing.T) {
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

	/* Success.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, next.Name, next.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"make_group"}).
			AddRow(groupID))
	res, err = r.InsertGroup(userID, next)
	assert.NoError(t, err)
	assert.Equal(t, groupID, res)

	/* User not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, next.Name, next.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.InsertGroup(userID, next)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Equal(t, "", res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, next.Name, next.Description).
		WillReturnError(&pq.Error{})
	res, err = r.InsertGroup(userID, next)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func TestFetchGroupByID(t *testing.T) {
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
			IsArchived:  false,
			ArchivedAt:  time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{
			"id", "owner_id", "name", "description", "is_archived",
			"archived_at", "created_at", "updated_at"}
	)

	/* Success.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(
				group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived,
				group.ArchivedAt, group.CreatedAt, group.UpdatedAt))
	res, err = r.FetchGroupByID(userID, groupID)
	assert.NoError(t, err)
	assert.Equal(t, group, res)

	/* User not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.FetchGroupByID(userID, groupID)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Nil(t, res)

	/* Group not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with ID"})
	res, err = r.FetchGroupByID(userID, groupID)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.Nil(t, res)

	/* Deadline (5s) exceeded.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.FetchGroupByID(userID, groupID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnError(&pq.Error{})
	res, err = r.FetchGroupByID(userID, groupID)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestFetchGroups(t *testing.T) {
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
					 "is_archived",
					 "archived_at",
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
			IsArchived:  false,
			ArchivedAt:  time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{
			"id", "owner_id", "name", "description", "is_archived",
			"archived_at", "created_at", "updated_at"}
	)

	/* Success with 2 records.  */

	page, rpp = 1, 2
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt))
	res, err = r.FetchGroups(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	/* Success with the default number of records (10).  */

	page, rpp = 1, -1000
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt))
	res, err = r.FetchGroups(userID, page, rpp, needle, sortBy) /* Should set `rpp' to 10.  */
	assert.NoError(t, err)
	assert.Len(t, res, 10)

	/* Success with custom pagination and RPP.  */

	page, rpp = 2, 5
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt))
	res, err = r.FetchGroups(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 5)

	/* Success with searching.  */

	page, rpp, needle = 1, 7, "name"
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt))
	res, err = r.FetchGroups(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 7)

	/* There should not be a response for a weird needle and neither should be
	   an error.  */

	page, rpp, needle = 1, 5, "aljfkjaksjpiwquramakjsfasjfkjwpoijefj"
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.NewRows(columns))
	res, err = r.FetchGroups(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)

	/* User not found.  */

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.FetchGroups(userID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Nil(t, res)

	/* Deadline (5s) exceeded.  */

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.FetchGroups(userID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{})
	res, err = r.FetchGroups(userID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)

	/* Unexpected scanning error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"group_id", "owner_id", "name", "description", "is_archived",
				"archived_at", "created_at", "updated_at"}).
			AddRow(group.ID, group.OwnerID, group.Name, group.Description, group.IsArchived, group.ArchivedAt, group.CreatedAt, group.UpdatedAt))
	res, err = r.FetchGroups(userID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestUpdateGroup(t *testing.T) {
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

	/* Success.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, up.Name, up.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_group"}).
			AddRow(true))
	res, err = r.UpdateGroup(userID, groupID, up)
	assert.True(t, res)
	assert.NoError(t, err)

	/* Did not update and no error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, up.Name, up.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_group"}).
			AddRow(false))
	res, err = r.UpdateGroup(userID, groupID, up)
	assert.False(t, res)
	assert.NoError(t, err)

	/* User not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, up.Name, up.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.UpdateGroup(userID, groupID, up)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.False(t, res)

	/* Group not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, up.Name, up.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with ID"})
	res, err = r.UpdateGroup(userID, groupID, up)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.False(t, res)

	/* Deadline (5s) exceeded.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, up.Name, up.Description).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.UpdateGroup(userID, groupID, up)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.False(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, up.Name, up.Description).
		WillReturnError(&pq.Error{})
	res, err = r.UpdateGroup(userID, groupID, up)
	assert.Error(t, err)
	assert.False(t, res)
}

func TestDeleteGroup(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`SELECT delete_group ($1, $2);`)
		res   bool
		err   error
	)

	/* Success.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnRows(sqlmock.
			NewRows([]string{"delete_group"}).
			AddRow(true))
	res, err = r.DeleteGroup(userID, groupID)
	assert.True(t, res)
	assert.NoError(t, err)

	/* Did not delete and no error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnRows(sqlmock.
			NewRows([]string{"delete_group"}).
			AddRow(false))
	res, err = r.DeleteGroup(userID, groupID)
	assert.False(t, res)
	assert.NoError(t, err)

	/* User not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.DeleteGroup(userID, groupID)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.False(t, res)

	/* Group not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with ID"})
	res, err = r.DeleteGroup(userID, groupID)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.False(t, res)

	/* Deadline (5s) exceeded.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.DeleteGroup(userID, groupID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.False(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID).
		WillReturnError(&pq.Error{})
	res, err = r.DeleteGroup(userID, groupID)
	assert.Error(t, err)
	assert.False(t, res)
}
