package repository

import (
	"github.com/google/uuid"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const userID string = "9039f725-e31f-4f04-bdb1-7b74e7f72d59"

func TestUserRepository_Save(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
			WillReturnRows(sqlmock.
				NewRows([]string{"make_user"}).
				AddRow(userID))
		res, err = r.Save(n)
		assert.NoError(t, err)
		assert.Equal(t, res, userID)
	})

	t.Run("got an invalid email", func(t *testing.T) {
		n.Email = "invalid-email"
		mock.
			ExpectQuery(query).
			WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
			WillReturnError(&pq.Error{Code: "23514", Message: "value for domain email_t violates check constraint \"email_t_check\""})
		res, err = r.Save(n)
		assert.Error(t, err)
		assert.Equal(t, "", res)
	})

	t.Run("got a duplicated email", func(t *testing.T) {
		n.Email = "mail@mail.com"
		mock.
			ExpectQuery(query).
			WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
			WillReturnError(&pq.Error{Code: "23505", Message: "duplicate key value violates unique constraint \"user_email_key\""})
		res, err = r.Save(n)
		assert.ErrorIs(t, err, noda.ErrSameEmail)
		assert.Equal(t, "", res)
	})

	t.Run("got an unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(n.FirstName, n.MiddleName, n.LastName, n.Surname, n.Email, n.Password).
			WillReturnError(&pq.Error{})
		res, err = r.Save(n)
		assert.Error(t, err)
		assert.Equal(t, "", res)
	})
}

func TestUserRepository_FetchByID(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		res   *model.User
		err   error
		user  = &model.User{ID: uuid.MustParse(userID)}
		query = regexp.QuoteMeta(`
	  SELECT "user_id" AS "id",
	         "role_id" AS "role",
	         "first_name",
	         "middle_name",
	         "last_name",
	         "surname",
	         "picture_url",
	         "email",
	  			 "password",
	  			 "is_blocked",
	         "created_at",
	         "updated_at"
	    FROM fetch_user_by_id ($1);`)
	)

	t.Run("success", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"id", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password", "is_blocked", "created_at", "updated_at"}).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.Password, user.IsBlocked, user.CreatedAt, user.UpdatedAt)
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(rows)
		res, err = r.FetchByID(userID)
		assert.Equal(t, user, res)
		assert.NoError(t, err)
	})

	t.Run("got a scanning error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"unknown_column", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password", "is_blocked", "created_at", "updated_at"}))
		res, err = r.FetchByID(userID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("got an expected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.FetchByID(userID)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Nil(t, res)
	})

	t.Run("got an unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(new(pq.Error))
		res, err = r.FetchByID(userID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_Update(t *testing.T) {
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

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
			WillReturnRows(sqlmock.
				NewRows([]string{"update_user"}).
				AddRow(true))
		res, err = r.Update(userID, up)
		assert.NoError(t, err)
		assert.Equal(t, res, true)
	})

	t.Run("could not update but didn't get any error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
			WillReturnRows(sqlmock.
				NewRows([]string{"update_user"}).
				AddRow(false))
		res, err = r.Update(userID, up)
		assert.NoError(t, err)
		assert.Equal(t, res, false)
	})

	t.Run("user does not exist", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.Update(userID, up)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Equal(t, res, false)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, up.FirstName, up.MiddleName, up.LastName, up.Surname).
			WillReturnError(&pq.Error{})
		res, err = r.Update(userID, up)
		assert.Error(t, err)
		assert.Equal(t, res, false)
	})
}

func TestUserRepository_PromoteToAdmin(t *testing.T) {
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		query = regexp.QuoteMeta(`SELECT promote_user_to_admin ($1);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.
				NewRows([]string{"promote_user_to_admin"}).
				AddRow(true))
		res, err = r.PromoteToAdmin(userID)
		assert.NoError(t, err)
		assert.Equal(t, res, true)
	})
}
