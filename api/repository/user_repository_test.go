package repository

import (
	"noda/api/data/transfer"
	"noda/failure"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const userID string = "9039f725-e31f-4f04-bdb1-7b74e7f72d59"

func TestInsertUser_ValidInput(t *testing.T) {
	db, mock := newMock()
	r := NewUserRepository(db)
	query := `SELECT make_user ($1, $2, $3, $4, $5, $6);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("FirstName", "MiddleName", "LastName", "Surname", "Email", "Password").
		WillReturnRows(sqlmock.
			NewRows([]string{"make_user"}).
			AddRow(userID))
	res, err := r.InsertUser(&transfer.UserCreation{
		FirstName:  "FirstName",
		MiddleName: "MiddleName",
		LastName:   "LastName",
		Surname:    "Surname",
		Email:      "Email",
		Password:   "Password",
	})
	assert.NoError(t, err)
	assert.Equal(t, res, userID)
}

func TestInsertUser_InvalidEmail(t *testing.T) {
	defer beQuiet()()
	n := transfer.UserCreation{Email: "invalid-email"}
	db, mock := newMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := `SELECT make_user ($1, $2, $3, $4, $5, $6);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
		WillReturnError(&pq.Error{Code: "23514", Message: "value for domain email_t violates check constraint \"email_t_check\""})
	res, err := r.InsertUser(&n)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func TestInsertUser_DuplicatedEmail(t *testing.T) {
	defer beQuiet()()
	n := transfer.UserCreation{Email: "mail@mail.com"}
	db, mock := newMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := `SELECT make_user ($1, $2, $3, $4, $5, $6);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
		WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint \"user_email_key\""})
	res, err := r.InsertUser(&n)
	assert.ErrorIs(t, err, failure.ErrSameEmail)
	assert.Equal(t, "", res)
}

func TestInsertUser_UnexpectedDatabaseError(t *testing.T) {
	defer beQuiet()()
	n := transfer.UserCreation{}
	db, mock := newMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := `SELECT make_user ($1, $2, $3, $4, $5, $6);`
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
		WillReturnError(&pq.Error{})
	res, err := r.InsertUser(&n)
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func TestUpdateUser_ValidInput(t *testing.T) {
	up := transfer.UserUpdate{}
	db, mock := newMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_user"}).
			AddRow(true))
	res, err := r.UpdateUser(userID, &up)
	assert.NoError(t, err)
	assert.Equal(t, res, true)
}

func TestUpdateUser_FailureWithoutError(t *testing.T) {
	up := transfer.UserUpdate{}
	db, mock := newMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_user"}).
			AddRow(false))
	res, err := r.UpdateUser(userID, &up)
	assert.NoError(t, err)
	assert.Equal(t, res, false)
}

func TestUpdateUser_NonexistentUserID(t *testing.T) {
	up := transfer.UserUpdate{}
	db, mock := newMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err := r.UpdateUser(userID, &up)
	assert.ErrorIs(t, err, failure.ErrNotFound)
	assert.Equal(t, res, false)
}

func TestUpdateUser_UnexpectedDatabaseError(t *testing.T) {
	defer beQuiet()()
	up := transfer.UserUpdate{}
	db, mock := newMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnError(&pq.Error{})
	res, err := r.UpdateUser(userID, &up)
	assert.Error(t, err)
	assert.Equal(t, res, false)
}

func TestPromoteUserToAdmin_ValidInput(t *testing.T) {
	db, mock := newMock()
	defer db.Close()
	r := NewUserRepository(db)
	query := "SELECT promote_user_to_admin ($1);"
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID).
		WillReturnRows(sqlmock.
			NewRows([]string{"promote_user_to_admin"}).
			AddRow(true))
	got, err := r.PromoteUserToAdmin(userID)
	assert.NoError(t, err)
	assert.Equal(t, got, true)
}
