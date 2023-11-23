package repository

import (
	"noda"
	"noda/data/transfer"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const userID string = "9039f725-e31f-4f04-bdb1-7b74e7f72d59"

func TestGroupRepository_Save(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		query = regexp.QuoteMeta(`SELECT make_user ($1, $2, $3, $4, $5, $6);`)
		res   string
		err   error
		n     = &transfer.UserCreation{
			FirstName:  "FirstName",
			MiddleName: "MiddleName",
			LastName:   "LastName",
			Surname:    "Surname",
			Email:      "Email",
			Password:   "Password",
		}
	)

	/* Success.  */

	mock.
		ExpectQuery(query).
		WithArgs("FirstName", "MiddleName", "LastName", "Surname", "Email", "Password").
		WillReturnRows(sqlmock.
			NewRows([]string{"make_user"}).
			AddRow(userID))
	res, err = r.Save(n)
	assert.NoError(t, err)
	assert.Equal(t, res, userID)

	/* Invalid email.  */

	n.Email = "invalid-email"
	mock.
		ExpectQuery(query).
		WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
		WillReturnError(&pq.Error{Code: "23514", Message: "value for domain email_t violates check constraint \"email_t_check\""})
	res, err = r.Save(n)
	assert.Error(t, err)
	assert.Equal(t, "", res)

	/* Duplicated email.  */

	n.Email = "mail@mail.com"
	mock.
		ExpectQuery(query).
		WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
		WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint \"user_email_key\""})
	res, err = r.Save(n)
	assert.ErrorIs(t, err, noda.ErrSameEmail)
	assert.Equal(t, "", res)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
		WillReturnError(&pq.Error{})
	res, err = r.Save(n)
	assert.Error(t, err)
	assert.Equal(t, "", res)

}

func TestGroupRepository_Update(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		res   bool
		err   error
		query = regexp.QuoteMeta(`SELECT update_user ($1, $2, $3, $4, $5, NULL, NULL, NULL);`)
		up    = &transfer.UserUpdate{}
	)

	/* Success.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_user"}).
			AddRow(true))
	res, err = r.Update(userID, up)
	assert.NoError(t, err)
	assert.Equal(t, res, true)

	/* Could not update but didn't get any error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnRows(sqlmock.
			NewRows([]string{"update_user"}).
			AddRow(false))
	res, err = r.Update(userID, up)
	assert.NoError(t, err)
	assert.Equal(t, res, false)

	/* User does not exist.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
	res, err = r.Update(userID, up)
	assert.ErrorIs(t, err, noda.ErrUserNotFound)
	assert.Equal(t, res, false)

	/* Unexpected database error.  */

	mock.
		ExpectQuery(query).
		WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
		WillReturnError(&pq.Error{})
	res, err = r.Update(userID, up)
	assert.Error(t, err)
	assert.Equal(t, res, false)
}

func TestGroupRepository_PromoteToAdmin(t *testing.T) {
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		query = regexp.QuoteMeta(`SELECT promote_user_to_admin ($1);`)
		res   bool
		err   error
	)
	mock.
		ExpectQuery(query).
		WithArgs(userID).
		WillReturnRows(sqlmock.
			NewRows([]string{"promote_user_to_admin"}).
			AddRow(true))
	res, err = r.PromoteToAdmin(userID)
	assert.NoError(t, err)
	assert.Equal(t, res, true)
}
