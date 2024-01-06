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
	"time"
)

func TestTaskService_Save(t *testing.T) {
	defer beQuiet()()
	const routine = "Save"
	var (
		ownerID, listID, inserted = uuid.New(), uuid.New(), uuid.New()
		res                       uuid.UUID
		err                       error
	)

	t.Run("success", func(t *testing.T) {
		var c = &transfer.TaskCreation{
			Title:       "title",
			Headline:    "headline",
			Description: "description",
			RemindAt:    time.Now().Add(5 * time.Hour),
			DueDate:     time.Now().Add(10 * time.Hour),
			Status:      types.TaskStatusIncomplete,
			Priority:    types.TaskPriorityHigh,
		}
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, ownerID.String(), listID.String(), c).Return(inserted.String(), nil)
		res, err = NewTaskService(r).Save(ownerID, listID, c)
		assert.Equal(t, inserted, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not nil or uuid.Nil", func(t *testing.T) {
		var c = new(transfer.TaskCreation)

		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Save(uuid.Nil, listID, c)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Save", "ownerID").Error())
		})

		t.Run("\"listID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Save(ownerID, uuid.Nil, c)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Save", "listID").Error())
		})

		t.Run("\"creation\" != nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Save(ownerID, listID, nil)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Save", "creation").Error())
		})
	})

	t.Run("must trim all string fields in \"creation\"", func(t *testing.T) {
		var c = &transfer.TaskCreation{
			Title:       blankset + "Title" + blankset,
			Headline:    blankset + "Headline" + blankset,
			Description: blankset + "Description" + blankset,
		}
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything).Return(inserted.String(), nil)
		res, err = NewTaskService(r).Save(ownerID, listID, c)
		assert.Equal(t, inserted, res)
		assert.Equal(t, "Title", c.Title)
		assert.Equal(t, "Headline", c.Headline)
		assert.Equal(t, "Description", c.Description)
		assert.NoError(t, err)
	})

	t.Run("must default values in \"creation\"", func(t *testing.T) {
		var c = &transfer.TaskCreation{
			Headline:    blankset + "Headline" + blankset,
			Description: blankset + "Description" + blankset,
		}
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything).Return(inserted.String(), nil)
		res, err = NewTaskService(r).Save(ownerID, listID, c)
		assert.Equal(t, inserted, res)
		assert.Equal(t, "Untitled", c.Title)
		assert.Equal(t, types.TaskPriorityMedium, c.Priority)
		assert.Equal(t, types.TaskStatusIncomplete, c.Status)
		assert.NoError(t, err)
	})

	t.Run("satisfies...", func(t *testing.T) {
		var c = new(transfer.TaskCreation)

		t.Run("128 < len(creation.Title)", func(t *testing.T) {
			c.Title = strings.Repeat("x", 129)
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Save(ownerID, listID, c)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Title", "creation", 128).Error())
			assert.Equal(t, uuid.Nil, res)
			c.Title = ""
		})

		t.Run("64 < len(creation.Headline)", func(t *testing.T) {
			c.Headline = strings.Repeat("x", 65)
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Save(ownerID, listID, c)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Headline", "creation", 64).Error())
			assert.Equal(t, uuid.Nil, res)
			c.Headline = ""
		})

		t.Run("512 < len(creation.Description)", func(t *testing.T) {
			c.Description = strings.Repeat("x", 513)
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Save(ownerID, listID, c)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Description", "creation", 512).Error())
			assert.Equal(t, uuid.Nil, res)
			c.Description = ""
		})
	})

	t.Run("got a repository error", func(t *testing.T) {
		var c = new(transfer.TaskCreation)
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything).Return("", unexpected)
		res, err = NewTaskService(r).Save(ownerID, listID, c)
		assert.ErrorIs(t, err, unexpected)
		assert.Equal(t, uuid.Nil, res)
	})
}

func TestTaskService_Duplicate(t *testing.T) {
	defer beQuiet()()
	const routine = "Duplicate"
	var (
		res                        uuid.UUID
		err                        error
		ownerID, taskID, replicaID = uuid.New(), uuid.New(), uuid.New()
	)

	t.Run("success", func(t *testing.T) {
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, ownerID.String(), taskID.String()).Return(replicaID.String(), nil)
		res, err = NewTaskService(r).Duplicate(ownerID, taskID)
		assert.Equal(t, replicaID, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not uuid.Nil", func(t *testing.T) {
		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Duplicate(uuid.Nil, taskID)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Duplicate", "ownerID").Error())
		})

		t.Run("\"taskID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Duplicate(ownerID, uuid.Nil)
			assert.Equal(t, uuid.Nil, res)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Duplicate", "taskID").Error())
		})
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything).Return("", unexpected)
		res, err = NewTaskService(r).Duplicate(ownerID, taskID)
		assert.ErrorIs(t, err, unexpected)
		assert.Equal(t, uuid.Nil, res)
	})
}

func TestTaskService_FetchByID(t *testing.T) {
	defer beQuiet()()
	const routine = "FetchByID"
	var (
		res  *model.Task
		err  error
		task = &model.Task{
			ID:      uuid.New(),
			OwnerID: uuid.New(),
			ListID:  uuid.New(),
		}
	)

	t.Run("success", func(t *testing.T) {
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, task.OwnerID.String(), task.ListID.String(), task.ID.String()).Return(task, nil)
		res, err = NewTaskService(r).FetchByID(task.OwnerID, task.ListID, task.ID)
		assert.Equal(t, task, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not uuid.Nil", func(t *testing.T) {
		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).FetchByID(uuid.Nil, task.ListID, task.ID)
			assert.ErrorContains(t, err, noda.NewNilParameterError("FetchByID", "ownerID").Error())
			assert.Nil(t, res)
		})

		t.Run("\"listID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).FetchByID(task.OwnerID, uuid.Nil, task.ID)
			assert.ErrorContains(t, err, noda.NewNilParameterError("FetchByID", "listID").Error())
			assert.Nil(t, res)
		})

		t.Run("\"taskID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).FetchByID(task.OwnerID, task.ListID, uuid.Nil)
			assert.ErrorContains(t, err, noda.NewNilParameterError("FetchByID", "taskID").Error())
			assert.Nil(t, res)
		})
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything).Return(nil, unexpected)
		res, err = NewTaskService(r).FetchByID(task.OwnerID, task.ListID, task.ID)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestTaskService_Fetch(t *testing.T) {
	defer beQuiet()()
	const routine = "Fetch"
	var (
		ownerID, listID = uuid.New(), uuid.New()
		res             *types.Result[model.Task]
		err             error
		page            int64 = 1
		rpp             int64 = 10
		needle                = "x"
		sortExpr              = "-title"
		pagination            = &types.Pagination{Page: 1, RPP: 10}
		tasks                 = []*model.Task{
			{ID: uuid.New(), OwnerID: ownerID, ListID: listID},
			{ID: uuid.New(), OwnerID: ownerID, ListID: listID},
			{ID: uuid.New(), OwnerID: ownerID, ListID: listID},
		}
	)

	t.Run("success", func(t *testing.T) {
		var result = &types.Result[model.Task]{
			Page:      page,
			RPP:       10,
			Retrieved: int64(len(tasks)),
			Payload:   tasks,
		}
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, ownerID.String(), listID.String(), page, rpp, needle, sortExpr).Return(tasks, nil)
		res, err = NewTaskService(r).Fetch(ownerID, listID, pagination, needle, sortExpr)
		assert.Equal(t, result, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not nil or uuid.Nil", func(t *testing.T) {
		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Fetch(uuid.Nil, listID, pagination, needle, sortExpr)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Fetch", "ownerID").Error())
			assert.Nil(t, res)
		})

		t.Run("\"listID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Fetch(ownerID, uuid.Nil, pagination, needle, sortExpr)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Fetch", "listID").Error())
			assert.Nil(t, res)
		})

		t.Run("\"pagination\" != nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Fetch(ownerID, listID, nil, needle, sortExpr)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Fetch", "pagination").Error())
			assert.Nil(t, res)
		})
	})

	t.Run("parameters are trimmed", func(n *testing.T) {
		t.Run("\"needle\" is trimmed", func(t *testing.T) {
			var n = blankset + needle + blankset
			var r = mocks.NewTaskRepositoryMock()
			r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, needle, mock.Anything).Return(tasks, nil)
			_, _ = NewTaskService(r).Fetch(ownerID, listID, pagination, n, sortExpr)
		})

		t.Run("\"sortExpr\" is trimmed", func(t *testing.T) {
			var s = blankset + sortExpr + blankset
			var r = mocks.NewTaskRepositoryMock()
			r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, sortExpr).Return(tasks, nil)
			_, _ = NewTaskService(r).Fetch(ownerID, listID, pagination, needle, s)
		})
	})

	t.Run("defaults pagination", func(t *testing.T) {
		const expectedPage, expectedRPP int64 = 1, 10
		pagination.Page = -1
		pagination.RPP = 0
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, expectedPage, expectedRPP, mock.Anything, mock.Anything).Return(tasks, nil)
		_, _ = NewTaskService(r).Fetch(ownerID, listID, pagination, needle, sortExpr)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.
			On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		res, err = NewTaskService(r).Fetch(ownerID, listID, pagination, needle, sortExpr)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestTaskService_FetchFromToday(t *testing.T) {
	defer beQuiet()()
	const routine = "FetchFromToday"
	var (
		ownerID    = uuid.New()
		res        *types.Result[model.Task]
		err        error
		page       int64 = 1
		rpp        int64 = 10
		needle           = "x"
		sortExpr         = "-title"
		pagination       = &types.Pagination{Page: 1, RPP: 10}
		tasks            = []*model.Task{
			{ID: uuid.New(), OwnerID: ownerID},
			{ID: uuid.New(), OwnerID: ownerID},
			{ID: uuid.New(), OwnerID: ownerID},
		}
	)

	t.Run("success", func(t *testing.T) {
		var result = &types.Result[model.Task]{
			Page:      page,
			RPP:       10,
			Retrieved: int64(len(tasks)),
			Payload:   tasks,
		}
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, ownerID.String(), page, rpp, needle, sortExpr).Return(tasks, nil)
		res, err = NewTaskService(r).FetchFromToday(ownerID, pagination, needle, sortExpr)
		assert.Equal(t, result, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not nil or uuid.Nil", func(t *testing.T) {
		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).FetchFromToday(uuid.Nil, pagination, needle, sortExpr)
			assert.ErrorContains(t, err, noda.NewNilParameterError("FetchFromToday", "ownerID").Error())
			assert.Nil(t, res)
		})

		t.Run("\"pagination\" != nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).FetchFromToday(ownerID, nil, needle, sortExpr)
			assert.ErrorContains(t, err, noda.NewNilParameterError("FetchFromToday", "pagination").Error())
			assert.Nil(t, res)
		})
	})

	t.Run("parameters are trimmed", func(n *testing.T) {
		t.Run("\"needle\" is trimmed", func(t *testing.T) {
			var n = blankset + needle + blankset
			var r = mocks.NewTaskRepositoryMock()
			r.On(routine, mock.Anything, mock.Anything, mock.Anything, needle, mock.Anything).Return(tasks, nil)
			_, _ = NewTaskService(r).FetchFromToday(ownerID, pagination, n, sortExpr)
		})

		t.Run("\"sortExpr\" is trimmed", func(t *testing.T) {
			var s = blankset + sortExpr + blankset
			var r = mocks.NewTaskRepositoryMock()
			r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, sortExpr).Return(tasks, nil)
			_, _ = NewTaskService(r).FetchFromToday(ownerID, pagination, needle, s)
		})
	})

	t.Run("defaults pagination", func(t *testing.T) {
		const expectedPage, expectedRPP int64 = 1, 10
		pagination.Page = -1
		pagination.RPP = 0
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, expectedPage, expectedRPP, mock.Anything, mock.Anything).Return(tasks, nil)
		_, _ = NewTaskService(r).FetchFromToday(ownerID, pagination, needle, sortExpr)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.
			On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		res, err = NewTaskService(r).FetchFromToday(ownerID, pagination, needle, sortExpr)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestTaskService_FetchFromTomorrow(t *testing.T) {
	defer beQuiet()()
	const routine = "FetchFromTomorrow"
	var (
		ownerID    = uuid.New()
		res        *types.Result[model.Task]
		err        error
		page       int64 = 1
		rpp        int64 = 10
		needle           = "x"
		sortExpr         = "-title"
		pagination       = &types.Pagination{Page: 1, RPP: 10}
		tasks            = []*model.Task{
			{ID: uuid.New(), OwnerID: ownerID},
			{ID: uuid.New(), OwnerID: ownerID},
			{ID: uuid.New(), OwnerID: ownerID},
		}
	)

	t.Run("success", func(t *testing.T) {
		var result = &types.Result[model.Task]{
			Page:      page,
			RPP:       10,
			Retrieved: int64(len(tasks)),
			Payload:   tasks,
		}
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, ownerID.String(), page, rpp, needle, sortExpr).Return(tasks, nil)
		res, err = NewTaskService(r).FetchFromTomorrow(ownerID, pagination, needle, sortExpr)
		assert.Equal(t, result, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not nil or uuid.Nil", func(t *testing.T) {
		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).FetchFromTomorrow(uuid.Nil, pagination, needle, sortExpr)
			assert.ErrorContains(t, err, noda.NewNilParameterError("FetchFromTomorrow", "ownerID").Error())
			assert.Nil(t, res)
		})

		t.Run("\"pagination\" != nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).FetchFromTomorrow(ownerID, nil, needle, sortExpr)
			assert.ErrorContains(t, err, noda.NewNilParameterError("FetchFromTomorrow", "pagination").Error())
			assert.Nil(t, res)
		})
	})

	t.Run("parameters are trimmed", func(n *testing.T) {
		t.Run("\"needle\" is trimmed", func(t *testing.T) {
			var n = blankset + needle + blankset
			var r = mocks.NewTaskRepositoryMock()
			r.On(routine, mock.Anything, mock.Anything, mock.Anything, needle, mock.Anything).Return(tasks, nil)
			_, _ = NewTaskService(r).FetchFromTomorrow(ownerID, pagination, n, sortExpr)
		})

		t.Run("\"sortExpr\" is trimmed", func(t *testing.T) {
			var s = blankset + sortExpr + blankset
			var r = mocks.NewTaskRepositoryMock()
			r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, sortExpr).Return(tasks, nil)
			_, _ = NewTaskService(r).FetchFromTomorrow(ownerID, pagination, needle, s)
		})
	})

	t.Run("defaults pagination", func(t *testing.T) {
		const expectedPage, expectedRPP int64 = 1, 10
		pagination.Page = -1
		pagination.RPP = 0
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, expectedPage, expectedRPP, mock.Anything, mock.Anything).Return(tasks, nil)
		_, _ = NewTaskService(r).FetchFromTomorrow(ownerID, pagination, needle, sortExpr)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.
			On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		res, err = NewTaskService(r).FetchFromTomorrow(ownerID, pagination, needle, sortExpr)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestTaskService_FetchFromDeferred(t *testing.T) {
	defer beQuiet()()
	const routine = "FetchFromDeferred"
	var (
		ownerID    = uuid.New()
		res        *types.Result[model.Task]
		err        error
		page       int64 = 1
		rpp        int64 = 10
		needle           = "x"
		sortExpr         = "-title"
		pagination       = &types.Pagination{Page: 1, RPP: 10}
		tasks            = []*model.Task{
			{ID: uuid.New(), OwnerID: ownerID},
			{ID: uuid.New(), OwnerID: ownerID},
			{ID: uuid.New(), OwnerID: ownerID},
		}
	)

	t.Run("success", func(t *testing.T) {
		var result = &types.Result[model.Task]{
			Page:      page,
			RPP:       10,
			Retrieved: int64(len(tasks)),
			Payload:   tasks,
		}
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, ownerID.String(), page, rpp, needle, sortExpr).Return(tasks, nil)
		res, err = NewTaskService(r).FetchFromDeferred(ownerID, pagination, needle, sortExpr)
		assert.Equal(t, result, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not nil or uuid.Nil", func(t *testing.T) {
		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).FetchFromDeferred(uuid.Nil, pagination, needle, sortExpr)
			assert.ErrorContains(t, err, noda.NewNilParameterError("FetchFromDeferred", "ownerID").Error())
			assert.Nil(t, res)
		})

		t.Run("\"pagination\" != nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).FetchFromDeferred(ownerID, nil, needle, sortExpr)
			assert.ErrorContains(t, err, noda.NewNilParameterError("FetchFromDeferred", "pagination").Error())
			assert.Nil(t, res)
		})
	})

	t.Run("parameters are trimmed", func(n *testing.T) {
		t.Run("\"needle\" is trimmed", func(t *testing.T) {
			var n = blankset + needle + blankset
			var r = mocks.NewTaskRepositoryMock()
			r.On(routine, mock.Anything, mock.Anything, mock.Anything, needle, mock.Anything).Return(tasks, nil)
			_, _ = NewTaskService(r).FetchFromDeferred(ownerID, pagination, n, sortExpr)
		})

		t.Run("\"sortExpr\" is trimmed", func(t *testing.T) {
			var s = blankset + sortExpr + blankset
			var r = mocks.NewTaskRepositoryMock()
			r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, sortExpr).Return(tasks, nil)
			_, _ = NewTaskService(r).FetchFromDeferred(ownerID, pagination, needle, s)
		})
	})

	t.Run("defaults pagination", func(t *testing.T) {
		const expectedPage, expectedRPP int64 = 1, 10
		pagination.Page = -1
		pagination.RPP = 0
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, expectedPage, expectedRPP, mock.Anything, mock.Anything).Return(tasks, nil)
		_, _ = NewTaskService(r).FetchFromDeferred(ownerID, pagination, needle, sortExpr)
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.
			On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		res, err = NewTaskService(r).FetchFromDeferred(ownerID, pagination, needle, sortExpr)
		assert.ErrorIs(t, err, unexpected)
		assert.Nil(t, res)
	})
}

func TestTaskService_Update(t *testing.T) {
	defer beQuiet()()
	const routine = "Update"
	var (
		ownerID, listID, taskID = uuid.New(), uuid.New(), uuid.New()
		res                     bool
		err                     error
	)

	t.Run("success", func(t *testing.T) {
		var u = &transfer.TaskUpdate{
			Title:       "Title",
			Description: "Description",
			Headline:    "Headline",
		}
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, ownerID.String(), listID.String(), taskID.String(), u).Return(true, nil)
		res, err = NewTaskService(r).Update(ownerID, listID, taskID, u)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not nil or uuid.Nil", func(t *testing.T) {
		var u = new(transfer.TaskUpdate)

		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Update(uuid.Nil, listID, taskID, u)
			assert.False(t, res)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Update", "ownerID").Error())
		})

		t.Run("\"listID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Update(ownerID, uuid.Nil, taskID, u)
			assert.False(t, res)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Update", "listID").Error())
		})

		t.Run("\"update\" != nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Update(ownerID, listID, taskID, nil)
			assert.False(t, res)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Update", "update").Error())
		})
	})

	t.Run("trims all string fields in \"update\"", func(t *testing.T) {
		var u = &transfer.TaskUpdate{
			Title:       blankset + "Title" + blankset,
			Headline:    blankset + "Headline" + blankset,
			Description: blankset + "Description" + blankset,
		}
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
		res, err = NewTaskService(r).Update(ownerID, listID, taskID, u)
		assert.Equal(t, "Title", u.Title)
		assert.Equal(t, "Headline", u.Headline)
		assert.Equal(t, "Description", u.Description)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("satisfies...", func(t *testing.T) {
		var u = new(transfer.TaskUpdate)

		t.Run("128 < len(update.Title)", func(t *testing.T) {
			u.Title = strings.Repeat("x", 129)
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Update(ownerID, listID, taskID, u)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Title", "update", 128).Error())
			assert.False(t, res)
			u.Title = ""
		})

		t.Run("64 < len(update.Headline)", func(t *testing.T) {
			u.Headline = strings.Repeat("x", 65)
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Update(ownerID, listID, taskID, u)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Headline", "update", 64).Error())
			assert.False(t, res)
			u.Headline = ""
		})

		t.Run("512 < len(update.Description)", func(t *testing.T) {
			u.Description = strings.Repeat("x", 513)
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Update(ownerID, listID, taskID, u)
			assert.ErrorContains(t, err, noda.ErrTooLong.Clone().FormatDetails("Description", "update", 512).Error())
			assert.False(t, res)
			u.Description = ""
		})
	})

	t.Run("got a repository error", func(t *testing.T) {
		var u = new(transfer.TaskUpdate)
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, unexpected)
		res, err = NewTaskService(r).Update(ownerID, listID, taskID, u)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}

func TestTaskService_Reorder(t *testing.T) {
	defer beQuiet()()
	const routine = "Reorder"
	var (
		ownerID, listID, taskID = uuid.New(), uuid.New(), uuid.New()
		res                     bool
		err                     error
	)

	t.Run("success", func(t *testing.T) {
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, ownerID.String(), listID.String(), taskID.String(), uint64(10)).Return(true, nil)
		res, err = NewTaskService(r).Reorder(ownerID, listID, taskID, 10)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not uuid.Nil", func(t *testing.T) {
		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Reorder(uuid.Nil, listID, taskID, 1)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Reorder", "ownerID").Error())
			assert.False(t, res)
		})

		t.Run("\"listID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Reorder(ownerID, uuid.Nil, taskID, 1)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Reorder", "listID").Error())
			assert.False(t, res)
		})

		t.Run("\"taskID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).Reorder(ownerID, listID, uuid.Nil, 1)
			assert.ErrorContains(t, err, noda.NewNilParameterError("Reorder", "taskID").Error())
			assert.False(t, res)
		})
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, unexpected)
		res, err = NewTaskService(r).Reorder(ownerID, listID, taskID, 1)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}

func TestTaskService_SetReminder(t *testing.T) {
	defer beQuiet()()
	const routine = "SetReminder"
	var (
		ownerID, listID, taskID = uuid.New(), uuid.New(), uuid.New()
		res                     bool
		err                     error
		tm                      = time.Now().Add(5 * time.Hour)
	)

	t.Run("success", func(t *testing.T) {
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, ownerID.String(), listID.String(), taskID.String(), tm).Return(true, nil)
		res, err = NewTaskService(r).SetReminder(ownerID, listID, taskID, tm)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("parameters are not uuid.Nil", func(t *testing.T) {
		t.Run("\"ownerID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).SetReminder(uuid.Nil, listID, taskID, tm)
			assert.ErrorContains(t, err, noda.NewNilParameterError("SetReminder", "ownerID").Error())
			assert.False(t, res)
		})

		t.Run("\"listID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).SetReminder(ownerID, uuid.Nil, taskID, tm)
			assert.ErrorContains(t, err, noda.NewNilParameterError("SetReminder", "listID").Error())
			assert.False(t, res)
		})

		t.Run("\"taskID\" != uuid.Nil", func(t *testing.T) {
			var r = mocks.NewTaskRepositoryMock()
			r.AssertNotCalled(t, routine)
			res, err = NewTaskService(r).SetReminder(ownerID, listID, uuid.Nil, tm)
			assert.ErrorContains(t, err, noda.NewNilParameterError("SetReminder", "taskID").Error())
			assert.False(t, res)
		})
	})

	t.Run("got a repository error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var r = mocks.NewTaskRepositoryMock()
		r.On(routine, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, unexpected)
		res, err = NewTaskService(r).SetReminder(ownerID, listID, taskID, tm)
		assert.ErrorIs(t, err, unexpected)
		assert.False(t, res)
	})
}
