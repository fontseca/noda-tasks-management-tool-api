package mocks

import (
	"github.com/stretchr/testify/mock"
	"noda/data/model"
	"noda/data/transfer"
)

type GroupRepository struct {
	mock.Mock
}

func NewGroupRepositoryMock() *GroupRepository {
	return new(GroupRepository)
}

func (o *GroupRepository) Save(ownerID string, next *transfer.GroupCreation) (string, error) {
	args := o.Called(ownerID, next)
	return args.String(0), args.Error(1)
}

func (o *GroupRepository) FetchByID(ownerID, groupID string) (*model.Group, error) {
	args := o.Called(ownerID, groupID)
	var group *model.Group
	arg1 := args.Get(0)
	if nil != arg1 {
		group = arg1.(*model.Group)
	}
	return group, args.Error(1)
}

func (o *GroupRepository) Fetch(ownerID string, page, rpp int64, needle, sortBy string) ([]*model.Group, error) {
	args := o.Called(ownerID, page, rpp, needle, sortBy)
	var groups []*model.Group
	arg1 := args.Get(0)
	if nil != arg1 {
		groups = arg1.([]*model.Group)
	}
	return groups, args.Error(1)
}

func (o *GroupRepository) Update(ownerID, groupID string, up *transfer.GroupUpdate) (ok bool, err error) {
	args := o.Called(ownerID, groupID, up)
	return args.Bool(0), args.Error(1)
}

func (o *GroupRepository) Remove(ownerID, groupID string) (ok bool, err error) {
	args := o.Called(ownerID, groupID)
	return args.Bool(0), args.Error(1)
}
