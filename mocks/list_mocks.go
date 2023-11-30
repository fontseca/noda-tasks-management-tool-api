package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
)

type ListService struct {
	mock.Mock
}

func NewListServiceMock() *ListService {
	return new(ListService)
}

func (o *ListService) Save(ownerID, groupID uuid.UUID, next *transfer.ListCreation) (insertedID uuid.UUID, err error) {
	var args = o.Called(ownerID, groupID, next)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *ListService) FetchByID(ownerID, groupID, listID uuid.UUID) (list *model.List, err error) {
	var args = o.Called(ownerID, groupID, listID)
	var arg1 = args.Get(0)
	if nil != arg1 {
		list = arg1.(*model.List)
	}
	return list, args.Error(1)
}

func (o *ListService) GetTodayListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	var args = o.Called(ownerID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *ListService) GetTomorrowListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	var args = o.Called(ownerID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *ListService) Fetch(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (lists *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		lists = arg1.(*types.Result[model.List])
	}
	return lists, args.Error(1)
}

func (o *ListService) FetchGrouped(ownerID, groupID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	var args = o.Called(ownerID, groupID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		result = arg1.(*types.Result[model.List])
	}
	return result, args.Error(1)
}

func (o *ListService) FetchScattered(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		result = arg1.(*types.Result[model.List])
	}
	return result, args.Error(1)
}

func (o *ListService) Remove(ownerID, groupID, listID uuid.UUID) error {
	var args = o.Called(ownerID, groupID, listID)
	return args.Error(0)
}

func (o *ListService) Duplicate(ownerID, listID uuid.UUID) (replicaID uuid.UUID, err error) {
	var args = o.Called(ownerID, listID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *ListService) Scatter(ownerID, listID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID)
	return args.Bool(0), args.Error(1)
}

func (o *ListService) Move(ownerID, listID, targetGroupID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, targetGroupID)
	return args.Bool(0), args.Error(1)
}

func (o *ListService) Update(ownerID, groupID, listID uuid.UUID, up *transfer.ListUpdate) (ok bool, err error) {
	var args = o.Called(ownerID, groupID, listID, up)
	return args.Bool(0), args.Error(1)
}

type ListRepository struct {
	mock.Mock
}

func NewListRepositoryMock() *ListRepository {
	return new(ListRepository)
}

func (o *ListRepository) Save(ownerID, groupID string, next *transfer.ListCreation) (string, error) {
	args := o.Called(ownerID, groupID, next)
	return args.String(0), args.Error(1)
}

func (o *ListRepository) FetchByID(ownerID, groupID, listID string) (*model.List, error) {
	args := o.Called(ownerID, groupID, listID)
	arg1 := args.Get(0)
	var list *model.List
	if nil != arg1 {
		list = arg1.(*model.List)
	}
	return list, args.Error(1)
}

func (o *ListRepository) GetTodayListID(ownerID string) (string, error) {
	args := o.Called(ownerID)
	return args.String(0), args.Error(1)
}

func (o *ListRepository) GetTomorrowListID(ownerID string) (string, error) {
	args := o.Called(ownerID)
	return args.String(0), args.Error(1)
}

func (o *ListRepository) Fetch(ownerID string, page, rpp int64, needle, sortExpr string) ([]*model.List, error) {
	args := o.Called(ownerID, page, rpp, needle, sortExpr)
	arg1 := args.Get(0)
	var lists []*model.List
	if nil != arg1 {
		lists = arg1.([]*model.List)
	}
	return lists, args.Error(1)
}

func (o *ListRepository) FetchGrouped(ownerID, groupID string, page, rpp int64, needle, sortExpr string) ([]*model.List, error) {
	args := o.Called(ownerID, groupID, page, rpp, needle, sortExpr)
	arg1 := args.Get(0)
	var lists []*model.List
	if nil != arg1 {
		lists = arg1.([]*model.List)
	}
	return lists, args.Error(1)
}

func (o *ListRepository) FetchScattered(ownerID string, page, rpp int64, needle, sortExpr string) ([]*model.List, error) {
	args := o.Called(ownerID, page, rpp, needle, sortExpr)
	arg1 := args.Get(0)
	var lists []*model.List
	if nil != arg1 {
		lists = arg1.([]*model.List)
	}
	return lists, args.Error(1)
}

func (o *ListRepository) Remove(ownerID, groupID, listID string) (bool, error) {
	args := o.Called(ownerID, groupID, listID)
	return args.Bool(0), args.Error(1)
}

func (o *ListRepository) Duplicate(ownerID, listID string) (string, error) {
	args := o.Called(ownerID, listID)
	return args.String(0), args.Error(1)
}

func (o *ListRepository) Scatter(ownerID, listID string) (bool, error) {
	args := o.Called(ownerID, listID)
	return args.Bool(0), args.Error(1)
}

func (o *ListRepository) Move(ownerID, listID, targetGroupID string) (bool, error) {
	args := o.Called(ownerID, listID, targetGroupID)
	return args.Bool(0), args.Error(1)
}

func (o *ListRepository) Update(ownerID, groupID, listID string, up *transfer.ListUpdate) (bool, error) {
	args := o.Called(ownerID, groupID, listID, up)
	return args.Bool(0), args.Error(1)
}
