package repository

import (
	"noda/api/data/transfer"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const groupID string = "942d76f4-28b2-44be-8339-232b62c0ef22"

func TestInsertGroup_ValidInput(t *testing.T) {
	next := &transfer.GroupCreation{
		Name:        "name",
		Description: "desc",
	}
	db, mock := newMock()
	r := NewGroupRepository(db)
	query := `SELECT make_group ($1, $2, $3);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, next.Name, next.Description).
		WillReturnRows(sqlmock.
			NewRows([]string{"make_group"}).
			AddRow(groupID))
	res, err := r.InsertGroup(userID, next)
	assert.NoError(t, err)
	assert.Equal(t, groupID, res)
}

func TestInsertGroup_UserNotFound(t *testing.T) {
	next := &transfer.GroupCreation{
		Name:        "name",
		Description: "desc",
	}
	db, mock := newMock()
	r := NewGroupRepository(db)
	query := `SELECT make_group ($1, $2, $3);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, next.Name, next.Description).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err := r.InsertGroup(userID, next)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func TestInsertGroup_UnexpectedDatabaseError(t *testing.T) {
	defer beQuiet()()
	next := &transfer.GroupCreation{
		Name:        "name",
		Description: "desc",
	}
	db, mock := newMock()
	r := NewGroupRepository(db)
	query := `SELECT make_group ($1, $2, $3);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, next.Name, next.Description).
		WillReturnError(&pq.Error{})
	res, err := r.InsertGroup(userID, next)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}
