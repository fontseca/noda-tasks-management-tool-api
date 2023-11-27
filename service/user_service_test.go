package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"testing"
)

type userRepositoryMock struct {
	mock.Mock
}

func newUserRepositoryMock() *userRepositoryMock {
	return new(userRepositoryMock)
}

func (o *userRepositoryMock) Save(creation *transfer.UserCreation) (insertedID string, err error) {
	var args = o.Called(creation)
	return args.String(0), args.Error(1)
}

func (o *userRepositoryMock) FetchByID(id string) (user *model.User, err error) {
	var args = o.Called(id)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*model.User)
	}
	return user, args.Error(1)
}

func (o *userRepositoryMock) FetchShallowUserByID(id string) (user *transfer.User, err error) {
	var args = o.Called(id)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*transfer.User)
	}
	return user, args.Error(1)
}

func (o *userRepositoryMock) FetchByEmail(email string) (user *model.User, err error) {
	var args = o.Called(email)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*model.User)
	}
	return user, args.Error(1)
}

func (o *userRepositoryMock) FetchShallowUserByEmail(email string) (user *transfer.User, err error) {
	var args = o.Called(email)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*transfer.User)
	}
	return user, args.Error(1)
}

func (o *userRepositoryMock) Fetch(page, rpp int64, needle, sortExpr string) (users []*transfer.User, err error) {
	var args = o.Called(page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		users = arg0.([]*transfer.User)
	}
	return users, args.Error(1)
}

func (o *userRepositoryMock) FetchBlocked(page, rpp int64, needle, sortExpr string) (users []*transfer.User, err error) {
	var args = o.Called(page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		users = arg0.([]*transfer.User)
	}
	return users, args.Error(1)
}

func (o *userRepositoryMock) FetchSettings(userID string, page, rpp int64, needle, sortExpr string) (settings []*transfer.UserSetting, err error) {
	var args = o.Called(userID, page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		settings = arg0.([]*transfer.UserSetting)
	}
	return settings, args.Error(1)
}

func (o *userRepositoryMock) FetchOneSetting(userID string, settingKey string) (setting *transfer.UserSetting, err error) {
	var args = o.Called(userID, settingKey)
	var arg0 = args.Get(0)
	if nil != arg0 {
		setting = arg0.(*transfer.UserSetting)
	}
	return setting, args.Error(1)
}

func (o *userRepositoryMock) Search(page, rpp int64, needle, sortExpr string) (users []*transfer.User, err error) {
	var args = o.Called(page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		users = arg0.([]*transfer.User)
	}
	return users, args.Error(1)
}

func (o *userRepositoryMock) Update(id string, update *transfer.UserUpdate) (ok bool, err error) {
	var args = o.Called(id, update)
	return args.Bool(0), args.Error(1)
}

func (o *userRepositoryMock) UpdateUserSetting(userID, settingKey, newValue string) (ok bool, err error) {
	var args = o.Called(userID, settingKey, newValue)
	return args.Bool(0), args.Error(1)
}

func (o *userRepositoryMock) Block(id string) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *userRepositoryMock) Unblock(id string) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *userRepositoryMock) PromoteToAdmin(id string) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *userRepositoryMock) DegradeToUser(id string) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *userRepositoryMock) RemoveHardly(id string) error {
	var args = o.Called(id)
	return args.Error(0)
}

func (o *userRepositoryMock) RemoveSoftly(id string) error {
	var args = o.Called(id)
	return args.Error(0)
}

func TestUserService_Save(t *testing.T) {
	defer beQuiet()()
	const (
		routine         = "Save"
		correctPassword = "Xxxxxx*0"
	)
	var (
		res                 uuid.UUID
		err                 error
		inserted            = uuid.New()
		correctUserCreation = &transfer.UserCreation{Password: correctPassword}
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, correctUserCreation).Return(inserted.String(), nil)
		res, err = NewUserService(r).Save(correctUserCreation)
		assert.Equal(t, inserted, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"creation\" cannot be nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).Save(nil)
		assert.Equal(t, uuid.Nil, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("Save", "creation").Error())
	})

	t.Run("must trim all string fields", func(t *testing.T) {
		var creation = &transfer.UserCreation{
			FirstName:  "  \a\b\f\n\r\t\vFirst Name\a\b\f\n\r\t\v  ",
			MiddleName: "  \a\b\f\n\r\t\vMiddle Name\a\b\f\n\r\t\v  ",
			LastName:   "  \a\b\f\n\r\t\vLast Name\a\b\f\n\r\t\v  ",
			Surname:    "  \a\b\f\n\r\t\vSurname\a\b\f\n\r\t\v  ",
			Email:      "  \a\b\f\n\r\t\vfoo@bar.com\a\b\f\n\r\t\v  ",
			Password:   correctPassword,
		}
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything).Return(inserted.String(), nil)
		res, err = NewUserService(r).Save(creation)
		assert.Equal(t, inserted, res)
		assert.Equal(t, "First Name", creation.FirstName)
		assert.Equal(t, "Middle Name", creation.MiddleName)
		assert.Equal(t, "Last Name", creation.LastName)
		assert.Equal(t, "Surname", creation.Surname)
		assert.Equal(t, "foo@bar.com", creation.Email)
		assert.NoError(t, err)
	})

	t.Run("must trim and bcrypt password", func(t *testing.T) {
		var creation = &transfer.UserCreation{Password: "  \a\b\f\r\t\v" + correctPassword + "\a\b\f\r\t\v  "}
		var trimmedPassword = correctPassword
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything).Return(inserted.String(), nil)
		res, err = NewUserService(r).Save(creation)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(creation.Password), []byte(trimmedPassword)))
		assert.NoError(t, err)
	})

	t.Run("satisfies...", func(t *testing.T) {
		var max = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxX"
		var creation = &transfer.UserCreation{Password: correctPassword}

		t.Run("50 < creation.FirstName", func(t *testing.T) {
			creation.FirstName = max
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Save(creation)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("FirstName", "user", 50).Error())
			assert.Equal(t, uuid.Nil, res)
			creation.FirstName = ""
		})

		t.Run("50 < creation.MiddleName", func(t *testing.T) {
			creation.MiddleName = max
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Save(creation)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("MiddleName", "user", 50).Error())
			assert.Equal(t, uuid.Nil, res)
			creation.MiddleName = ""
		})

		t.Run("50 < creation.LastName", func(t *testing.T) {
			creation.LastName = max
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Save(creation)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("LastName", "user", 50).Error())
			assert.Equal(t, uuid.Nil, res)
			creation.LastName = ""
		})

		t.Run("50 < creation.Surname", func(t *testing.T) {
			creation.Surname = max
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Save(creation)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Surname", "user", 50).Error())
			assert.Equal(t, uuid.Nil, res)
			creation.Surname = ""
		})

		max = max + max + max + max + max // >250
		t.Run("240 < creation.Email", func(t *testing.T) {
			creation.Email = max
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Save(creation)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Email", "user", 240).Error())
			assert.Equal(t, uuid.Nil, res)
			creation.Email = ""
		})

		t.Run("72 < creation.Password", func(t *testing.T) {
			creation.Password = max + "0*"
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Save(creation)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Password", "user", 72).Error())
			assert.Equal(t, uuid.Nil, res)
		})
	})

	t.Run("password must...", func(t *testing.T) {
		var (
			criterion string
			creation  = new(transfer.UserCreation)
		)

		creation.Password = "Xxxxx0*"
		criterion = "be at least 8 characters long"
		t.Run(criterion, func(t *testing.T) {
			var r = newUserRepositoryMock()
			r.On(routine, creation).Return(inserted.String(), nil)
			res, err = NewUserService(r).Save(creation)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, criterion)
		})

		creation.Password = "Xxxxxxx*"
		criterion = "contain at least one digit"
		t.Run(criterion, func(t *testing.T) {
			var r = newUserRepositoryMock()
			r.On(routine, creation).Return(inserted.String(), nil)
			res, err = NewUserService(r).Save(creation)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, criterion)
		})

		creation.Password = "xxxxxx0*"
		criterion = "contain at least one uppercase letter"
		t.Run(criterion, func(t *testing.T) {
			var r = newUserRepositoryMock()
			r.On(routine, creation).Return(inserted.String(), nil)
			res, err = NewUserService(r).Save(creation)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, criterion)
		})

		creation.Password = "XXXXXX0*"
		criterion = "contain at least one lowercase letter"
		t.Run(criterion, func(t *testing.T) {
			var r = newUserRepositoryMock()
			r.On(routine, creation).Return(inserted.String(), nil)
			res, err = NewUserService(r).Save(creation)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, criterion)
		})

		creation.Password = "XXXXXXx0"
		criterion = "contain at least one special character"
		t.Run(criterion, func(t *testing.T) {
			var r = newUserRepositoryMock()
			r.On(routine, creation).Return(inserted.String(), nil)
			res, err = NewUserService(r).Save(creation)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, criterion)
		})

		creation.Password = correctPassword
		creation.Email = correctPassword + "@example.com"
		t.Run("differ from email", func(t *testing.T) {
			var r = newUserRepositoryMock()
			r.On(routine, creation).Return(inserted.String(), nil)
			res, err = NewUserService(r).Save(creation)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, "seems to be similar to email")
		})
	})

	t.Run("did not parse inserted UUID", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything).Return("x", nil)
		res, err = NewUserService(r).Save(correctUserCreation)
		assert.ErrorContains(t, err, "invalid UUID length: 1")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything).Return("", unexpected)
		res, err = NewUserService(r).Save(correctUserCreation)
		assert.ErrorIs(t, err, unexpected)
		assert.Equal(t, uuid.Nil, res)
	})
}

func TestUserService_FetchByID(t *testing.T) {
	defer beQuiet()()
	const routine = "FetchShallowUserByID"
	var (
		res    *transfer.User
		err    error
		userID = uuid.New()
		user   = &transfer.User{ID: userID}
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(user, nil)
		res, err = NewUserService(r).FetchByID(userID)
		assert.Equal(t, user, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"id\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).FetchByID(uuid.Nil)
		assert.Nil(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("FetchByID", "id").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything).Return(nil, unexpected)
		res, err = NewUserService(r).FetchByID(userID)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}
