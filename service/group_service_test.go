package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/failure"
	"noda/mocks"
	"strings"
	"testing"
)

func TestGroupService_SaveGroup(t *testing.T) {
	var (
		ownerID = uuid.New()
		next    = new(transfer.GroupCreation)
		s       GroupService
		res     uuid.UUID
		err     error
	)

	t.Run("success", func(t *testing.T) {
		var m = mocks.NewGroupRepositoryMock()
		m.On("Save", ownerID.String(), next).
			Return(ownerID.String(), nil)
		s = NewGroupService(m)
		res, err = s.Save(ownerID, next)
		assert.Equal(t, ownerID, res)
		assert.NoError(t, err)
	})

	t.Run("name too long: max length is 32 characters", func(t *testing.T) {
		var previousName = next.Name
		next.Name = strings.Repeat("x", 1+32)
		var m = mocks.NewGroupRepositoryMock()
		m.AssertNotCalled(t, "Save")
		s = NewGroupService(m)
		res, err = s.Save(ownerID, next)
		assert.Equal(t, uuid.Nil, res)
		assert.ErrorContains(t, err, failure.ErrTooLong.Clone().FormatDetails("name", "group", 32).Error())
		next.Name = previousName
	})

	t.Run("description too long: max length is 512 characters", func(t *testing.T) {
		var previousDescription = next.Description
		next.Description = strings.Repeat("x", 1+512)
		var m = mocks.NewGroupRepositoryMock()
		m.AssertNotCalled(t, "Save")
		s = NewGroupService(m)
		res, err = s.Save(ownerID, next)
		assert.ErrorContains(t, err, failure.ErrTooLong.Clone().FormatDetails("description", "group", 512).Error())
		assert.Equal(t, uuid.Nil, res)
		next.Description = previousDescription
	})

	t.Run("got an error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewGroupRepositoryMock()
		m.On("Save", ownerID.String(), next).Return("", unexpected)
		s = NewGroupService(m)
		res, err = s.Save(ownerID, next)
		assert.Equal(t, uuid.Nil, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestGroupService_FindGroupByID(t *testing.T) {
	var (
		ownerID, groupID = uuid.New(), uuid.New()
		s                GroupService
		res              *model.Group
		err              error
	)

	t.Run("success", func(t *testing.T) {
		current := new(model.Group)
		var m = mocks.NewGroupRepositoryMock()
		m.On("FetchByID", ownerID.String(), groupID.String()).
			Return(current, nil)
		s = NewGroupService(m)
		res, err = s.FetchByID(ownerID, groupID)
		assert.Equal(t, current, res)
		assert.NoError(t, err)
	})

	t.Run("got an error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewGroupRepositoryMock()
		m.On("FetchByID", ownerID.String(), groupID.String()).
			Return(nil, unexpected)
		s = NewGroupService(m)
		res, err = s.FetchByID(ownerID, groupID)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestGroupService_FindGroups(t *testing.T) {
	var (
		ownerID = uuid.New()
		s       GroupService
		err     error
		res     *types.Result[model.Group]
		pag     = &types.Pagination{Page: 1, RPP: 10}
	)

	t.Run("success", func(t *testing.T) {
		var groups = make([]*model.Group, 0)
		current := &types.Result[model.Group]{
			Page:      1,
			RPP:       10,
			Payload:   groups,
			Retrieved: int64(len(groups)),
		}
		var m = mocks.NewGroupRepositoryMock()
		m.On("Fetch", ownerID.String(), pag.Page, pag.RPP, "", "").
			Return(groups, nil)
		s = NewGroupService(m)
		res, err = s.Fetch(ownerID, pag, "", "")
		assert.Equal(t, current, res)
		assert.NoError(t, err)
	})

	t.Run("got an error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewGroupRepositoryMock()
		m.On("Fetch", ownerID.String(), pag.Page, pag.RPP, "", "").
			Return(nil, unexpected)
		s = NewGroupService(m)
		res, err = s.Fetch(ownerID, pag, "", "")
		assert.Nil(t, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestGroupService_UpdateGroup(t *testing.T) {
	var (
		ownerID, groupID = uuid.New(), uuid.New()
		s                GroupService
		res              bool
		err              error
		up               = new(transfer.GroupUpdate)
	)

	t.Run("success", func(t *testing.T) {
		var m = mocks.NewGroupRepositoryMock()
		m.On("Update", ownerID.String(), groupID.String(), up).
			Return(true, nil)
		s = NewGroupService(m)
		res, err = s.Update(ownerID, groupID, up)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("name too long: max length is 32 characters", func(t *testing.T) {
		var previousName = up.Name
		up.Name = strings.Repeat("x", 1+32)
		var m = mocks.NewGroupRepositoryMock()
		m.AssertNotCalled(t, "Update")
		s = NewGroupService(m)
		res, err = s.Update(ownerID, groupID, up)
		assert.False(t, res)
		assert.ErrorContains(t, err, failure.ErrTooLong.Clone().FormatDetails("name", "group", 32).Error())
		up.Name = previousName
	})

	t.Run("description too long: max length is 512 characters", func(t *testing.T) {
		var previousDescription = up.Description
		up.Description = strings.Repeat("x", 1+512)
		var m = mocks.NewGroupRepositoryMock()
		m.AssertNotCalled(t, "Update")
		s = NewGroupService(m)
		res, err = s.Update(ownerID, groupID, up)
		assert.ErrorContains(t, err, failure.ErrTooLong.Clone().FormatDetails("description", "group", 512).Error())
		assert.False(t, res)
		up.Description = previousDescription
	})

	t.Run("got an error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewGroupRepositoryMock()
		m.On("Update", ownerID.String(), groupID.String(), up).
			Return(false, unexpected)
		s = NewGroupService(m)
		res, err = s.Update(ownerID, groupID, up)
		assert.False(t, res)
		assert.ErrorIs(t, err, unexpected)
	})
}

func TestGroupService_DeleteGroup(t *testing.T) {
	var (
		ownerID, groupID = uuid.New(), uuid.New()
		s                GroupService
		res              bool
		err              error
	)

	t.Run("success", func(t *testing.T) {
		var m = mocks.NewGroupRepositoryMock()
		m.On("Remove", ownerID.String(), groupID.String()).
			Return(true, nil)
		s = NewGroupService(m)
		res, err = s.Remove(ownerID, groupID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("got an error", func(t *testing.T) {
		unexpected := errors.New("unexpected error")
		var m = mocks.NewGroupRepositoryMock()
		m.On("Remove", ownerID.String(), groupID.String()).
			Return(false, unexpected)
		s = NewGroupService(m)
		res, err = s.Remove(ownerID, groupID)
		assert.False(t, res)
		assert.ErrorIs(t, err, unexpected)
	})
}
