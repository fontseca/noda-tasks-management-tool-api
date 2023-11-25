package service

import (
	"github.com/stretchr/testify/mock"
	"noda/data/model"
	"noda/data/transfer"
)

type userRepositoryMock struct {
	mock.Mock
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
