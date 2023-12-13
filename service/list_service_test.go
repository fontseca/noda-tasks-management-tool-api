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
	"noda/mocks"
	"strings"
	"testing"
)

func TestListService_Save(t *testing.T) {
	defer beQuiet()()
	var (
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
		var m = mocks.NewListRepositoryMock()
		m.On("Save", ownerID.String(), groupID.String(), next).
			Return(insertedID.String(), nil)
		s = NewListService(m)
		res, err = s.Save(ownerID, groupID, next)
		assert.Equal(t, insertedID, res)
		assert.NoError(t, err)
	})

	t.Run("success for scattered list", func(t *testing.T) {
		insertedID := uuid.New()
		var m = mocks.NewListRepositoryMock()
		m.On("Save", ownerID.String(), "", next).
			Return(insertedID.String(), nil)
		s = NewListService(m)
		res, err = s.Save(ownerID, uuid.Nil, next)
		assert.Equal(t, insertedID, res)
		assert.NoError(t, err)
	})

	t.Run("got UUID parsing error", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("Save", mock.Anything, mock.Anything, mock.Anything).
			Return("x", nil)
		s = NewListService(m)
		res, err = s.Save(ownerID, groupID, next)
		assert.ErrorContains(t, err, "invalid UUID length: 1")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("did parse UUID", func(t *testing.T) {
		parsed := uuid.MustParse("4fedb41f-5e44-4e63-9266-4b094bd7ba2d")
		var m = mocks.NewListRepositoryMock()
		m.On("Save", mock.Anything, mock.Anything, mock.Anything).
			Return(parsed.String(), nil)
		s = NewListService(m)
		res, err = s.Save(ownerID, groupID, next)
		assert.Equal(t, parsed, res)
		assert.NoError(t, err)
	})

	t.Run("name cannot be empty", func(t *testing.T) {
		var previousName = next.Name
		next.Name = "  		  \n"
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Save")
		s = NewListService(m)
		res, err = s.Save(ownerID, groupID, next)
		next.Name = previousName
		assert.ErrorContains(t, err, "name cannot be an empty string")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Save")
		s = NewListService(m)
		res, err = s.Save(uuid.Nil, groupID, next)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Save", "ownerID").Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("parameter next cannot be nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Save")
		s = NewListService(m)
		res, err = s.Save(ownerID, groupID, nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Save", "creation").Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("name too long: max length is 32 characters", func(t *testing.T) {
		var previousName = next.Name
		next.Name = strings.Repeat("x", 1+32)
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Save")
		s = NewListService(m)
		res, err = s.Save(ownerID, groupID, next)
		next.Name = previousName
		assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("name", "list", 32).Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("description too long: max length is 512 characters", func(t *testing.T) {
		var description = next.Description
		next.Description = strings.Repeat("x", 1+512)
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Save")
		s = NewListService(m)
		res, err = s.Save(ownerID, groupID, next)
		next.Description = description
		assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("description", "list", 512).Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("next.Name and next.Description must be trimmed", func(t *testing.T) {
		var previousName, previousDesc = next.Name, next.Description
		var insertedID = uuid.New()
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Save")
		s = NewListService(m)
		m.On("Save", mock.Anything, mock.Anything, mock.Anything).
			Return(insertedID.String(), nil)
		s = NewListService(m)
		res, err = s.Save(ownerID, groupID, next)
		assert.Equal(t, "list name", next.Name)
		assert.Equal(t, "description", next.Description)
		next.Name, next.Description = previousName, previousDesc
		assert.Equal(t, insertedID, res)
		assert.NoError(t, err)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("Save", mock.Anything, mock.Anything, mock.Anything).
			Return("", unexpected)
		s = NewListService(m)
		res, err = s.Save(ownerID, groupID, next)
		assert.ErrorIs(t, err, unexpected)
		assert.Equal(t, uuid.Nil, res)
	})
}

func TestListService_FetchByID(t *testing.T) {
	defer beQuiet()()
	var (
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
		var m = mocks.NewListRepositoryMock()
		m.On("FetchByID", ownerID.String(), groupID.String(), listID.String()).
			Return(actual, nil)
		s = NewListService(m)
		res, err = s.FetchByID(ownerID, groupID, listID)
		assert.NoError(t, err)
		assert.Equal(t, actual, res)
	})

	t.Run("success for scattered list", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("FetchByID", ownerID.String(), "", listID.String()).
			Return(actual, nil)
		s = NewListService(m)
		res, err = s.FetchByID(ownerID, uuid.Nil, listID)
		assert.NoError(t, err)
		assert.Equal(t, actual, res)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "FetchByID")
		s = NewListService(m)
		res, err = s.FetchByID(uuid.Nil, groupID, listID)
		assert.Nil(t, res)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FetchByID", "ownerID").Error())
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "FetchByID")
		s = NewListService(m)
		res, err = s.FetchByID(ownerID, groupID, uuid.Nil)
		assert.Nil(t, res)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FetchByID", "listID").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("FetchByID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		s = NewListService(m)
		res, err = s.FetchByID(ownerID, groupID, listID)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestListService_GetTodayListID(t *testing.T) {
	defer beQuiet()()
	var (
		s               ListService
		res             uuid.UUID
		err             error
		ownerID, listID = uuid.New(), uuid.New()
	)

	t.Run("success", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("GetTodayListID", mock.Anything).
			Return(listID.String(), nil)
		s = NewListService(m)
		res, err = s.GetTodayListID(ownerID)
		assert.Equal(t, listID, res)
		assert.NoError(t, err)
	})

	t.Run("got UUID parsing error", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("GetTodayListID", mock.Anything).
			Return("x", nil)
		s = NewListService(m)
		res, err = s.GetTodayListID(ownerID)
		assert.ErrorContains(t, err, "invalid UUID length: 1")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("did parse UUID", func(t *testing.T) {
		var id = uuid.MustParse("4fedb41f-5e44-4e63-9266-4b094bd7ba2d")
		var m = mocks.NewListRepositoryMock()
		m.On("GetTodayListID", mock.Anything).
			Return(id.String(), nil)
		s = NewListService(m)
		res, err = s.GetTodayListID(ownerID)
		assert.Equal(t, id, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "GetTodayListID")
		s = NewListService(m)
		res, err = s.GetTodayListID(uuid.Nil)
		assert.Equal(t, uuid.Nil, res)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("GetTodayListID", "ownerID").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
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
		s               ListService
		res             uuid.UUID
		err             error
		ownerID, listID = uuid.New(), uuid.New()
	)

	t.Run("success", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("GetTomorrowListID", mock.Anything).
			Return(listID.String(), nil)
		s = NewListService(m)
		res, err = s.GetTomorrowListID(ownerID)
		assert.Equal(t, listID, res)
		assert.NoError(t, err)
	})

	t.Run("got UUID parsing error", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("GetTomorrowListID", mock.Anything).
			Return("x", nil)
		s = NewListService(m)
		res, err = s.GetTomorrowListID(ownerID)
		assert.ErrorContains(t, err, "invalid UUID length: 1")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("did parse UUID", func(t *testing.T) {
		var id = uuid.MustParse("4fedb41f-5e44-4e63-9266-4b094bd7ba2d")
		var m = mocks.NewListRepositoryMock()
		m.On("GetTomorrowListID", mock.Anything).
			Return(id.String(), nil)
		s = NewListService(m)
		res, err = s.GetTomorrowListID(ownerID)
		assert.Equal(t, id, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "GetTomorrowListID")
		s = NewListService(m)
		res, err = s.GetTomorrowListID(uuid.Nil)
		assert.Equal(t, uuid.Nil, res)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("GetTomorrowListID", "ownerID").Error())
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("GetTomorrowListID", mock.Anything).
			Return("", unexpected)
		s = NewListService(m)
		res, err = s.GetTomorrowListID(ownerID)
		assert.ErrorIs(t, err, unexpected)
		assert.Equal(t, uuid.Nil, res)
	})
}

func TestListService_Fetch(t *testing.T) {
	defer beQuiet()()
	var (
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
		var m = mocks.NewListRepositoryMock()
		m.On("Fetch",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.Fetch(ownerID, pagination, "", "")
		assert.Equal(t, current, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Fetch")
		s = NewListService(m)
		res, err = s.Fetch(uuid.Nil, pagination, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Fetch", "ownerID").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter pagination cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Fetch")
		s = NewListService(m)
		res, err = s.Fetch(ownerID, nil, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Fetch", "pagination").Error())
		assert.Nil(t, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("Fetch",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		s = NewListService(m)
		res, err = s.Fetch(ownerID, pagination, "", "")
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestListService_FetchGrouped(t *testing.T) {
	defer beQuiet()()
	var (
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
		var m = mocks.NewListRepositoryMock()
		m.On("FetchGrouped",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, groupID, pagination, "", "")
		assert.Equal(t, current, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "FetchGrouped")
		s = NewListService(m)
		res, err = s.FetchGrouped(uuid.Nil, groupID, pagination, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FetchGrouped", "ownerID").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter groupID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "FetchGrouped")
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, uuid.Nil, pagination, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FetchGrouped", "groupID").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter pagination cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "FetchGrouped")
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, groupID, nil, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FetchGrouped", "pagination").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter needle must be trimmed", func(t *testing.T) {
		var (
			lists  = make([]*model.List, 0)
			needle = "\n		needle 		\n"
		)
		var m = mocks.NewListRepositoryMock()
		m.On("FetchGrouped",
			ownerID.String(), groupID.String(), pagination.Page, pagination.RPP,
			strings.Trim(needle, " \n\t"), "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, groupID, pagination, needle, "")
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

		var m = mocks.NewListRepositoryMock()
		m.On("FetchGrouped",
			ownerID.String(), groupID.String(), expectedPageNumber, pag.RPP, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, groupID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)

		/* when page<0 */

		pag.Page = -1
		m = mocks.NewListRepositoryMock()
		m.On("FetchGrouped",
			ownerID.String(), groupID.String(), expectedPageNumber, pag.RPP, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, groupID, pag, "", "")
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

		var m = mocks.NewListRepositoryMock()
		m.On("FetchGrouped",
			ownerID.String(), groupID.String(), pag.Page, expectedRPPNumber, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, groupID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)

		/* when RPP<0 */

		pag.RPP = -1
		m = mocks.NewListRepositoryMock()
		m.On("FetchGrouped",
			ownerID.String(), groupID.String(), pag.Page, expectedRPPNumber, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, groupID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter sortExpr must be trimmed", func(t *testing.T) {
		var (
			lists    = make([]*model.List, 0)
			sortExpr = "\n		+first_name 		\n"
		)
		var m = mocks.NewListRepositoryMock()
		m.On("FetchGrouped",
			ownerID.String(), groupID.String(), pagination.Page, pagination.RPP, "",
			strings.Trim(sortExpr, " \n\t")).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, groupID, pagination, "", sortExpr)
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("FetchGrouped",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		s = NewListService(m)
		res, err = s.FetchGrouped(ownerID, groupID, pagination, "", "")
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestListService_FetchScattered(t *testing.T) {
	defer beQuiet()()
	var (
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
		var m = mocks.NewListRepositoryMock()
		m.On("FetchScattered",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchScattered(ownerID, pagination, "", "")
		assert.Equal(t, current, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "FetchScattered")
		s = NewListService(m)
		res, err = s.FetchScattered(uuid.Nil, pagination, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FetchScattered", "ownerID").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter pagination cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "FetchScattered")
		s = NewListService(m)
		res, err = s.FetchScattered(ownerID, nil, "", "")
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("FetchScattered", "pagination").Error())
		assert.Nil(t, res)
	})

	t.Run("parameter needle must be trimmed", func(t *testing.T) {
		var (
			lists  = make([]*model.List, 0)
			needle = "\n		needle 		\n"
		)
		var m = mocks.NewListRepositoryMock()
		m.On("FetchScattered",
			ownerID.String(), pagination.Page, pagination.RPP,
			strings.Trim(needle, " \n\t"), "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchScattered(ownerID, pagination, needle, "")
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

		var m = mocks.NewListRepositoryMock()
		m.On("FetchScattered",
			ownerID.String(), expectedPageNumber, pag.RPP, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchScattered(ownerID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)

		/* when page<0 */

		pag.Page = -1
		m = mocks.NewListRepositoryMock()
		m.On("FetchScattered",
			ownerID.String(), expectedPageNumber, pag.RPP, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchScattered(ownerID, pag, "", "")
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

		var m = mocks.NewListRepositoryMock()
		m.On("FetchScattered",
			ownerID.String(), pag.Page, expectedRPPNumber, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchScattered(ownerID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)

		/* when RPP<0 */

		pag.RPP = -1
		m = mocks.NewListRepositoryMock()
		m.On("FetchScattered",
			ownerID.String(), pag.Page, expectedRPPNumber, "", "").
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchScattered(ownerID, pag, "", "")
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter sortExpr must be trimmed", func(t *testing.T) {
		var (
			lists    = make([]*model.List, 0)
			sortExpr = "\n		+first_name 		\n"
		)
		var m = mocks.NewListRepositoryMock()
		m.On("FetchScattered",
			ownerID.String(), pagination.Page, pagination.RPP, "",
			strings.Trim(sortExpr, " \n\t")).
			Return(lists, nil)
		s = NewListService(m)
		res, err = s.FetchScattered(ownerID, pagination, "", sortExpr)
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("FetchScattered",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		s = NewListService(m)
		res, err = s.FetchScattered(ownerID, pagination, "", "")
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestListService_Remove(t *testing.T) {
	defer beQuiet()()
	var (
		s                        ListService
		err                      error
		ownerID, groupID, listID = uuid.New(), uuid.New(), uuid.New()
	)

	t.Run("success for grouped list", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("Remove", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		s = NewListService(m)
		err = s.Remove(ownerID, groupID, listID)
		assert.NoError(t, err)
	})

	t.Run("success for scattered list (groupID=uuid.Nil)", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("Remove", ownerID.String(), "", listID.String()).
			Return(true, nil)
		s = NewListService(m)
		err = s.Remove(ownerID, uuid.Nil, listID)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Remove")
		s = NewListService(m)
		err = s.Remove(uuid.Nil, groupID, listID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Remove", "ownerID").Error())
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Remove")
		s = NewListService(m)
		err = s.Remove(ownerID, groupID, uuid.Nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Remove", "listID").Error())
	})

	t.Run("got a repository error (list could not be deleted)", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("Remove", mock.Anything, mock.Anything, mock.Anything).
			Return(false, unexpected)
		s = NewListService(m)
		err = s.Remove(ownerID, groupID, listID)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestListService_Duplicate(t *testing.T) {
	defer beQuiet()()
	var (
		s               ListService
		res             uuid.UUID
		err             error
		ownerID, listID = uuid.New(), uuid.New()
	)

	t.Run("success", func(t *testing.T) {
		var replicaID = uuid.New()
		var m = mocks.NewListRepositoryMock()
		m.On("Duplicate", mock.Anything, mock.Anything, mock.Anything).
			Return(replicaID.String(), nil)
		s = NewListService(m)
		res, err = s.Duplicate(ownerID, listID)
		assert.Equal(t, replicaID, res)
		assert.NoError(t, err)
	})

	t.Run("got UUID parsing error", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("Duplicate", mock.Anything, mock.Anything, mock.Anything).
			Return("x", nil)
		s = NewListService(m)
		res, err = s.Duplicate(ownerID, listID)
		assert.ErrorContains(t, err, "invalid UUID length: 1")
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("did parse UUID", func(t *testing.T) {
		var id = uuid.New()
		var m = mocks.NewListRepositoryMock()
		m.On("Duplicate", mock.Anything, mock.Anything, mock.Anything).
			Return(id.String(), nil)
		s = NewListService(m)
		res, err = s.Duplicate(ownerID, listID)
		assert.Equal(t, id, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Duplicate")
		s = NewListService(m)
		res, err = s.Duplicate(uuid.Nil, listID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Duplicate", "ownerID").Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Duplicate")
		s = NewListService(m)
		res, err = s.Duplicate(ownerID, uuid.Nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Duplicate", "listID").Error())
		assert.Equal(t, uuid.Nil, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("Duplicate", mock.Anything, mock.Anything, mock.Anything).
			Return("", unexpected)
		s = NewListService(m)
		res, err = s.Duplicate(ownerID, listID)
		assert.Equal(t, uuid.Nil, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestListService_Scatter(t *testing.T) {
	defer beQuiet()()
	var (
		s               ListService
		res             bool
		err             error
		ownerID, listID = uuid.New(), uuid.New()
	)

	t.Run("success list", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("Scatter", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		s = NewListService(m)
		res, err = s.Scatter(ownerID, listID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Scatter")
		s = NewListService(m)
		res, err = s.Scatter(uuid.Nil, listID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Scatter", "ownerID").Error())
		assert.False(t, res)
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Scatter")
		s = NewListService(m)
		res, err = s.Scatter(ownerID, uuid.Nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Scatter", "listID").Error())
		assert.False(t, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("Scatter", mock.Anything, mock.Anything, mock.Anything).
			Return(false, unexpected)
		s = NewListService(m)
		res, err = s.Scatter(ownerID, listID)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}

func TestListService_Move(t *testing.T) {
	defer beQuiet()()
	var (
		s                        ListService
		res                      bool
		err                      error
		ownerID, listID, groupID = uuid.New(), uuid.New(), uuid.New()
	)

	t.Run("success list", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("Move", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		s = NewListService(m)
		res, err = s.Move(ownerID, listID, groupID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Move")
		s = NewListService(m)
		res, err = s.Move(uuid.Nil, listID, groupID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Move", "ownerID").Error())
		assert.False(t, res)
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Move")
		s = NewListService(m)
		res, err = s.Move(ownerID, uuid.Nil, groupID)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Move", "listID").Error())
		assert.False(t, res)
	})

	t.Run("parameter targetGroupID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Move")
		s = NewListService(m)
		res, err = s.Move(ownerID, listID, uuid.Nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Move", "targetGroupID").Error())
		assert.False(t, res)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("Move", mock.Anything, mock.Anything, mock.Anything).
			Return(false, unexpected)
		s = NewListService(m)
		res, err = s.Move(ownerID, listID, groupID)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}

func TestListService_Update(t *testing.T) {
	defer beQuiet()()
	var (
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
		var m = mocks.NewListRepositoryMock()
		m.On("Update",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)
		s = NewListService(m)
		res, err = s.Update(ownerID, groupID, listID, up)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("success for scattered list", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.On("Update",
			ownerID.String(), "", listID.String(), up).
			Return(true, nil)
		s = NewListService(m)
		res, err = s.Update(ownerID, uuid.Nil, listID, up)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameter ownerID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Update")
		s = NewListService(m)
		res, err = s.Update(uuid.Nil, groupID, listID, up)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Update", "ownerID").Error())
		assert.False(t, res)
	})

	t.Run("parameter listID cannot be uuid.Nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Update")
		s = NewListService(m)
		res, err = s.Update(ownerID, groupID, uuid.Nil, up)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Update", "listID").Error())
		assert.False(t, res)
	})

	t.Run("parameter up cannot be nil", func(t *testing.T) {
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Update")
		s = NewListService(m)
		res, err = s.Update(ownerID, groupID, listID, nil)
		assert.ErrorContains(t, err,
			noda.NewNilParameterError("Update", "up").Error())
		assert.False(t, res)
	})

	t.Run("name too long: max length is 32 characters", func(t *testing.T) {
		var previousName = up.Name
		up.Name = strings.Repeat("x", 1+32)
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Update")
		s = NewListService(m)
		res, err = s.Update(ownerID, groupID, listID, up)
		up.Name = previousName
		assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("name", "list", 32).Error())
		assert.False(t, res)
	})

	t.Run("description too long: max length is 512 characters", func(t *testing.T) {
		var previousDescription = up.Description
		up.Description = strings.Repeat("x", 1+512)
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Update")
		s = NewListService(m)
		res, err = s.Update(ownerID, groupID, listID, up)
		up.Description = previousDescription
		assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("description", "list", 512).Error())
		assert.False(t, res)
	})

	t.Run("next.Name and next.Description must be trimmed", func(t *testing.T) {
		var previousName, previousDesc = up.Name, up.Description
		var m = mocks.NewListRepositoryMock()
		m.AssertNotCalled(t, "Update")
		s = NewListService(m)
		m.On("Update",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil)
		s = NewListService(m)
		res, err = s.Update(ownerID, groupID, listID, up)
		assert.Equal(t, "list name", up.Name)
		assert.Equal(t, "description", up.Description)
		up.Name, up.Description = previousName, previousDesc
		assert.False(t, res)
		assert.NoError(t, err)
	})

	t.Run("got a repository error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewListRepositoryMock()
		m.On("Update",
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, unexpected)
		s = NewListService(m)
		res, err = s.Update(ownerID, groupID, listID, up)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}
