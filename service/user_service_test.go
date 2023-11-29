package service

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
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
			FirstName:  blankset + "First Name" + blankset,
			MiddleName: blankset + "Middle Name" + blankset,
			LastName:   blankset + "Last Name" + blankset,
			Surname:    blankset + "Surname" + blankset,
			Email:      blankset + "foo@bar.com" + blankset,
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
		var creation = &transfer.UserCreation{Password: blankset + correctPassword + blankset}
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

func TestUserService_FetchByEmail(t *testing.T) {
	defer beQuiet()()
	const (
		routine = "FetchShallowUserByEmail"
		email   = "foo@bar.com"
	)
	var (
		res  *transfer.User
		err  error
		user = &transfer.User{ID: uuid.New(), Email: email}
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, email).Return(user, nil)
		res, err = NewUserService(r).FetchByEmail(email)
		assert.Equal(t, user, res)
		assert.NoError(t, err)
	})

	t.Run("must trim parameter \"email\"", func(t *testing.T) {
		var e = blankset + email + blankset
		var r = newUserRepositoryMock()
		r.On(routine, email).Return(user, nil)
		res, err = NewUserService(r).FetchByEmail(e)
		assert.Equal(t, user, res)
		assert.NoError(t, err)
	})

	t.Run("empty email? then noda.ErrUserNotFound", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).FetchByEmail(blankset)
		assert.ErrorContains(t, err, noda.ErrUserNotFound.Error())
		assert.Nil(t, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything).Return(nil, unexpected)
		res, err = NewUserService(r).FetchByEmail(email)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestUserService_FetchRawUserByEmail(t *testing.T) {
	defer beQuiet()()
	const (
		routine = "FetchByEmail"
		email   = "foo@bar.com"
	)
	var (
		res  *model.User
		err  error
		user = &model.User{ID: uuid.New(), Email: email}
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, email).Return(user, nil)
		res, err = NewUserService(r).FetchRawUserByEmail(email)
		assert.Equal(t, user, res)
		assert.NoError(t, err)
	})

	t.Run("must trim parameter \"email\"", func(t *testing.T) {
		var e = blankset + email + blankset
		var r = newUserRepositoryMock()
		r.On(routine, email).Return(user, nil)
		res, err = NewUserService(r).FetchRawUserByEmail(e)
		assert.Equal(t, user, res)
		assert.NoError(t, err)
	})

	t.Run("empty email? then noda.ErrUserNotFound", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).FetchRawUserByEmail(blankset)
		assert.ErrorContains(t, err, noda.ErrUserNotFound.Error())
		assert.Nil(t, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything).Return(nil, unexpected)
		res, err = NewUserService(r).FetchRawUserByEmail(email)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestUserService_Fetch(t *testing.T) {
	defer beQuiet()()
	const routine = "Fetch"
	var (
		res        *types.Result[transfer.User]
		err        error
		needle     = "user"
		sortExpr   = "+first_name"
		users      = make([]*transfer.User, 2)
		pagination = &types.Pagination{Page: 1, RPP: 10}
	)

	t.Run("success", func(t *testing.T) {
		var result = &types.Result[transfer.User]{
			Page:      pagination.Page,
			RPP:       pagination.RPP,
			Retrieved: int64(len(users)),
			Payload:   users,
		}
		var r = newUserRepositoryMock()
		r.On(routine, pagination.Page, pagination.RPP, needle, sortExpr).Return(users, nil)
		res, err = NewUserService(r).Fetch(pagination, needle, sortExpr)
		assert.Equal(t, result, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"pagination\" cannot be nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).Fetch(nil, needle, sortExpr)
		assert.ErrorContains(t, err, noda.NewNilParameterError("Fetch", "pagination").Error())
		assert.Nil(t, res)
	})

	t.Run("must default pagination fields", func(t *testing.T) {
		const expectedPage, expectedRPP int64 = 1, 10
		pagination.Page = -1
		pagination.RPP = 0
		var r = newUserRepositoryMock()
		r.On(routine, expectedPage, expectedRPP, mock.Anything, mock.Anything).Return(users, nil)
		_, _ = NewUserService(r).Fetch(pagination, needle, sortExpr)
	})

	t.Run("must trim \"needle\" parameter", func(t *testing.T) {
		var n = blankset + needle + blankset
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, needle, mock.Anything).Return(users, nil)
		_, _ = NewUserService(r).Fetch(pagination, n, sortExpr)
	})

	t.Run("must trim \"sortExpr\" parameter", func(t *testing.T) {
		var s = blankset + sortExpr + blankset
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, sortExpr).Return(users, nil)
		_, _ = NewUserService(r).Fetch(pagination, needle, s)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, unexpected)
		res, err = NewUserService(r).Fetch(pagination, needle, sortExpr)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestUserService_FetchBlocked(t *testing.T) {
	defer beQuiet()()
	const routine = "FetchBlocked"
	var (
		res        *types.Result[transfer.User]
		err        error
		needle     = "user"
		sortExpr   = "+first_name"
		users      = make([]*transfer.User, 2)
		pagination = &types.Pagination{Page: 1, RPP: 10}
	)

	t.Run("success", func(t *testing.T) {
		var result = &types.Result[transfer.User]{
			Page:      pagination.Page,
			RPP:       pagination.RPP,
			Retrieved: int64(len(users)),
			Payload:   users,
		}
		var r = newUserRepositoryMock()
		r.On(routine, pagination.Page, pagination.RPP, needle, sortExpr).Return(users, nil)
		res, err = NewUserService(r).FetchBlocked(pagination, needle, sortExpr)
		assert.Equal(t, result, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"pagination\" cannot be nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).FetchBlocked(nil, needle, sortExpr)
		assert.ErrorContains(t, err, noda.NewNilParameterError("FetchBlocked", "pagination").Error())
		assert.Nil(t, res)
	})

	t.Run("must default pagination fields", func(t *testing.T) {
		const expectedPage, expectedRPP int64 = 1, 10
		pagination.Page = -1
		pagination.RPP = 0
		var r = newUserRepositoryMock()
		r.On(routine, expectedPage, expectedRPP, mock.Anything, mock.Anything).Return(users, nil)
		_, _ = NewUserService(r).FetchBlocked(pagination, needle, sortExpr)
	})

	t.Run("must trim \"needle\" parameter", func(t *testing.T) {
		var n = blankset + needle + blankset
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, needle, mock.Anything).Return(users, nil)
		_, _ = NewUserService(r).FetchBlocked(pagination, n, sortExpr)
	})

	t.Run("must trim \"sortExpr\" parameter", func(t *testing.T) {
		var s = blankset + sortExpr + blankset
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, sortExpr).Return(users, nil)
		_, _ = NewUserService(r).FetchBlocked(pagination, needle, s)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, unexpected)
		res, err = NewUserService(r).FetchBlocked(pagination, needle, sortExpr)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestUserService_FetchSettings(t *testing.T) {
	defer beQuiet()()
	const routine = "FetchSettings"
	var (
		res        *types.Result[transfer.UserSetting]
		err        error
		needle     = "setting"
		sortExpr   = "+first_name"
		pagination = &types.Pagination{Page: 1, RPP: 10}
		userID     = uuid.New()
		settings   = make([]*transfer.UserSetting, 2)
	)

	t.Run("success", func(t *testing.T) {
		var s = []*transfer.UserSetting{
			{
				Key:   "setting 1",
				Value: []byte("true"),
			},
			{
				Key:   "setting 2",
				Value: []byte("{\"key_1\":true,\"key_2\":\"value_2\",\"key_3\":null}"),
			},
			{
				Key:   "setting 3",
				Value: []byte("3.14"),
			},
		}
		var r = newUserRepositoryMock()
		r.On(routine, userID.String(), pagination.Page, pagination.RPP, needle, sortExpr).Return(s, nil)
		res, err = NewUserService(r).FetchSettings(userID, pagination, needle, sortExpr)
		assert.NoError(t, err)
		assert.NotNil(t, res.Payload)
		assert.Equal(t, int64(1), res.Page)
		assert.Equal(t, int64(10), res.RPP)
		assert.Equal(t, int64(3), res.Retrieved)
		var value, expected any
		value, expected = res.Payload[0].Value, true
		assert.IsType(t, expected, value)
		assert.Equal(t, expected, value)
		value, expected = res.Payload[1].Value, map[string]any{"key_1": true, "key_2": "value_2", "key_3": nil}
		assert.IsType(t, expected, value)
		assert.Equal(t, expected, value)
		value, expected = res.Payload[2].Value, 3.14
		assert.IsType(t, expected, value)
		assert.Equal(t, expected, value)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).FetchSettings(uuid.Nil, pagination, needle, sortExpr)
		assert.Nil(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("FetchSettings", "userID").Error())
	})

	t.Run("parameter \"pagination\" cannot be nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).FetchSettings(userID, nil, needle, sortExpr)
		assert.ErrorContains(t, err, noda.NewNilParameterError("FetchSettings", "pagination").Error())
		assert.Nil(t, res)
	})

	t.Run("must default pagination fields", func(t *testing.T) {
		const expectedPage, expectedRPP int64 = 1, 10
		pagination.Page = -1
		pagination.RPP = 0
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, expectedPage, expectedRPP, mock.Anything, mock.Anything).Return(settings, nil)
		_, _ = NewUserService(r).FetchSettings(userID, pagination, needle, sortExpr)
	})

	t.Run("must trim \"needle\" parameter", func(t *testing.T) {
		var n = blankset + needle + blankset
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, needle, mock.Anything).Return(settings, nil)
		_, _ = NewUserService(r).FetchSettings(userID, pagination, n, sortExpr)
	})

	t.Run("must trim \"sortExpr\" parameter", func(t *testing.T) {
		var s = blankset + sortExpr + blankset
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, sortExpr).Return(settings, nil)
		_, _ = NewUserService(r).FetchSettings(userID, pagination, needle, s)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, unexpected)
		res, err = NewUserService(r).FetchSettings(userID, pagination, needle, sortExpr)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestUserService_FetchOneSetting(t *testing.T) {
	defer beQuiet()()
	const routine = "FetchOneSetting"
	var (
		userID     = uuid.New()
		settingKey = "key"
		res        *transfer.UserSetting
		err        error
		setting    = &transfer.UserSetting{
			Key:   "key",
			Value: []byte("\"yeah\""),
		}
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String(), settingKey).Return(setting, nil)
		res, err = NewUserService(r).FetchOneSetting(userID, settingKey)
		assert.NoError(t, err)
		assert.Equal(t, "yeah", setting.Value)
		assert.Equal(t, setting, res)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).FetchOneSetting(uuid.Nil, settingKey)
		assert.Nil(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("FetchOneSetting", "userID").Error())
	})

	t.Run("must trim \"settingKey\" parameter", func(t *testing.T) {
		setting.Value = []byte("\"yeah\"")
		var k = blankset + settingKey + blankset
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, settingKey).Return(setting, nil)
		_, _ = NewUserService(r).FetchOneSetting(userID, k)
	})

	t.Run("got a repository error", func(t *testing.T) {
		setting.Value = []byte("\"yeah\"")
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything).Return(nil, unexpected)
		res, err = NewUserService(r).FetchOneSetting(userID, settingKey)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestUserService_Update(t *testing.T) {
	defer beQuiet()()
	const routine = "Update"
	var (
		res         bool
		err         error
		placeholder = &transfer.UserUpdate{}
		userID      = uuid.New()
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String(), placeholder).Return(true, nil)
		res, err = NewUserService(r).Update(userID, placeholder)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).Update(uuid.Nil, placeholder)
		assert.False(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("Update", "userID").Error())
	})

	t.Run("parameter \"update\" cannot be nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).Update(userID, nil)
		assert.False(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("Update", "update").Error())
	})

	t.Run("must trim all string fields", func(t *testing.T) {
		var update = &transfer.UserUpdate{
			FirstName:  blankset + "First Name" + blankset,
			MiddleName: blankset + "Middle Name" + blankset,
			LastName:   blankset + "Last Name" + blankset,
			Surname:    blankset + "Surname" + blankset,
		}
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything).Return(true, nil)
		res, err = NewUserService(r).Update(userID, update)
		assert.True(t, res)
		assert.Equal(t, "First Name", update.FirstName)
		assert.Equal(t, "Middle Name", update.MiddleName)
		assert.Equal(t, "Last Name", update.LastName)
		assert.Equal(t, "Surname", update.Surname)
		assert.NoError(t, err)
	})

	t.Run("if not empty, satisfies...", func(t *testing.T) {
		var max = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxX"
		var update = &transfer.UserUpdate{}

		t.Run("50 < update.FirstName", func(t *testing.T) {
			update.FirstName = max
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Update(userID, update)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("FirstName", "user", 50).Error())
			assert.False(t, res)
			update.FirstName = ""
		})

		t.Run("50 < update.MiddleName", func(t *testing.T) {
			update.MiddleName = max
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Update(userID, update)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("MiddleName", "user", 50).Error())
			assert.False(t, res)
			update.MiddleName = ""
		})

		t.Run("50 < update.LastName", func(t *testing.T) {
			update.LastName = max
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Update(userID, update)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("LastName", "user", 50).Error())
			assert.False(t, res)
			update.LastName = ""
		})

		t.Run("50 < update.Surname", func(t *testing.T) {
			update.Surname = max
			var r = newUserRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewUserService(r).Update(userID, update)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Surname", "user", 50).Error())
			assert.False(t, res)
		})
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything).Return(false, unexpected)
		res, err = NewUserService(r).Update(userID, placeholder)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}

func TestUserService_UpdateUserSetting(t *testing.T) {
	defer beQuiet()()
	const routine = "UpdateUserSetting"
	var (
		res           bool
		err           error
		placeholder   = &transfer.UserSettingUpdate{Value: "{\"setting\":\"value\"}"}
		settingKey    = "key"
		userID        = uuid.New()
		buf, _        = json.Marshal(placeholder.Value)
		expectedValue = string(buf)
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String(), settingKey, expectedValue).Return(true, nil)
		res, err = NewUserService(r).UpdateUserSetting(userID, settingKey, placeholder)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).UpdateUserSetting(uuid.Nil, settingKey, placeholder)
		assert.False(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("UpdateUserSetting", "userID").Error())
	})

	t.Run("parameter \"update\" cannot be nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).UpdateUserSetting(userID, settingKey, nil)
		assert.False(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("UpdateUserSetting", "update").Error())
	})

	t.Run("empty \"settingKey\"? then do nothing", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).UpdateUserSetting(userID, blankset, placeholder)
		assert.False(t, res)
		assert.NoError(t, err)
	})

	t.Run("must trim \"settingKey\" parameter", func(t *testing.T) {
		var s = blankset + settingKey + blankset
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, settingKey, mock.Anything).Return(true, nil)
		_, _ = NewUserService(r).UpdateUserSetting(userID, s, placeholder)
	})

	t.Run("if value is string, trim it", func(t *testing.T) {
		var value = "setting value"
		var update = &transfer.UserSettingUpdate{
			Value: blankset + value + blankset,
		}
		var buf, _ = json.Marshal(value)
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, string(buf)).Return(true, nil)
		res, err = NewUserService(r).UpdateUserSetting(userID, settingKey, update)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything).Return(false, unexpected)
		res, err = NewUserService(r).UpdateUserSetting(userID, settingKey, placeholder)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}

func TestUserService_Block(t *testing.T) {
	const routine = "Block"
	var (
		userID = uuid.New()
		res    bool
		err    error
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(true, nil)
		res, err = NewUserService(r).Block(userID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).Block(uuid.Nil)
		assert.False(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("Block", "userID").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(false, unexpected)
		res, err = NewUserService(r).Block(userID)
		assert.False(t, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestUserService_Unblock(t *testing.T) {
	const routine = "Unblock"
	var (
		userID = uuid.New()
		res    bool
		err    error
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(true, nil)
		res, err = NewUserService(r).Unblock(userID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).Unblock(uuid.Nil)
		assert.False(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("Unblock", "userID").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(false, unexpected)
		res, err = NewUserService(r).Unblock(userID)
		assert.False(t, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestUserService_PromoteToAdmin(t *testing.T) {
	const routine = "PromoteToAdmin"
	var (
		userID = uuid.New()
		res    bool
		err    error
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(true, nil)
		res, err = NewUserService(r).PromoteToAdmin(userID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).PromoteToAdmin(uuid.Nil)
		assert.False(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("PromoteToAdmin", "userID").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(false, unexpected)
		res, err = NewUserService(r).PromoteToAdmin(userID)
		assert.False(t, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestUserService_DegradeToUser(t *testing.T) {
	const routine = "DegradeToUser"
	var (
		userID = uuid.New()
		res    bool
		err    error
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(true, nil)
		res, err = NewUserService(r).DegradeToUser(userID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		res, err = NewUserService(r).DegradeToUser(uuid.Nil)
		assert.False(t, res)
		assert.ErrorContains(t, err, noda.NewNilParameterError("DegradeToUser", "userID").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(false, unexpected)
		res, err = NewUserService(r).DegradeToUser(userID)
		assert.False(t, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestUserService_RemoveHardly(t *testing.T) {
	const routine = "RemoveHardly"
	var (
		userID = uuid.New()
		err    error
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(nil)
		err = NewUserService(r).RemoveHardly(userID)
		assert.NoError(t, err)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		err = NewUserService(r).RemoveHardly(uuid.Nil)
		assert.ErrorContains(t, err, noda.NewNilParameterError("RemoveHardly", "id").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(unexpected)
		err = NewUserService(r).RemoveHardly(userID)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestUserService_RemoveSoftly(t *testing.T) {
	const routine = "RemoveSoftly"
	var (
		userID = uuid.New()
		err    error
	)

	t.Run("success", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(nil)
		err = NewUserService(r).RemoveSoftly(userID)
		assert.NoError(t, err)
	})

	t.Run("parameter \"userID\" cannot be uuid.Nil", func(t *testing.T) {
		var r = newUserRepositoryMock()
		r.AssertNotCalled(t, routine)
		err = NewUserService(r).RemoveSoftly(uuid.Nil)
		assert.ErrorContains(t, err, noda.NewNilParameterError("RemoveSoftly", "id").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = newUserRepositoryMock()
		r.On(routine, userID.String()).Return(unexpected)
		err = NewUserService(r).RemoveSoftly(userID)
		assert.ErrorIs(t, err, unexpected)
	})
}
