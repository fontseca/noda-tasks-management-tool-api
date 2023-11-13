package handler

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
)

type mockListService struct {
	mock.Mock
}

func (o *mockListService) SaveList(ownerID, groupID uuid.UUID, next *transfer.ListCreation) (insertedID uuid.UUID, err error) {
	var args = o.Called(ownerID, groupID, next)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *mockListService) FindListByID(ownerID, groupID, listID uuid.UUID) (list *model.List, err error) {
	var args = o.Called(ownerID, groupID, listID)
	var arg1 = args.Get(0)
	if nil != arg1 {
		list = arg1.(*model.List)
	}
	return list, args.Error(1)
}

func (o *mockListService) GetTodayListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	var args = o.Called(ownerID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *mockListService) GetTomorrowListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	var args = o.Called(ownerID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *mockListService) FindLists(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (lists *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		lists = arg1.(*types.Result[model.List])
	}
	return lists, args.Error(1)
}

func (o *mockListService) FindGroupedLists(ownerID, groupID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		result = arg1.(*types.Result[model.List])
	}
	return result, args.Error(1)
}

func (o *mockListService) FindScatteredLists(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		result = arg1.(*types.Result[model.List])
	}
	return result, args.Error(1)
}

func (o *mockListService) DeleteList(ownerID, groupID, listID uuid.UUID) error {
	var args = o.Called(ownerID, groupID, listID)
	return args.Error(1)
}

func (o *mockListService) DuplicateList(ownerID, listID uuid.UUID) (replicaID uuid.UUID, err error) {
	var args = o.Called(ownerID, listID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *mockListService) ConvertToScatteredList(ownerID, listID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID)
	return args.Bool(0), args.Error(1)
}

func (o *mockListService) MoveList(ownerID, listID, targetGroupID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, targetGroupID)
	return args.Bool(0), args.Error(1)
}

func (o *mockListService) UpdateList(ownerID, groupID, listID uuid.UUID, up *transfer.ListUpdate) (ok bool, err error) {
	var args = o.Called(ownerID, groupID, listID)
	return args.Bool(0), args.Error(1)
}
