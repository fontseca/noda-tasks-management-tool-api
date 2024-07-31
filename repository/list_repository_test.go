package repository

import (
	"errors"
	"github.com/google/uuid"
	"noda/data/model"
	"noda/data/transfer"
	"noda/failure"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const listID = "7d7b997f-a593-4ecd-a09f-039453321a51"

func TestListRepository_Save(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT "lists"."make" ($1, $2, $3, $4);`)
		res   string
		err   error
		next  = &transfer.ListCreation{Name: "list name", Description: "list desc"}
	)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, next.Name, next.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"make_list"}).
			AddRow(listID))
	res, err = r.Save(userID, groupID, next)
	assert.NoError(t, err)
	assert.Equal(t, listID, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, nil, next.Name, next.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"make_list"}).
			AddRow(listID))
	res, err = r.Save(userID, "", next)
	assert.NoError(t, err)
	assert.Equal(t, listID, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, next.Name, next.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.Save(userID, groupID, next)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.Equal(t, "", res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, next.Name, next.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.Save(userID, groupID, next)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.Equal(t, "", res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, next.Name, next.Description).
		WillReturnError(&pq.Error{})
	res, err = r.Save(userID, groupID, next)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func TestListRepository_FetchByID(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT * FROM "lists"."fetch_by_uuid" ($1, $2, $3);`)
		res   *model.List
		err   error
		list  = &model.List{
			UUID:        uuid.MustParse(listID),
			OwnerUUID:   uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{"list_id", "owner_id", "group_uuid", "name", "description", "created_at", "updated_at"}
	)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchByID(userID, groupID, listID)
	assert.NoError(t, err)
	assert.Equal(t, list, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, nil, listID).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchByID(userID, "", listID)
	assert.NoError(t, err)
	assert.Equal(t, list, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.FetchByID(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with UUID"})
	res, err = r.FetchByID(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent list with UUID"})
	res, err = r.FetchByID(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrListNotFound)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, nil, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent list with UUID"})
	res, err = r.FetchByID(userID, "", listID)
	assert.Error(t, err, failure.ErrListNotFound)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.FetchByID(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{})
	res, err = r.FetchByID(userID, groupID, listID)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestListRepository_GetTodayListID(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT "lists"."get_today_list_uuid" ($1);`)
		res   string
		err   error
	)

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnRows(sqlmock.
			NewRows([]string{"get_today_list_id"}).
			AddRow("the actual UUID"))
	res, err = r.GetTodayListID(userID)
	assert.NoError(t, err)
	assert.Equal(t, "the actual UUID", res)

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.GetTodayListID(userID)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.Empty(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.GetTodayListID(userID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Empty(t, res)

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
		query = regexp.QuoteMeta(`SELECT "lists"."get_tomorrow_list_uuid" ($1);`)
		res   string
		err   error
	)

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnRows(sqlmock.
			NewRows([]string{"get_tomorrow_list_id"}).
			AddRow("the actual UUID"))
	res, err = r.GetTomorrowListID(userID)
	assert.NoError(t, err)
	assert.Equal(t, "the actual UUID", res)

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.GetTomorrowListID(userID)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.Equal(t, "", res)

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.GetTomorrowListID(userID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Empty(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnError(&pq.Error{})
	res, err = r.GetTomorrowListID(userID)
	assert.Error(t, err)
	assert.Empty(t, res)
}

func TestListRepository_Fetch(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`
		SELECT "list_uuid" AS "uuid",
		       "owner_uuid",
		       "group_uuid",
		       "name",
		       "description",
		       "created_at",
		       "updated_at"
      FROM "lists"."fetch" ($1, $2, $3, $4, $5);`)
		res       []*model.List
		err       error
		page, rpp int64
		needle    = ""
		sortBy    = ""
		list      = &model.List{
			UUID:        uuid.MustParse(listID),
			OwnerUUID:   uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{"uuid", "owner_uuid", "group_uuid", "name", "description", "created_at", "updated_at"}
	)

	page, rpp = 1, 2
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.Fetch(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	page, rpp = 1, -1000
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.Fetch(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 10)

	page, rpp = 2, 5
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.Fetch(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 5)

	page, rpp, needle = 1, 7, "name"
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.Fetch(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 7)

	page, rpp, needle = 1, 5, "aljfkjaksjpiwquramakjsfasjfkjwpoijefj"
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.NewRows(columns))
	res, err = r.Fetch(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.Fetch(userID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.Nil(t, res)

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.Fetch(userID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{})
	res, err = r.Fetch(userID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "unknown_column", "owner_id", "name", "description", "created_at", "updated_at"}).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.Fetch(userID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestListRepository_FetchGrouped(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`
		SELECT "list_uuid" AS "uuid",
           "owner_uuid",
           "group_uuid",
		       "name",
		       "description",
		       "created_at",
		       "updated_at"
      FROM "lists"."fetch_grouped" ($1, $2, $3, $4, $5, $6);`)
		res       []*model.List
		err       error
		page, rpp int64
		needle    = ""
		sortBy    = ""
		list      = &model.List{
			UUID:        uuid.MustParse(listID),
			OwnerUUID:   uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{"uuid", "owner_uuid", "group_uuid", "name", "description", "created_at", "updated_at"}
	)

	page, rpp = 1, 2
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	page, rpp = 1, -1000
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 10)

	page, rpp = 2, 5
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 5)

	page, rpp, needle = 1, 7, "name"
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 7)

	page, rpp, needle = 1, 5, "aljfkjaksjpiwquramakjsfasjfkjwpoijefj"
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.NewRows(columns))
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.Nil(t, res)

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with UUID"})
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.Nil(t, res)

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{})
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "unknown_column", "owner_id", "name", "description", "created_at", "updated_at"}).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchGrouped(userID, groupID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestListRepository_FetchScattered(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`
		SELECT "list_uuid" AS "uuid",
           "owner_uuid",
           "group_uuid",
		       "name",
		       "description",
		       "created_at",
		       "updated_at"
      FROM "lists"."fetch_scattered" ($1, $2, $3, $4, $5);`)
		res       []*model.List
		err       error
		page, rpp int64
		needle    = ""
		sortBy    = ""
		list      = &model.List{
			UUID:        uuid.MustParse(listID),
			OwnerUUID:   uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		columns = []string{"uuid", "owner_uuid", "group_uuid", "name", "description", "created_at", "updated_at"}
	)

	/* Success with 2 records.  */

	page, rpp = 1, 2
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScattered(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	/* Success with the default number of records (10).  */

	page, rpp = 1, -1000
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScattered(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 10)

	/* Success with custom pagination and RPP.  */

	page, rpp = 2, 5
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScattered(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 5)

	page, rpp, needle = 1, 7, "name"
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows(columns).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScattered(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.Len(t, res, 7)

	page, rpp, needle = 1, 5, "some random text"
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.NewRows(columns))
	res, err = r.FetchScattered(userID, page, rpp, needle, sortBy)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 0)

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.FetchScattered(userID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.Nil(t, res)

	page, rpp = 1, 10
	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.FetchScattered(userID, page, rpp, needle, sortBy)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnError(&pq.Error{})
	res, err = r.FetchScattered(userID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, page, rpp, needle, sortBy).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "unknown_column", "owner_id", "name", "description", "created_at", "updated_at"}).
			AddRow(list.UUID, list.OwnerUUID, list.GroupUUID, list.Name, list.Description, list.CreatedAt, list.UpdatedAt))
	res, err = r.FetchScattered(userID, page, rpp, needle, sortBy)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestListRepository_Remove(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT "lists"."delete" ($1, $2, $3);`)
		res   bool
		err   error
	)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnRows(sqlmock.
			NewRows([]string{"delete_list"}).
			AddRow(true))
	res, err = r.Remove(userID, groupID, listID)
	assert.True(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, nil, listID).
		WillReturnRows(sqlmock.
			NewRows([]string{"delete_list"}).
			AddRow(true))
	res, err = r.Remove(userID, "", listID)
	assert.True(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnRows(sqlmock.
			NewRows([]string{"delete_list"}).
			AddRow(false))
	res, err = r.Remove(userID, groupID, listID)
	assert.False(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.Remove(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with UUID"})
	res, err = r.Remove(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent list with UUID"})
	res, err = r.Remove(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrListNotFound)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.Remove(userID, groupID, listID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID).
		WillReturnError(&pq.Error{})
	res, err = r.Remove(userID, groupID, listID)
	assert.Error(t, err)
	assert.False(t, res)
}

func TestListRepository_Duplicate(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r         = NewListRepository(db)
		query     = regexp.QuoteMeta(`SELECT "lists"."duplicate" ($1, $2);`)
		replicaID = uuid.New().String()
		res       string
		err       error
	)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnRows(sqlmock.
			NewRows([]string{"duplicate_list"}).
			AddRow(replicaID))
	res, err = r.Duplicate(userID, listID)
	assert.Equal(t, replicaID, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.Duplicate(userID, listID)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.Empty(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent list with UUID"})
	res, err = r.Duplicate(userID, listID)
	assert.ErrorIs(t, err, failure.ErrListNotFound)
	assert.Empty(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.Duplicate(userID, listID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.Empty(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnError(&pq.Error{})
	res, err = r.Duplicate(userID, listID)
	assert.Error(t, err)
	assert.Empty(t, res)
}

func TestListRepository_Scatter(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT "lists"."convert_to_scattered_list" ($1, $2);`)
		res   bool
		err   error
	)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnRows(sqlmock.
			NewRows([]string{"convert_to_scattered_list"}).
			AddRow(true))
	res, err = r.Scatter(userID, listID)
	assert.True(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.Scatter(userID, listID)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent list with UUID"})
	res, err = r.Scatter(userID, listID)
	assert.False(t, res)
	assert.ErrorIs(t, err, failure.ErrListNotFound)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnRows(sqlmock.
			NewRows([]string{"convert_to_scattered_list"}).
			AddRow(false))
	res, err = r.Scatter(userID, listID)
	assert.False(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.Scatter(userID, listID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID).
		WillReturnError(&pq.Error{})
	res, err = r.Scatter(userID, listID)
	assert.Error(t, err)
	assert.False(t, res)
}

func TestListRepository_Move(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT "lists"."move" ($1, $2, $3);`)
		res   bool
		err   error
	)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID, groupID).
		WillReturnRows(sqlmock.
			NewRows([]string{"move_list"}).
			AddRow(true))
	res, err = r.Move(userID, listID, groupID)
	assert.True(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID, groupID).
		WillReturnRows(sqlmock.
			NewRows([]string{"move_list"}).
			AddRow(false))
	res, err = r.Move(userID, listID, groupID)
	assert.False(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID, groupID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.Move(userID, listID, groupID)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID, groupID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent list with UUID"})
	res, err = r.Move(userID, listID, groupID)
	assert.False(t, res)
	assert.ErrorIs(t, err, failure.ErrListNotFound)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID, groupID).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with UUID"})
	res, err = r.Move(userID, listID, groupID)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID, groupID).
		WillReturnRows(sqlmock.
			NewRows([]string{"move_list"}).
			AddRow(false))
	res, err = r.Move(userID, listID, groupID)
	assert.False(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID, groupID).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.Move(userID, listID, groupID)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, listID, groupID).
		WillReturnError(&pq.Error{})
	res, err = r.Move(userID, listID, groupID)
	assert.Error(t, err)
	assert.False(t, res)
}

func TestListRepository_Update(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewListRepository(db)
		query = regexp.QuoteMeta(`SELECT "lists"."update" ($1, $2, $3, $4, $5);`)
		res   bool
		err   error
		up    = new(transfer.ListUpdate)
	)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID, up.Name, up.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_list"}).
			AddRow(true))
	res, err = r.Update(userID, groupID, listID, up)
	assert.True(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, nil, listID, up.Name, up.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_list"}).
			AddRow(true))
	res, err = r.Update(userID, "", listID, up)
	assert.True(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID, up.Name, up.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_list"}).
			AddRow(false))
	res, err = r.Update(userID, groupID, listID, up)
	assert.False(t, res)
	assert.NoError(t, err)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID, up.Name, up.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
	res, err = r.Update(userID, groupID, listID, up)
	assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID, up.Name, up.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with UUID"})
	res, err = r.Update(userID, groupID, listID, up)
	assert.ErrorIs(t, err, failure.ErrGroupNotFound)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID, up.Name, up.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent list with UUID"})
	res, err = r.Update(userID, groupID, listID, up)
	assert.ErrorIs(t, err, failure.ErrListNotFound)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID, up.Name, up.Description).
		WillReturnError(errors.New("context deadline exceeded"))
	res, err = r.Update(userID, groupID, listID, up)
	assert.ErrorIs(t, err, failure.ErrDeadlineExceeded)
	assert.False(t, res)

	mock.
		ExpectQuery(query).
		WithArgs(userID, groupID, listID, up.Name, up.Description).
		WillReturnError(new(pq.Error))
	res, err = r.Update(userID, groupID, listID, up)
	assert.Error(t, err)
	assert.False(t, res)
}
