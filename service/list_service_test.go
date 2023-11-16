package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"noda"
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
		s                ListService
		res              uuid.UUID
		err              error
		ownerID, groupID = uuid.New(), uuid.New()
		next             = &transfer.ListCreation{
			Name:        "\t   list name\n   ",
			Description: "\n  description  \n",
		}
	)

	t.Run("success for grouped list", func(t *testing.T) {
		insertedID := uuid.New()
		m = new(listRepositoryMock)
		m.On("InsertList", ownerID.String(), groupID.String(), next).
			Return(insertedID.String(), nil)
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, next)
		assert.Equal(t, insertedID, res)
		assert.NoError(t, err)
	})

	t.Run("success for scattered list", func(t *testing.T) {
		insertedID := uuid.New()
		m = new(listRepositoryMock)
		m.On("InsertList", ownerID.String(), "", next).
			Return(insertedID.String(), nil)
		s = NewListService(m)
		res, err = s.SaveList(ownerID, uuid.Nil, next)
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
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("SaveList", "ownerID").Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("parameter next cannot be nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "InsertList")
		s = NewListService(m)
		res, err = s.SaveList(ownerID, groupID, nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("SaveList", "next").Error())
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
		s                        ListService
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

	t.Run("success for grouped list", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("FetchListByID", ownerID.String(), groupID.String(), listID.String()).
			Return(actual, nil)
		s = NewListService(m)
		res, err = s.FindListByID(ownerID, groupID, listID)
		assert.NoError(t, err)
		assert.Equal(t, actual, res)
	})

	t.Run("success for scattered list", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("FetchListByID", ownerID.String(), "", listID.String()).
			Return(actual, nil)
		s = NewListService(m)
		res, err = s.FindListByID(ownerID, uuid.Nil, listID)
		assert.NoError(t, err)
		assert.Equal(t, actual, res)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchListByID")
		s = NewListService(m)
		res, err = s.FindListByID(uuid.Nil, groupID, listID)
		assert.Nil(t, res)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FindListByID", "ownerID").Error())
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchListByID")
		s = NewListService(m)
		res, err = s.FindListByID(ownerID, groupID, uuid.Nil)
		assert.Nil(t, res)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FindListByID", "listID").Error())
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
		s               ListService
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
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("GetTodayListID", "ownerID").Error())
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
		s               ListService
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
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("GetTomorrowListID", "ownerID").Error())
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
		s          ListService
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
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FindLists", "ownerID").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter pagination cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchLists")
		s = NewListService(m)
		res, err = s.FindLists(ownerID, nil, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FindLists", "pagination").Error())
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
		s                ListService
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
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FindGroupedLists", "ownerID").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter groupID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchGroupedLists")
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, uuid.Nil, pagination, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FindGroupedLists", "groupID").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter pagination cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchGroupedLists")
		s = NewListService(m)
		res, err = s.FindGroupedLists(ownerID, groupID, nil, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FindGroupedLists", "pagination").Error())
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

func TestListService_FindScatteredLists(t *testing.T) {
	defer beQuiet()()
	var (
		m          *listRepositoryMock
		s          ListService
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
		m.On("FetchScatteredLists",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindScatteredLists(ownerID, pagination, "", "")
		assert.Equal(t, current, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchScatteredLists")
		s = NewListService(m)
		res, err = s.FindScatteredLists(uuid.Nil, pagination, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FindScatteredLists", "ownerID").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter pagination cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "FetchScatteredLists")
		s = NewListService(m)
		res, err = s.FindScatteredLists(ownerID, nil, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FindScatteredLists", "pagination").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter needle must be trimmed", func(t *testing.T) {
		var (
			lists  = make([]*model.List, 0)
			needle = "\n		needle 		\n"
		)
		m = new(listRepositoryMock)
		m.On("FetchScatteredLists",
			ownerID.String(), pagination.Page, pagination.RPP,
			strings.Trim(needle, " \n\t"), "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindScatteredLists(ownerID, pagination, needle, "")
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
		m.On("FetchScatteredLists",
			ownerID.String(), expectedPageNumber, pag.RPP, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindScatteredLists(ownerID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)

		/* when page<0 */

		pag.Page = -1
		m = new(listRepositoryMock)
		m.On("FetchScatteredLists",
			ownerID.String(), expectedPageNumber, pag.RPP, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindScatteredLists(ownerID, pag, "", "")
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
		m.On("FetchScatteredLists",
			ownerID.String(), pag.Page, expectedRPPNumber, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindScatteredLists(ownerID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)

		/* when RPP<0 */

		pag.RPP = -1
		m = new(listRepositoryMock)
		m.On("FetchScatteredLists",
			ownerID.String(), pag.Page, expectedRPPNumber, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindScatteredLists(ownerID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter sortBy must be trimmed", func(t *testing.T) {
		var (
			lists  = make([]*model.List, 0)
			sortBy = "\n		+first_name 		\n"
		)
		m = new(listRepositoryMock)
		m.On("FetchScatteredLists",
			ownerID.String(), pagination.Page, pagination.RPP, "",
			strings.Trim(sortBy, " \n\t")).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FindScatteredLists(ownerID, pagination, "", sortBy)
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("FetchScatteredLists",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		s = NewListService(m)
		res, err = s.FindScatteredLists(ownerID, pagination, "", "")
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestListService_DeleteList(t *testing.T) {
	defer beQuiet()()
	var (
		m                        *listRepositoryMock
		s                        ListService
		err                      error
		ownerID, groupID, listID = uuid.New(), uuid.New(), uuid.New()
	)

	t.Run("success for grouped list", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("DeleteList", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		s = NewListService(m)
		err = s.DeleteList(ownerID, groupID, listID)
		assert.NoError(t, err)
	})

	t.Run("success for scattered list (groupID=uuid.Nil)", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("DeleteList", ownerID.String(), "", listID.String()).
			Return(true, nil)
		s = NewListService(m)
		err = s.DeleteList(ownerID, uuid.Nil, listID)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "DeleteList")
		s = NewListService(m)
		err = s.DeleteList(uuid.Nil, groupID, listID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("DeleteList", "ownerID").Error())
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "DeleteList")
		s = NewListService(m)
		err = s.DeleteList(ownerID, groupID, uuid.Nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("DeleteList", "listID").Error())
	})

	t.Run("got a repository error (list could not be deleted)", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("DeleteList", mock.Anything, mock.Anything, mock.Anything).
			Return(false, unexpected)
		s = NewListService(m)
		err = s.DeleteList(ownerID, groupID, listID)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestListService_DuplicateList(t *testing.T) {
	defer beQuiet()()
	var (
		m               *listRepositoryMock
		s               ListService
		res             uuid.UUID
		err             error
		ownerID, listID = uuid.New(), uuid.New()
	)

	t.Run("success", func(t *testing.T) {
		var replicaID = uuid.New()
		m = new(listRepositoryMock)
		m.On("DuplicateList", mock.Anything, mock.Anything, mock.Anything).
			Return(replicaID.String(), nil)
		s = NewListService(m)
		res, err = s.DuplicateList(ownerID, listID)
		assert.Equal(t, replicaID, res)
		assert.NoError(t, err)
	})

	t.Run("got UUID parsing error", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("DuplicateList", mock.Anything, mock.Anything, mock.Anything).
			Return("x", nil)
		s = NewListService(m)
		res, err = s.DuplicateList(ownerID, listID)
		assert.ErrorContains(t, err, "invalid UUID length: 1")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("did parse UUID", func(t *testing.T) {
		var id = uuid.New()
		m = new(listRepositoryMock)
		m.On("DuplicateList", mock.Anything, mock.Anything, mock.Anything).
			Return(id.String(), nil)
		s = NewListService(m)
		res, err = s.DuplicateList(ownerID, listID)
		assert.Equal(t, id, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "DuplicateList")
		s = NewListService(m)
		res, err = s.DuplicateList(uuid.Nil, listID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("DuplicateList", "ownerID").Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "DuplicateList")
		s = NewListService(m)
		res, err = s.DuplicateList(ownerID, uuid.Nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("DuplicateList", "listID").Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("DuplicateList", mock.Anything, mock.Anything, mock.Anything).
			Return("", unexpected)
		s = NewListService(m)
		res, err = s.DuplicateList(ownerID, listID)
		assert.Equal(t, uuid.Nil, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestListService_ConvertToScatteredList(t *testing.T) {
	defer beQuiet()()
	var (
		m               *listRepositoryMock
		s               ListService
		res             bool
		err             error
		ownerID, listID = uuid.New(), uuid.New()
	)

	t.Run("success list", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("ConvertToScatteredList", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		s = NewListService(m)
		res, err = s.ConvertToScatteredList(ownerID, listID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "ConvertToScatteredList")
		s = NewListService(m)
		res, err = s.ConvertToScatteredList(uuid.Nil, listID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("ConvertToScatteredList", "ownerID").Error())
		assert.False(t, res)
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "ConvertToScatteredList")
		s = NewListService(m)
		res, err = s.ConvertToScatteredList(ownerID, uuid.Nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("ConvertToScatteredList", "listID").Error())
		assert.False(t, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("ConvertToScatteredList", mock.Anything, mock.Anything, mock.Anything).
			Return(false, unexpected)
		s = NewListService(m)
		res, err = s.ConvertToScatteredList(ownerID, listID)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}

func TestListService_MoveList(t *testing.T) {
	defer beQuiet()()
	var (
		m                        *listRepositoryMock
		s                        ListService
		res                      bool
		err                      error
		ownerID, listID, groupID = uuid.New(), uuid.New(), uuid.New()
	)

	t.Run("success list", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("MoveList", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		s = NewListService(m)
		res, err = s.MoveList(ownerID, listID, groupID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "MoveList")
		s = NewListService(m)
		res, err = s.MoveList(uuid.Nil, listID, groupID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("MoveList", "ownerID").Error())
		assert.False(t, res)
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "MoveList")
		s = NewListService(m)
		res, err = s.MoveList(ownerID, uuid.Nil, groupID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("MoveList", "listID").Error())
		assert.False(t, res)
	})

	t.Run("parameter targetGroupID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "MoveList")
		s = NewListService(m)
		res, err = s.MoveList(ownerID, listID, uuid.Nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("MoveList", "targetGroupID").Error())
		assert.False(t, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("MoveList", mock.Anything, mock.Anything, mock.Anything).
			Return(false, unexpected)
		s = NewListService(m)
		res, err = s.MoveList(ownerID, listID, groupID)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}

func TestListService_UpdateList(t *testing.T) {
	defer beQuiet()()
	var (
		m                        *listRepositoryMock
		s                        ListService
		res                      bool
		err                      error
		ownerID, groupID, listID = uuid.New(), uuid.New(), uuid.New()
		up                       = &transfer.ListUpdate{
			Name:        "\t   list name\n   ",
			Description: "\n  description  \n",
		}
	)

	t.Run("success for grouped list", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("UpdateList",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		s = NewListService(m)
		res, err = s.UpdateList(ownerID, groupID, listID, up)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("success for scattered list", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.On("UpdateList",
			ownerID.String(), "", listID.String(), up).
			Return(true, nil)
		s = NewListService(m)
		res, err = s.UpdateList(ownerID, uuid.Nil, listID, up)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "UpdateList")
		s = NewListService(m)
		res, err = s.UpdateList(uuid.Nil, groupID, listID, up)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("UpdateList", "ownerID").Error())
		assert.False(t, res)
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "UpdateList")
		s = NewListService(m)
		res, err = s.UpdateList(ownerID, groupID, uuid.Nil, up)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("UpdateList", "listID").Error())
		assert.False(t, res)
	})

	t.Run("parameter up cannot be nil", func(t *testing.T) {
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "UpdateList")
		s = NewListService(m)
		res, err = s.UpdateList(ownerID, groupID, listID, nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("UpdateList", "up").Error())
		assert.False(t, res)
	})

	t.Run("name too long: max length is 50 characters", func(t *testing.T) {
		var previousName = up.Name
		up.Name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxX"
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "UpdateList")
		s = NewListService(m)
		res, err = s.UpdateList(ownerID, groupID, listID, up)
		up.Name = previousName
		assert.ErrorContains(t, err,
			noda.ErrTooLong.Clone().FormatDetails("name", "list", 50).Error())
		assert.False(t, res)
	})

	t.Run("next.Name and next.Description must be trimmed", func(t *testing.T) {
		var previousName, previousDesc = up.Name, up.Description
		m = new(listRepositoryMock)
		m.AssertNotCalled(t, "UpdateList")
		s = NewListService(m)
		m.On("UpdateList",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil)
		s = NewListService(m)
		res, err = s.UpdateList(ownerID, groupID, listID, up)
		assert.Equal(t, "list name", up.Name)
		assert.Equal(t, "description", up.Description)
		up.Name, up.Description = previousName, previousDesc
		assert.False(t, res)
		assert.NoError(t, err)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		m = new(listRepositoryMock)
		m.On("UpdateList",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, unexpected)
		s = NewListService(m)
		res, err = s.UpdateList(ownerID, groupID, listID, up)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}
