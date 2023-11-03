package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"noda/api/data/model"
	"noda/api/data/transfer"
	"testing"
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
	var list *model.List
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
	var lists []*model.List
	if nil != arg1 {
		lists = arg1.([]*model.List)
	}
	return lists, args.Error(1)
}

func (o *listRepositoryMock) FetchGroupedLists(ownerID, groupID string, page, rpp int64, needle, sortBy string) ([]*model.List, error) {
	args := o.Called(ownerID, groupID, page, rpp, needle, sortBy)
	arg1 := args.Get(0)
	var lists []*model.List
	if nil != arg1 {
		lists = arg1.([]*model.List)
	}
	return lists, args.Error(1)
}

func (o *listRepositoryMock) FetchScatteredLists(ownerID, groupID string, page, rpp int64, needle, sortBy string) ([]*model.List, error) {
	args := o.Called(ownerID, groupID, page, rpp, needle, sortBy)
	arg1 := args.Get(0)
	var lists []*model.List
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

func TestListService_SaveList(t *testing.T) {
	var (
		m                *listRepositoryMock
		s                *ListService
		res              uuid.UUID
		err              error
		ownerID, groupID = uuid.New(), uuid.New()
		next             = &transfer.ListCreation{
			Name:        "\t   list name\n   ",
			Description: "\n  description  \n",
		}
	)

	t.Run("success", func(t *testing.T) {
		insertedID := uuid.New()
		m = new(listRepositoryMock)
		m.On("InsertList", mock.Anything, mock.Anything, mock.Anything).
			Return(insertedID.String(), nil)
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, next)
		assert.Equal(t, insertedID, res)
		assert.NoError(t, err)
	})

	t.Run("got UUID parsing error", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("InsertList", mock.Anything, mock.Anything, mock.Anything).
			Return("x", nil)
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, next)
		assert.ErrorContains(t, err, "invalid UUID length: 1")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("did parse UUID", func(t *testing.T) {
		parsed := uuid.MustParse("4fedb41f-5e44-4e63-9266-4b094bd7ba2d")
		m = new(listRepositoryMock)
		m.On("InsertList", mock.Anything, mock.Anything, mock.Anything).
			Return(parsed.String(), nil)
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, next)
		assert.Equal(t, parsed, res)
		assert.NoError(t, err)
	})

	t.Run("name cannot be empty", func(t *testing.T) {
		var previousName = next.Name
		next.Name = "  		  \n"
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "InsertList")
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, next)
		next.Name = previousName
		assert.ErrorContains(t, err, "name cannot be an empty string")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "InsertList")
		s = NewListService(m)
		res, err = s.SaveList(uuid.Nil, groupID, next)
		assert.ErrorContains(t, err, "parameter \"ownerID\" on function \"SaveList\" cannot be uuid.Nil or nil")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("parameter groupID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "InsertList")
		s = NewListService(m)
		res, err = s.SaveList(ownerID, uuid.Nil, next)
		assert.ErrorContains(t, err, "parameter \"groupID\" on function \"SaveList\" cannot be uuid.Nil or nil")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("parameter next cannot be nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "InsertList")
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, nil)
		assert.ErrorContains(t, err, "parameter \"next\" on function \"SaveList\" cannot be uuid.Nil or nil")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("name too long: max length is 50 characters", func(t *testing.T) {
		var previousName = next.Name
		next.Name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxX"
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "InsertList")
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, next)
		next.Name = previousName
		assert.ErrorContains(t, err, "name too long")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("next.Name and next.Description must be trimmed", func(t *testing.T) {
		var previousName, previousDesc = next.Name, next.Description
		var insertedID = uuid.New()
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "InsertList")
		s = NewListService(m)
		m.On("InsertList", mock.Anything, mock.Anything, mock.Anything).
			Return(insertedID.String(), nil)
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, next)
		assert.Equal(t, "list name", next.Name)
		assert.Equal(t, "description", next.Description)
		next.Name, next.Description = previousName, previousDesc
		assert.Equal(t, insertedID, res)
		assert.NoError(t, err)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("InsertList", mock.Anything, mock.Anything, mock.Anything).
			Return("", unexpected)
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, next)
		assert.ErrorIs(t, err, unexpected)
		assert.Equal(t, uuid.Nil, res)
	})
}
