package service

import (
	"github.com/stretchr/testify/mock"
	"noda/api/data/model"
	"noda/api/data/transfer"
)

type listRepositoryMock struct {
	mock.Mock
}

func (o *listRepositoryMock) InsertList(ownerID, groupID string, next *transfer.ListCreation) (string, error) {
	args := o.Called(ownerID, groupID, next)
	return args.String(0), args.Error(1)
}

func (o *listRepositoryMock) FetchListByID(ownerID, groupID, listID string) (*model.List, error) {
	args := o.Called(ownerID, groupID, listID)
	arg1 := args.Get(0)
	var list = new(model.List)
	if nil != arg1 {
		list = arg1.(*model.List)
	}
	return list, args.Error(1)
}

func (o *listRepositoryMock) GetTodayListID(ownerID string) (string, error) {
	args := o.Called(ownerID)
	return args.String(0), args.Error(1)
}

func (o *listRepositoryMock) GetTomorrowListID(ownerID string) (string, error) {
	args := o.Called(ownerID)
	return args.String(0), args.Error(1)
}

func (o *listRepositoryMock) FetchLists(ownerID string, page, rpp int64, needle, sortExpr string) ([]*model.List, error) {
	args := o.Called(ownerID, page, rpp, needle, sortExpr)
	arg1 := args.Get(0)
	var lists = make([]*model.List, 0)
	if nil != arg1 {
		lists = arg1.([]*model.List)
	}
	return lists, args.Error(1)
}

func (o *listRepositoryMock) FetchGroupedLists(ownerID, groupID string, page, rpp int64, needle, sortBy string) ([]*model.List, error) {
	args := o.Called(ownerID, groupID, page, rpp, needle, sortBy)
	arg1 := args.Get(0)
	var lists = make([]*model.List, 0)
	if nil != arg1 {
		lists = arg1.([]*model.List)
	}
	return lists, args.Error(1)
}

func (o *listRepositoryMock) FetchScatteredLists(ownerID, groupID string, page, rpp int64, needle, sortBy string) ([]*model.List, error) {
	args := o.Called(ownerID, groupID, page, rpp, needle, sortBy)
	arg1 := args.Get(0)
	var lists = make([]*model.List, 0)
	if nil != arg1 {
		lists = arg1.([]*model.List)
	}
	return lists, args.Error(1)
}

func (o *listRepositoryMock) DeleteList(ownerID, groupID, listID string) (bool, error) {
	args := o.Called(ownerID, groupID, listID)
	return args.Bool(0), args.Error(1)
}

func (o *listRepositoryMock) DuplicateList(ownerID, listID string) (string, error) {
	args := o.Called(ownerID, listID)
	return args.String(0), args.Error(1)
}

func (o *listRepositoryMock) ConvertToScatteredList(ownerID, listID string) (bool, error) {
	args := o.Called(ownerID, listID)
	return args.Bool(0), args.Error(1)
}

func (o *listRepositoryMock) MoveList(ownerID, listID, targetGroupID string) (bool, error) {
	args := o.Called(ownerID, listID, targetGroupID)
	return args.Bool(0), args.Error(1)
}

func (o *listRepositoryMock) UpdateList(ownerID, groupID, listID string, up *transfer.ListUpdate) (bool, error) {
	args := o.Called(ownerID, groupID, listID, up)
	return args.Bool(0), args.Error(1)
}
