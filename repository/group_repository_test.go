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
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`SELECT "groups"."make" ($1, $2, $3);`)
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
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
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
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`SELECT * FROM "groups"."fetch" (p_owner_uuid := $1,
                                                              p_group_uuid := $2,
                                                              p_needle := NULL,
                                                              p_page := NULL,
                                                              p_rpp := NULL);`)
		res   *model.Group
		err   error
		group = &model.Group{
			UUID:        uuid.MustParse(groupID),
			OwnerUUID:   uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			CreatedAt:   nil,
			UpdatedAt:   nil,
		}
		columns = []string{"uuid", "owner_uuid", "name", "description", "created_at", "updated_at"}
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnRows(sqlmock.
				NewRows(columns).
				AddRow(group.UUID, group.OwnerUUID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt))
		res, err = r.FetchByID(userID, groupID)
		assert.NoError(t, err)
		assert.Equal(t, group, res)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
		res, err = r.FetchByID(userID, groupID)
		assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
		assert.Nil(t, res)
	})

	t.Run("group not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnRows(sqlmock.NewRows(columns))
		res, err = r.FetchByID(userID, groupID)
		assert.ErrorIs(t, err, failure.ErrGroupNotFound)
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
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`
		SELECT "group_uuid" AS "uuid",
         "owner_uuid",
         "name",
         "description",
         "created_at",
         "updated_at"
	  FROM "groups"."fetch" (p_owner_uuid := $1,
	                         p_group_uuid := $2,
	                         p_needle := $3,
	                         p_page := $4,
	                         p_rpp := $5);`)
		res       []*model.Group
		err       error
		page, rpp int64
		needle    = "name"
		group     = model.Group{
			UUID:        uuid.New(),
			OwnerUUID:   uuid.MustParse(userID),
			Name:        "name",
			Description: "desc",
			CreatedAt:   nil,
			UpdatedAt:   nil,
		}
		columns = []string{"uuid", "owner_uuid", "name", "description", "created_at", "updated_at"}
	)

	t.Run("success", func(t *testing.T) {
		page, rpp = 1, 2
		mock.
			ExpectQuery(query).
			WithArgs(userID, nil, needle, page, rpp).
			WillReturnRows(sqlmock.
				NewRows(columns).
				AddRow(group.UUID, group.OwnerUUID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt).
				AddRow(group.UUID, group.OwnerUUID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt))
		res, err = r.Fetch(userID, page, rpp, needle, "")
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("user not found", func(t *testing.T) {
		page, rpp = 1, 10
		mock.
			ExpectQuery(query).
			WithArgs(userID, nil, needle, page, rpp).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
		res, err = r.Fetch(userID, page, rpp, needle, "")
		assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
		assert.Nil(t, res)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, nil, needle, page, rpp).
			WillReturnError(&pq.Error{})
		res, err = r.Fetch(userID, page, rpp, needle, "")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("unexpected scanning error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, nil, needle, page, rpp).
			WillReturnRows(sqlmock.
				NewRows([]string{"group_uuid", "owner_id", "name", "description", "created_at", "updated_at"}).
				AddRow(group.UUID, group.OwnerUUID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt))
		res, err = r.Fetch(userID, page, rpp, needle, "")
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGroupRepository_Update(t *testing.T) {
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`SELECT "groups"."update" ($1, $2, $3, $4);`)
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
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
		res, err = r.Update(userID, groupID, up)
		assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
		assert.False(t, res)
	})

	t.Run("group not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID, up.Name, up.Description).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with UUID"})
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
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewGroupRepository(db)
		query = regexp.QuoteMeta(`SELECT "groups"."delete" ($1, $2);`)
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
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with UUID"})
		res, err = r.Remove(userID, groupID)
		assert.ErrorIs(t, err, failure.ErrUserNoLongerExists)
		assert.False(t, res)
	})

	t.Run("group not found", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, groupID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent group with UUID"})
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
