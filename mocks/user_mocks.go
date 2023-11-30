package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
)

type UserRepository struct {
	mock.Mock
}

func NewUserRepositoryMock() *UserRepository {
	return new(UserRepository)
}

func (o *UserRepository) Save(creation *transfer.UserCreation) (insertedID string, err error) {
	var args = o.Called(creation)
	return args.String(0), args.Error(1)
}

func (o *UserRepository) FetchByID(id string) (user *model.User, err error) {
	var args = o.Called(id)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*model.User)
	}
	return user, args.Error(1)
}

func (o *UserRepository) FetchShallowUserByID(id string) (user *transfer.User, err error) {
	var args = o.Called(id)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*transfer.User)
	}
	return user, args.Error(1)
}

func (o *UserRepository) FetchByEmail(email string) (user *model.User, err error) {
	var args = o.Called(email)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*model.User)
	}
	return user, args.Error(1)
}

func (o *UserRepository) FetchShallowUserByEmail(email string) (user *transfer.User, err error) {
	var args = o.Called(email)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*transfer.User)
	}
	return user, args.Error(1)
}

func (o *UserRepository) Fetch(page, rpp int64, needle, sortExpr string) (users []*transfer.User, err error) {
	var args = o.Called(page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		users = arg0.([]*transfer.User)
	}
	return users, args.Error(1)
}

func (o *UserRepository) FetchBlocked(page, rpp int64, needle, sortExpr string) (users []*transfer.User, err error) {
	var args = o.Called(page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		users = arg0.([]*transfer.User)
	}
	return users, args.Error(1)
}

func (o *UserRepository) FetchSettings(userID string, page, rpp int64, needle, sortExpr string) (settings []*transfer.UserSetting, err error) {
	var args = o.Called(userID, page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		settings = arg0.([]*transfer.UserSetting)
	}
	return settings, args.Error(1)
}

func (o *UserRepository) FetchOneSetting(userID string, settingKey string) (setting *transfer.UserSetting, err error) {
	var args = o.Called(userID, settingKey)
	var arg0 = args.Get(0)
	if nil != arg0 {
		setting = arg0.(*transfer.UserSetting)
	}
	return setting, args.Error(1)
}

func (o *UserRepository) Search(page, rpp int64, needle, sortExpr string) (users []*transfer.User, err error) {
	var args = o.Called(page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		users = arg0.([]*transfer.User)
	}
	return users, args.Error(1)
}

func (o *UserRepository) Update(id string, update *transfer.UserUpdate) (ok bool, err error) {
	var args = o.Called(id, update)
	return args.Bool(0), args.Error(1)
}

func (o *UserRepository) UpdateUserSetting(userID, settingKey, newValue string) (ok bool, err error) {
	var args = o.Called(userID, settingKey, newValue)
	return args.Bool(0), args.Error(1)
}

func (o *UserRepository) Block(id string) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *UserRepository) Unblock(id string) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *UserRepository) PromoteToAdmin(id string) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *UserRepository) DegradeToUser(id string) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *UserRepository) RemoveHardly(id string) error {
	var args = o.Called(id)
	return args.Error(0)
}

func (o *UserRepository) RemoveSoftly(id string) error {
	var args = o.Called(id)
	return args.Error(0)
}

type UserService struct {
	mock.Mock
}

func NewUserServiceMock() *UserService {
	return new(UserService)
}

func (o *UserService) Save(creation *transfer.UserCreation) (insertedID uuid.UUID, err error) {
	var args = o.Called(creation)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *UserService) FetchByID(id uuid.UUID) (user *transfer.User, err error) {
	var args = o.Called(id)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*transfer.User)
	}
	return user, args.Error(1)
}

func (o *UserService) FetchByEmail(email string) (user *transfer.User, err error) {
	var args = o.Called(email)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*transfer.User)
	}
	return user, args.Error(1)
}

func (o *UserService) FetchRawUserByEmail(email string) (user *model.User, err error) {
	var args = o.Called(email)
	var arg0 = args.Get(0)
	if nil != arg0 {
		user = arg0.(*model.User)
	}
	return user, args.Error(1)
}

func (o *UserService) Fetch(pagination *types.Pagination, needle, sortExpr string) (result *types.Result[transfer.User], err error) {
	var args = o.Called(pagination, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		result = arg0.(*types.Result[transfer.User])
	}
	return result, args.Error(1)
}

func (o *UserService) FetchBlocked(pagination *types.Pagination, needle, sortExpr string) (result *types.Result[transfer.User], err error) {
	var args = o.Called(pagination, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		result = arg0.(*types.Result[transfer.User])
	}
	return result, args.Error(1)
}

func (o *UserService) FetchSettings(userID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[transfer.UserSetting], err error) {
	var args = o.Called(userID, pagination, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		result = arg0.(*types.Result[transfer.UserSetting])
	}
	return result, args.Error(1)
}

func (o *UserService) FetchOneSetting(userID uuid.UUID, settingKey string) (setting *transfer.UserSetting, err error) {
	var args = o.Called(userID, settingKey)
	var arg0 = args.Get(0)
	if nil != arg0 {
		setting = arg0.(*transfer.UserSetting)
	}
	return setting, args.Error(1)
}

func (o *UserService) Search(pagination *types.Pagination, needle, sortExpr string) (users *types.Result[transfer.User], err error) {
	var args = o.Called(pagination, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		users = arg0.(*types.Result[transfer.User])
	}
	return users, args.Error(1)
}

func (o *UserService) Update(id uuid.UUID, update *transfer.UserUpdate) (ok bool, err error) {
	var args = o.Called(id, update)
	return args.Bool(0), args.Error(1)
}

func (o *UserService) UpdateUserSetting(userID uuid.UUID, settingKey string, update *transfer.UserSettingUpdate) (ok bool, err error) {
	var args = o.Called(userID, settingKey, update)
	return args.Bool(0), args.Error(1)
}

func (o *UserService) Block(id uuid.UUID) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *UserService) Unblock(id uuid.UUID) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *UserService) PromoteToAdmin(id uuid.UUID) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *UserService) DegradeToUser(id uuid.UUID) (ok bool, err error) {
	var args = o.Called(id)
	return args.Bool(0), args.Error(1)
}

func (o *UserService) RemoveHardly(id uuid.UUID) error {
	return o.Called(id).Error(0)
}

func (o *UserService) RemoveSoftly(id uuid.UUID) error {
	return o.Called(id).Error(1)
}
