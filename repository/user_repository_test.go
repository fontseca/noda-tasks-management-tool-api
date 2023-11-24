package repository

import (
	"database/sql"
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

func TestUserRepository_FetchShallowUserByID(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		res   *transfer.User
		err   error
		user  = &transfer.User{ID: uuid.MustParse(userID)}
		query = regexp.QuoteMeta(`
	  SELECT "user_id" AS "id",
	         "role_id" AS "role",
	         "first_name",
	         "middle_name",
	         "last_name",
	         "surname",
	         "picture_url",
	         "email",
	  			 "is_blocked",
	         "created_at",
	         "updated_at"
	    FROM fetch_user_by_id ($1);`)
	)

	t.Run("success", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"id", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "is_blocked", "created_at", "updated_at"}).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt)
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(rows)
		res, err = r.FetchShallowUserByID(userID)
		assert.Equal(t, user, res)
		assert.NoError(t, err)
	})

	t.Run("got a scanning error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"unknown_column", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "is_blocked", "created_at", "updated_at"}))
		res, err = r.FetchShallowUserByID(userID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("got an expected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.FetchShallowUserByID(userID)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Nil(t, res)
	})

	t.Run("got an unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(new(pq.Error))
		res, err = r.FetchShallowUserByID(userID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_FetchByEmail(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		res   *model.User
		err   error
		email = "foo@bar.com"
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
	    FROM fetch_user_by_email ($1);`)
	)

	t.Run("success", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"id", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password", "is_blocked", "created_at", "updated_at"}).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.Password, user.IsBlocked, user.CreatedAt, user.UpdatedAt)
		mock.
			ExpectQuery(query).
			WithArgs(email).
			WillReturnRows(rows)
		res, err = r.FetchByEmail(email)
		assert.Equal(t, user, res)
		assert.NoError(t, err)
	})

	t.Run("got a scanning error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"unknown_column", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password", "is_blocked", "created_at", "updated_at"}))
		res, err = r.FetchByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("got an expected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(email).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user email"})
		res, err = r.FetchByEmail(email)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Nil(t, res)
	})

	t.Run("got an unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(email).
			WillReturnError(new(pq.Error))
		res, err = r.FetchByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_FetchShallowUserByEmail(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		res   *transfer.User
		err   error
		email = "foo@bar.com"
		user  = &transfer.User{ID: uuid.MustParse(userID)}
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
	    FROM fetch_user_by_email ($1);`)
	)

	t.Run("success", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"id", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password", "is_blocked", "created_at", "updated_at"}).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, "password", user.IsBlocked, user.CreatedAt, user.UpdatedAt)
		mock.
			ExpectQuery(query).
			WithArgs(email).
			WillReturnRows(rows)
		res, err = r.FetchShallowUserByEmail(email)
		assert.Equal(t, user, res)
		assert.NoError(t, err)
	})

	t.Run("got a scanning error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"unknown_column", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "password", "is_blocked", "created_at", "updated_at"}))
		res, err = r.FetchShallowUserByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("got an expected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(email).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user email"})
		res, err = r.FetchShallowUserByEmail(email)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Nil(t, res)
	})

	t.Run("got an unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(email).
			WillReturnError(new(pq.Error))
		res, err = r.FetchShallowUserByEmail(email)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_Fetch(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	var (
		r         = NewUserRepository(db)
		res       []*transfer.User
		err       error
		page, rpp int64
		needle    = "foo"
		sortExpr  = "-first_name"
		user      = &transfer.User{ID: uuid.MustParse(userID)}
		query     = regexp.QuoteMeta(`
  	SELECT "user_id" AS "id",
  	       "role_id" AS "role",
  	       "first_name",
  	       "middle_name",
  	       "last_name",
  	       "surname",
  	       "picture_url",
  	       "email",
  	       "is_blocked",
  	       "created_at",
  	       "updated_at"
      FROM fetch_users ($1, $2, $3, $4);`)
	)

	t.Run("success with rpp=2", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"id", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "is_blocked", "created_at", "updated_at"}).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt)
		rpp = 2
		mock.
			ExpectQuery(query).
			WithArgs(page, rpp, needle, sortExpr).
			WillReturnRows(rows)
		res, err = r.Fetch(page, rpp, needle, sortExpr)
		assert.Len(t, res, 2)
		assert.NoError(t, err)
	})

	t.Run("success with rpp=3", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"id", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "is_blocked", "created_at", "updated_at"}).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt)
		rpp = 3
		mock.
			ExpectQuery(query).
			WithArgs(page, rpp, needle, sortExpr).
			WillReturnRows(rows)
		res, err = r.Fetch(page, rpp, needle, sortExpr)
		assert.Len(t, res, 3)
		assert.NoError(t, err)
	})

	t.Run("got a scanning error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(page, rpp, needle, sortExpr).
			WillReturnRows(
				sqlmock.NewRows([]string{"unknown_column", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "is_blocked", "created_at", "updated_at"}).
					AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt))
		res, err = r.Fetch(page, rpp, needle, sortExpr)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("got an unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(page, rpp, needle, sortExpr).
			WillReturnError(new(pq.Error))
		res, err = r.Fetch(page, rpp, needle, sortExpr)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_FetchBlocked(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	var (
		r         = NewUserRepository(db)
		res       []*transfer.User
		err       error
		page, rpp int64
		needle    = "foo"
		sortExpr  = "-first_name"
		user      = &transfer.User{ID: uuid.MustParse(userID)}
		query     = regexp.QuoteMeta(`
  	SELECT "user_id" AS "id",
  	       "role_id" AS "role",
  	       "first_name",
  	       "middle_name",
  	       "last_name",
  	       "surname",
  	       "picture_url",
  	       "email",
  	       "is_blocked",
  	       "created_at",
  	       "updated_at"
      FROM fetch_blocked_users ($1, $2, $3, $4);`)
	)

	t.Run("success with rpp=2", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"id", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "is_blocked", "created_at", "updated_at"}).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt)
		rpp = 2
		mock.
			ExpectQuery(query).
			WithArgs(page, rpp, needle, sortExpr).
			WillReturnRows(rows)
		res, err = r.FetchBlocked(page, rpp, needle, sortExpr)
		assert.Len(t, res, 2)
		assert.NoError(t, err)
	})

	t.Run("success with rpp=3", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"id", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "is_blocked", "created_at", "updated_at"}).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt).
			AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt)
		rpp = 3
		mock.
			ExpectQuery(query).
			WithArgs(page, rpp, needle, sortExpr).
			WillReturnRows(rows)
		res, err = r.FetchBlocked(page, rpp, needle, sortExpr)
		assert.Len(t, res, 3)
		assert.NoError(t, err)
	})

	t.Run("got a scanning error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(page, rpp, needle, sortExpr).
			WillReturnRows(
				sqlmock.NewRows([]string{"unknown_column", "role", "first_name", "middle_name", "last_name", "surname", "picture_url", "email", "is_blocked", "created_at", "updated_at"}).
					AddRow(user.ID, user.Role, user.FirstName, user.MiddleName, user.LastName, user.Surname, user.PictureUrl, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt))
		res, err = r.FetchBlocked(page, rpp, needle, sortExpr)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("got an unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(page, rpp, needle, sortExpr).
			WillReturnError(new(pq.Error))
		res, err = r.FetchBlocked(page, rpp, needle, sortExpr)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_FetchSettings(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	var (
		r         = NewUserRepository(db)
		res       []*transfer.UserSetting
		err       error
		page, rpp int64
		needle    = "foo"
		sortExpr  = "-first_name"
		setting   = new(transfer.UserSetting)
		query     = regexp.QuoteMeta("SELECT * FROM fetch_user_settings ($1, $2, $3, $4, $5);")
	)

	t.Run("success with rpp=2", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"key", "description", "value", "created_at", "updated_at"}).
			AddRow(setting.Key, setting.Description, setting.Value, setting.CreatedAt, setting.UpdatedAt).
			AddRow(setting.Key, setting.Description, setting.Value, setting.CreatedAt, setting.UpdatedAt)
		rpp = 2
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortExpr).
			WillReturnRows(rows)
		res, err = r.FetchSettings(userID, page, rpp, needle, sortExpr)
		assert.Len(t, res, 2)
		assert.NoError(t, err)
	})

	t.Run("success with rpp=3", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"key", "description", "value", "created_at", "updated_at"}).
			AddRow(setting.Key, setting.Description, setting.Value, setting.CreatedAt, setting.UpdatedAt).
			AddRow(setting.Key, setting.Description, setting.Value, setting.CreatedAt, setting.UpdatedAt).
			AddRow(setting.Key, setting.Description, setting.Value, setting.CreatedAt, setting.UpdatedAt)
		rpp = 3
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortExpr).
			WillReturnRows(rows)
		res, err = r.FetchSettings(userID, page, rpp, needle, sortExpr)
		assert.Len(t, res, 3)
		assert.NoError(t, err)
	})

	t.Run("got a scanning error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortExpr).
			WillReturnRows(
				sqlmock.NewRows([]string{"unknown_column", "description", "value", "created_at", "updated_at"}).
					AddRow(setting.Key, setting.Description, setting.Value, setting.CreatedAt, setting.UpdatedAt))
		res, err = r.FetchSettings(userID, page, rpp, needle, sortExpr)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("got not found user error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortExpr).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.FetchSettings(userID, page, rpp, needle, sortExpr)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Nil(t, res)
	})

	t.Run("got an unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, page, rpp, needle, sortExpr).
			WillReturnError(new(pq.Error))
		res, err = r.FetchSettings(userID, page, rpp, needle, sortExpr)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_FetchOneSetting(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	var (
		r       = NewUserRepository(db)
		res     *transfer.UserSetting
		err     error
		setting = &transfer.UserSetting{
			Key:   "key",
			Value: "value",
		}
		query = regexp.QuoteMeta("SELECT * FROM fetch_one_user_setting ($1, $2);")
	)

	t.Run("success", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"key", "description", "value", "created_at", "updated_at"}).
			AddRow(setting.Key, setting.Description, setting.Value, setting.CreatedAt, setting.UpdatedAt)
		mock.
			ExpectQuery(query).
			WithArgs(userID, setting.Key).
			WillReturnRows(rows)
		res, err = r.FetchOneSetting(userID, setting.Key)
		assert.Equal(t, setting, res)
		assert.NoError(t, err)
	})

	t.Run("got a scanning error", func(t *testing.T) {
		var rows = sqlmock.
			NewRows([]string{"unknown_column", "description", "value", "created_at", "updated_at"}).
			AddRow(setting.Key, setting.Description, setting.Value, setting.CreatedAt, setting.UpdatedAt)
		mock.
			ExpectQuery(query).
			WithArgs(userID, setting.Key).
			WillReturnRows(rows)
		res, err = r.FetchOneSetting(userID, setting.Key)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("got not found user error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, setting.Key).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.FetchOneSetting(userID, setting.Key)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Nil(t, res)
	})

	t.Run("got not found user setting error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, setting.Key).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent predefined user setting key"})
		res, err = r.FetchOneSetting(userID, setting.Key)
		assert.ErrorIs(t, err, noda.ErrSettingNotFound)
		assert.Nil(t, res)
	})

	t.Run("got an unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, setting.Key).
			WillReturnError(new(pq.Error))
		res, err = r.FetchOneSetting(userID, setting.Key)
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

func TestUserRepository_UpdateUserSetting(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r            = NewUserRepository(db)
		res          bool
		err          error
		query        = regexp.QuoteMeta("SELECT update_user_setting ($1, $2, $3);")
		settingKey   = "key"
		settingValue = "new value"
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, settingKey, settingValue).
			WillReturnRows(sqlmock.NewRows([]string{"update_user_setting"}).AddRow(true))
		res, err = r.UpdateUserSetting(userID, settingKey, settingValue)
		assert.NoError(t, err)
		assert.Equal(t, res, true)
	})

	t.Run("could not update but didn't get any error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, settingKey, settingValue).
			WillReturnRows(sqlmock.NewRows([]string{"update_user_setting"}).AddRow(false))
		res, err = r.UpdateUserSetting(userID, settingKey, settingValue)
		assert.NoError(t, err)
		assert.Equal(t, res, false)
	})

	t.Run("got not found user error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, settingKey, settingValue).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.UpdateUserSetting(userID, settingKey, settingValue)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Equal(t, res, false)
	})

	t.Run("got not found user setting error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, settingKey, settingValue).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent predefined user setting key"})
		res, err = r.UpdateUserSetting(userID, settingKey, settingValue)
		assert.ErrorIs(t, err, noda.ErrSettingNotFound)
		assert.Equal(t, res, false)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, settingKey, settingValue).
			WillReturnError(&pq.Error{})
		res, err = r.UpdateUserSetting(userID, settingKey, settingValue)
		assert.Error(t, err)
		assert.Equal(t, res, false)
	})
}

func TestUserRepository_Block(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		res   bool
		err   error
		query = regexp.QuoteMeta("SELECT block_user ($1);")
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"block_user"}).AddRow(true))
		res, err = r.Block(userID)
		assert.NoError(t, err)
		assert.Equal(t, res, true)
	})

	t.Run("could not block but didn't get any error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"block_user"}).AddRow(false))
		res, err = r.Block(userID)
		assert.NoError(t, err)
		assert.Equal(t, res, false)
	})

	t.Run("got not found user error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.Block(userID)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Equal(t, res, false)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{})
		res, err = r.Block(userID)
		assert.Error(t, err)
		assert.Equal(t, res, false)
	})
}

func TestUserRepository_Unblock(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		res   bool
		err   error
		query = regexp.QuoteMeta("SELECT unblock_user ($1);")
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"unblock_user"}).AddRow(true))
		res, err = r.Unblock(userID)
		assert.NoError(t, err)
		assert.Equal(t, res, true)
	})

	t.Run("could not unblock but didn't get any error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"unblock_user"}).AddRow(false))
		res, err = r.Unblock(userID)
		assert.NoError(t, err)
		assert.Equal(t, res, false)
	})

	t.Run("got not found user error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.Unblock(userID)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Equal(t, res, false)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{})
		res, err = r.Unblock(userID)
		assert.Error(t, err)
		assert.Equal(t, res, false)
	})
}

func TestUserRepository_PromoteToAdmin(t *testing.T) {
	defer beQuiet()()
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

	t.Run("could not promote but didn't get any error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"promote_user_to_admin"}).AddRow(false))
		res, err = r.PromoteToAdmin(userID)
		assert.NoError(t, err)
		assert.Equal(t, res, false)
	})

	t.Run("got not found user error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.PromoteToAdmin(userID)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Equal(t, res, false)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{})
		res, err = r.PromoteToAdmin(userID)
		assert.Error(t, err)
		assert.Equal(t, res, false)
	})
}

func TestUserRepository_DegradeToUser(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		query = regexp.QuoteMeta(`SELECT degrade_admin_to_user ($1);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.
				NewRows([]string{"degrade_admin_to_user"}).
				AddRow(true))
		res, err = r.DegradeToUser(userID)
		assert.NoError(t, err)
		assert.Equal(t, res, true)
	})

	t.Run("could not degrade but didn't get any error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"degrade_admin_to_user"}).AddRow(false))
		res, err = r.DegradeToUser(userID)
		assert.NoError(t, err)
		assert.Equal(t, res, false)
	})

	t.Run("got not found user error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		res, err = r.DegradeToUser(userID)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
		assert.Equal(t, res, false)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{})
		res, err = r.DegradeToUser(userID)
		assert.Error(t, err)
		assert.Equal(t, res, false)
	})
}

func TestUserRepository_RemoveHardly(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		query = regexp.QuoteMeta(`SELECT delete_user_hardly ($1);`)
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.
				NewRows([]string{"delete_user_hardly"}).
				AddRow(true))
		err = r.RemoveHardly(userID)
		assert.NoError(t, err)
	})

	t.Run("could not remove user", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{Code: "P0001", Message: "nonexistent user with ID"})
		err = r.RemoveHardly(userID)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{})
		err = r.RemoveHardly(userID)
		assert.Error(t, err)
	})
}

func TestUserRepository_RemoveSoftly(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewUserRepository(db)
		query = regexp.QuoteMeta(`
		DELETE FROM "user"
					WHERE "user_id" = $1;`)
		err error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnRows(sqlmock.
				NewRows([]string{"delete_user_hardly"}).
				AddRow(true))
		err = r.RemoveSoftly(userID)
		assert.NoError(t, err)
	})

	t.Run("could not remove user", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)
		err = r.RemoveSoftly(userID)
		assert.ErrorIs(t, err, noda.ErrUserNotFound)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID).
			WillReturnError(&pq.Error{})
		err = r.RemoveSoftly(userID)
		assert.Error(t, err)
	})
}
