package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"strings"
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

func (o *listRepositoryMock) FetchScatteredLists(ownerID string, page, rpp int64, needle, sortBy string) ([]*model.List, error) {
	args := o.Called(ownerID, page, rpp, needle, sortBy)
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
	defer beQuiet()()
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

func TestListService_FetchListByID(t *testing.T) {
	defer beQuiet()()
	var (
		m                        *listRepositoryMock
		s                        *ListService
		res                      *model.List
		err                      error
		ownerID, groupID, listID = uuid.New(), uuid.New(), uuid.New()
		actual                   = &model.List{
			ID:          listID,
			OwnerID:     ownerID,
			GroupID:     groupID,
			Name:        "the list name (1)",
			Description: "description",
		}
	)

	t.Run("success", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("FetchListByID", mock.Anything, mock.Anything, mock.Anything).
			Return(actual, nil)
		s = NewListService(m)
		res, err = s.FindListByID(ownerID, groupID, listID)
		assert.NoError(t, err)
		assert.Equal(t, actual, res)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchListByID")
		s = NewListService(m)
		res, err = s.FindListByID(uuid.Nil, groupID, listID)
		assert.Nil(t, res)
		assert.ErrorContains(t, err, "parameter \"ownerID\" on function \"FindListByID\" cannot be uuid.Nil or nil")
	})

	t.Run("parameter groupID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchListByID")
		s = NewListService(m)
		res, err = s.FindListByID(ownerID, uuid.Nil, listID)
		assert.Nil(t, res)
		assert.ErrorContains(t, err, "parameter \"groupID\" on function \"FindListByID\" cannot be uuid.Nil or nil")
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchListByID")
		s = NewListService(m)
		res, err = s.FindListByID(ownerID, groupID, uuid.Nil)
		assert.Nil(t, res)
		assert.ErrorContains(t, err, "parameter \"listID\" on function \"FindListByID\" cannot be uuid.Nil or nil")
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("FetchListByID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		s = NewListService(m)
		res, err = s.FindListByID(ownerID, groupID, listID)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestListService_GetTodayListID(t *testing.T) {
	defer beQuiet()()
	var (
		m               *listRepositoryMock
		s               *ListService
		res             uuid.UUID
		err             error
		ownerID, listID = uuid.New(), uuid.New()
	)

	t.Run("success", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("GetTodayListID", mock.Anything).
			Return(listID.String(), nil)
		s = NewListService(m)
		res, err = s.GetTodayListID(ownerID)
		assert.Equal(t, listID, res)
		assert.NoError(t, err)
	})

	t.Run("got UUID parsing error", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("GetTodayListID", mock.Anything).
			Return("x", nil)
		s = NewListService(m)
		res, err = s.GetTodayListID(ownerID)
		assert.ErrorContains(t, err, "invalid UUID length: 1")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("did parse UUID", func(t *testing.T) {
		var id = uuid.MustParse("4fedb41f-5e44-4e63-9266-4b094bd7ba2d")
		m = new(listRepositoryMock)
		m.On("GetTodayListID", mock.Anything).
			Return(id.String(), nil)
		s = NewListService(m)
		res, err = s.GetTodayListID(ownerID)
		assert.Equal(t, id, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "GetTodayListID")
		s = NewListService(m)
		res, err = s.GetTodayListID(uuid.Nil)
		assert.Equal(t, uuid.Nil, res)
		assert.ErrorContains(t, err, "parameter \"ownerID\" on function \"GetTodayListID\" cannot be uuid.Nil or nil")
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("GetTodayListID", mock.Anything).
			Return("", unexpected)
		s = NewListService(m)
		res, err = s.GetTodayListID(ownerID)
		assert.ErrorIs(t, err, unexpected)
		assert.Equal(t, uuid.Nil, res)
	})
}

func TestListService_GetTomorrowListID(t *testing.T) {
	defer beQuiet()()
	var (
		m               *listRepositoryMock
		s               *ListService
		res             uuid.UUID
		err             error
		ownerID, listID = uuid.New(), uuid.New()
	)

	t.Run("success", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("GetTomorrowListID", mock.Anything).
			Return(listID.String(), nil)
		s = NewListService(m)
		res, err = s.GetTomorrowListID(ownerID)
		assert.Equal(t, listID, res)
		assert.NoError(t, err)
	})

	t.Run("got UUID parsing error", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("GetTomorrowListID", mock.Anything).
			Return("x", nil)
		s = NewListService(m)
		res, err = s.GetTomorrowListID(ownerID)
		assert.ErrorContains(t, err, "invalid UUID length: 1")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("did parse UUID", func(t *testing.T) {
		var id = uuid.MustParse("4fedb41f-5e44-4e63-9266-4b094bd7ba2d")
		m = new(listRepositoryMock)
		m.On("GetTomorrowListID", mock.Anything).
			Return(id.String(), nil)
		s = NewListService(m)
		res, err = s.GetTomorrowListID(ownerID)
		assert.Equal(t, id, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "GetTomorrowListID")
		s = NewListService(m)
		res, err = s.GetTomorrowListID(uuid.Nil)
		assert.Equal(t, uuid.Nil, res)
		assert.ErrorContains(t, err, "parameter \"ownerID\" on function \"GetTomorrowListID\" cannot be uuid.Nil or nil")
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("GetTomorrowListID", mock.Anything).
			Return("", unexpected)
		s = NewListService(m)
		res, err = s.GetTomorrowListID(ownerID)
		assert.ErrorIs(t, err, unexpected)
		assert.Equal(t, uuid.Nil, res)
	})
}

func TestListService_FindLists(t *testing.T) {
	defer beQuiet()()
	var (
		m          *listRepositoryMock
		s          *ListService
		res        *types.Result[model.List]
		err        error
		ownerID    = uuid.New()
		pagination = &types.Pagination{Page: 1, RPP: 10}
	)

	t.Run("success", func(t *testing.T) {
		var lists = make([]*model.List, 0)
		var current = &types.Result[model.List]{
			Page:      1,
			RPP:       10,
			Retrieved: int64(len(lists)),
			Payload:   lists,
		}
		m = new(listRepositoryMock)
		m.On("FetchLists",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindLists(ownerID, pagination, "", "")
		assert.Equal(t, current, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchLists")
		s = NewListService(m)
		res, err = s.FindLists(uuid.Nil, pagination, "", "")
		assert.ErrorContains(t, err, "parameter \"ownerID\" on function \"FindLists\" cannot be uuid.Nil or ni")
		assert.Nil(t, res)
	})

	t.Run("parameter pagination cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchLists")
		s = NewListService(m)
		res, err = s.FindLists(ownerID, nil, "", "")
		assert.ErrorContains(t, err, "parameter \"pagination\" on function \"FindLists\" cannot be uuid.Nil or ni")
		assert.Nil(t, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("FetchLists",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		s = NewListService(m)
		res, err = s.FindLists(ownerID, pagination, "", "")
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestListService_FindGroupedLists(t *testing.T) {
	defer beQuiet()()
	var (
		m                *listRepositoryMock
		s                *ListService
		res              *types.Result[model.List]
		err              error
		ownerID, groupID = uuid.New(), uuid.New()
		pagination       = &types.Pagination{Page: 1, RPP: 10}
	)

	t.Run("success", func(t *testing.T) {
		var lists = make([]*model.List, 0)
		var current = &types.Result[model.List]{
			Page:      1,
			RPP:       10,
			Retrieved: int64(len(lists)),
			Payload:   lists,
		}
		m = new(listRepositoryMock)
		m.On("FetchGroupedLists",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, pagination, "", "")
		assert.Equal(t, current, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchGroupedLists")
		s = NewListService(m)
		res, err = s.FindGroupedLists(uuid.Nil, groupID, pagination, "", "")
		assert.ErrorContains(t, err, "parameter \"ownerID\" on function \"FindGroupedLists\" cannot be uuid.Nil or ni")
		assert.Nil(t, res)
	})

	t.Run("parameter groupID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchGroupedLists")
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, uuid.Nil, pagination, "", "")
		assert.ErrorContains(t, err, "parameter \"groupID\" on function \"FindGroupedLists\" cannot be uuid.Nil or ni")
		assert.Nil(t, res)
	})

	t.Run("parameter pagination cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchGroupedLists")
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, nil, "", "")
		assert.ErrorContains(t, err, "parameter \"pagination\" on function \"FindGroupedLists\" cannot be uuid.Nil or ni")
		assert.Nil(t, res)
	})

	t.Run("parameter needle must be trimmed", func(t *testing.T) {
		var (
			lists  = make([]*model.List, 0)
			needle = "\n		needle 		\n"
		)
		m = new(listRepositoryMock)
		m.On("FetchGroupedLists",
			ownerID.String(), groupID.String(), pagination.Page, pagination.RPP,
			strings.Trim(needle, " \n\t"), "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, pagination, needle, "")
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("if pagination.Page is non-positive, then set it to 1", func(t *testing.T) {
		const expectedPageNumber int64 = 1
		var (
			lists = make([]*model.List, 0)
			pag   = &types.Pagination{Page: 0, RPP: 1}
		)

		/* when page=0 */

		m = new(listRepositoryMock)
		m.On("FetchGroupedLists",
			ownerID.String(), groupID.String(), expectedPageNumber, pag.RPP, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)

		/* when page<0 */

		pag.Page = -1
		m = new(listRepositoryMock)
		m.On("FetchGroupedLists",
			ownerID.String(), groupID.String(), expectedPageNumber, pag.RPP, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("if pagination.RPP is non-positive, then set it to 10", func(t *testing.T) {
		const expectedRPPNumber int64 = 10
		var (
			lists = make([]*model.List, 0)
			pag   = &types.Pagination{Page: 1, RPP: 0}
		)

		/* when RPP=0 */

		m = new(listRepositoryMock)
		m.On("FetchGroupedLists",
			ownerID.String(), groupID.String(), pag.Page, expectedRPPNumber, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)

		/* when RPP<0 */

		pag.RPP = -1
		m = new(listRepositoryMock)
		m.On("FetchGroupedLists",
			ownerID.String(), groupID.String(), pag.Page, expectedRPPNumber, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter sortBy must be trimmed", func(t *testing.T) {
		var (
			lists  = make([]*model.List, 0)
			sortBy = "\n		+first_name 		\n"
		)
		m = new(listRepositoryMock)
		m.On("FetchGroupedLists",
			ownerID.String(), groupID.String(), pagination.Page, pagination.RPP, "",
			strings.Trim(sortBy, " \n\t")).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, pagination, "", sortBy)
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("FetchGroupedLists",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, pagination, "", "")
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}
