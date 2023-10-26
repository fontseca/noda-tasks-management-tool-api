package repository

import (
	"errors"
	"github.com/google/uuid"
	"noda/api/data/model"
	"noda/api/data/transfer"
	"noda/failure"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const listID = "7d7b997f-a593-4ecd-a09f-039453321a51"

func TestListRepository_InsertList(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT make_list ($1, $2, $3, $4);`)
		res   string
		err   error
		next  = &transfer.ListCreation{Name: "list name", Description: "list desc"}
	)

	/* Success for grouped list.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, next.Name, next.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"make_list"}).
			AddRow(listID))
	res, err = r.InsertList(userID, groupID, next)
	assert.NoError(t, err)
	assert.Equal(t, listID, res)

	/* Success for scattered list.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, nil, next.Name, next.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"make_list"}).
			AddRow(listID))
	res, err = r.InsertList(userID, "", next)
	assert.NoError(t, err)
	assert.Equal(t, listID, res)

	/* User not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, next.Name, next.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.InsertList(userID, groupID, next)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Equal(t, "", res)

	/* Group not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, next.Name, next.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.InsertList(userID, groupID, next)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Equal(t, "", res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, next.Name, next.Description).
		WillReturnError(&pq.Error{})
	res, err = r.InsertList(userID, groupID, next)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func TestListRepository_FetchListByID(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT * FROM fetch_list_by_id ($1, $2, $3);`)
		res   *model.List
		err   error
		list  = &model.List{
			ID:          uuid.MustParse(listID),
			OwnerID:     uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			IsArchived:  false,
			ArchivedAt:  nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{"list_id", "owner_id", "group_id", "name", "description", "is_archived", "archived_at", "created_at", "updated_at"}
	)

	/* Success for grouped list.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchListByID(userID, groupID, listID)
	assert.NoError(t, err)
	assert.Equal(t, list, res)

	/* Success for scattered list.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, nil, listID).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchListByID(userID, "", listID)
	assert.NoError(t, err)
	assert.Equal(t, list, res)

	/* User was not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.FetchListByID(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Nil(t, res)

	/* Group was not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with ID"})
	res, err = r.FetchListByID(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.Nil(t, res)

	/* Grouped list not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent list with ID"})
	res, err = r.FetchListByID(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrListNotFound)
	assert.Nil(t, res)

	/* Scattered list not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, nil, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent list with ID"})
	res, err = r.FetchListByID(userID, "", listID)
	assert.Error(t, err, failure.ErrListNotFound)
	assert.Nil(t, res)

	/* Deadline (5s) exceeded.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.FetchListByID(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{})
	res, err = r.FetchListByID(userID, groupID, listID)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestListRepository_GetTodayListID(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT get_today_list_id ($1);`)
		res   string
		err   error
	)

	/* Success.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnRows(sqlmock.
			NewRows([]string{"get_today_list_id"}).
			AddRow("the actual ID"))
	res, err = r.GetTodayListID(userID)
	assert.NoError(t, err)
	assert.Equal(t, "the actual ID", res)

	/* User not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.GetTodayListID(userID)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Empty(t, res)

	/* Deadline (5s) exceeded.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.GetTodayListID(userID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Empty(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(&pq.Error{})
	res, err = r.GetTodayListID(userID)
	assert.Error(t, err)
	assert.Empty(t, res)
}

func TestListRepository_GetTomorrowListID(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT get_tomorrow_list_id ($1);`)
		res   string
		err   error
	)

	/* Success.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnRows(sqlmock.
			NewRows([]string{"get_tomorrow_list_id"}).
			AddRow("the actual ID"))
	res, err = r.GetTomorrowListID(userID)
	assert.NoError(t, err)
	assert.Equal(t, "the actual ID", res)

	/* User not found.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.GetTomorrowListID(userID)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Equal(t, "", res)

	/* Deadline (5s) exceeded.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.GetTomorrowListID(userID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Empty(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(&pq.Error{})
	res, err = r.GetTomorrowListID(userID)
	assert.Error(t, err)
	assert.Empty(t, res)
}

func TestListRepository_FetchLists(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`
		SELECT "list_id" AS "id",
		       "owner_id",
		       "group_id",
		       "name",
		       "description",
		       "is_archived",
		       "archived_at",
		       "created_at",
		       "updated_at"
     FROM fetch_lists ($1, $2, $3, $4, $5);`)
		res       []*model.List
		err       error
		page, rpp int64
		needle    = ""
		sortBy    = ""
		list      = &model.List{
			ID:          uuid.MustParse(listID),
			OwnerID:     uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			IsArchived:  false,
			ArchivedAt:  nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{"id", "owner_id", "group_id", "name", "description", "is_archived", "archived_at", "created_at", "updated_at"}
	)

	/* Success with 2 records.  */

	page, rpp = 1, 2
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchLists(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	/* Success with the default number of records (10).  */

	page, rpp = 1, -1000
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchLists(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 10)

	/* Success with custom pagination and RPP.  */

	page, rpp = 2, 5
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchLists(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 5)

	/* Success with searching.  */

	page, rpp, needle = 1, 7, "name"
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchLists(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 7)

	/* There should not be a response for a weird needle and neither should be
	   an error.  */

	page, rpp, needle = 1, 5, "aljfkjaksjpiwquramakjsfasjfkjwpoijefj"
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.NewRows(columns))
	res, err = r.FetchLists(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)

	/* User not found.  */

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.FetchLists(userID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Nil(t, res)

	/* Deadline (5s) exceeded.  */

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.FetchLists(userID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{})
	res, err = r.FetchLists(userID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)

	/* Unexpected scanning error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id", "unknown_column", "owner_id", "name", "description", "is_archived",
				"archived_at", "created_at", "updated_at"}).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchLists(userID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestListRepository_FetchGroupedLists(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`
		SELECT "list_id" AS "id",
		       "owner_id",
		       "group_id",
		       "name",
		       "description",
		       "is_archived",
		       "archived_at",
		       "created_at",
		       "updated_at"
      FROM fetch_grouped_lists ($1, $2, $3, $4, $5, $6);`)
		res       []*model.List
		err       error
		page, rpp int64
		needle    = ""
		sortBy    = ""
		list      = &model.List{
			ID:          uuid.MustParse(listID),
			OwnerID:     uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			IsArchived:  false,
			ArchivedAt:  nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{"id", "owner_id", "group_id", "name", "description", "is_archived", "archived_at", "created_at", "updated_at"}
	)

	/* Success with 2 records.  */

	page, rpp = 1, 2
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	/* Success with the default number of records (10).  */

	page, rpp = 1, -1000
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 10)

	/* Success with custom pagination and RPP.  */

	page, rpp = 2, 5
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 5)

	/* Success with searching.  */

	page, rpp, needle = 1, 7, "name"
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 7)

	/* There should not be a response for a weird needle and neither should be
	   an error.  */

	page, rpp, needle = 1, 5, "aljfkjaksjpiwquramakjsfasjfkjwpoijefj"
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.NewRows(columns))
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)

	/* User not found.  */

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Nil(t, res)

	/* Group not found.  */

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with ID"})
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.Nil(t, res)

	/* Deadline (5s) exceeded.  */

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{})
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)

	/* Unexpected scanning error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id", "unknown_column", "owner_id", "name", "description", "is_archived",
				"archived_at", "created_at", "updated_at"}).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGroupedLists(userID, groupID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestListRepository_FetchScatteredLists(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`
		SELECT "list_id" AS "id",
		       "owner_id",
		       "group_id",
		       "name",
		       "description",
		       "is_archived",
		       "archived_at",
		       "created_at",
		       "updated_at"
      FROM fetch_scattered_lists ($1, $2, $3, $4, $5);`)
		res       []*model.List
		err       error
		page, rpp int64
		needle    = ""
		sortBy    = ""
		list      = &model.List{
			ID:          uuid.MustParse(listID),
			OwnerID:     uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			IsArchived:  false,
			ArchivedAt:  nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{"id", "owner_id", "group_id", "name", "description", "is_archived", "archived_at", "created_at", "updated_at"}
	)

	/* Success with 2 records.  */

	page, rpp = 1, 2
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScatteredLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	/* Success with the default number of records (10).  */

	page, rpp = 1, -1000
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScatteredLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 10)

	/* Success with custom pagination and RPP.  */

	page, rpp = 2, 5
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScatteredLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 5)

	/* Success with searching.  */

	page, rpp, needle = 1, 7, "name"
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScatteredLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 7)

	/* There should not be a response for a weird needle and neither should be
	   an error.  */

	page, rpp, needle = 1, 5, "aljfkjaksjpiwquramakjsfasjfkjwpoijefj"
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.NewRows(columns))
	res, err = r.FetchScatteredLists(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)

	/* User not found.  */

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.FetchScatteredLists(userID, groupID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Nil(t, res)

	/* Deadline (5s) exceeded.  */

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.FetchScatteredLists(userID, groupID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{})
	res, err = r.FetchScatteredLists(userID, groupID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)

	/* Unexpected scanning error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows([]string{
				"id", "unknown_column", "owner_id", "name", "description", "is_archived",
				"archived_at", "created_at", "updated_at"}).
			AddRow(list.ID, list.OwnerID, list.GroupID, list.Name, list.Description, list.IsArchived, list.ArchivedAt, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScatteredLists(userID, groupID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)
}
