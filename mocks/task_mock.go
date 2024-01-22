package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"time"
)

type TaskRepository struct {
	mock.Mock
}

func NewTaskRepositoryMock() *TaskRepository {
	return new(TaskRepository)
}

func (o *TaskRepository) Save(ownerID, taskID string, creation *transfer.TaskCreation) (insertedID string, err error) {
	var args = o.Called(ownerID, taskID, creation)
	return args.String(0), args.Error(1)
}

func (o *TaskRepository) Duplicate(ownerID, taskID string) (replicaID string, err error) {
	var args = o.Called(ownerID, taskID)
	return args.String(0), args.Error(1)
}

func (o *TaskRepository) FetchByID(ownerID, listID, taskID string) (task *model.Task, err error) {
	var args = o.Called(ownerID, listID, taskID)
	var arg0 = args.Get(0)
	if nil != arg0 {
		task = arg0.(*model.Task)
	}
	return task, args.Error(1)
}

func (o *TaskRepository) Fetch(ownerID, listID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	var args = o.Called(ownerID, listID, page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		tasks = arg0.([]*model.Task)
	}
	return tasks, args.Error(1)
}

func (o *TaskRepository) FetchFromToday(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	var args = o.Called(ownerID, page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		tasks = arg0.([]*model.Task)
	}
	return tasks, args.Error(1)
}

func (o *TaskRepository) FetchFromTomorrow(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	var args = o.Called(ownerID, page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		tasks = arg0.([]*model.Task)
	}
	return tasks, args.Error(1)
}

func (o *TaskRepository) FetchFromDeferred(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	var args = o.Called(ownerID, page, rpp, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		tasks = arg0.([]*model.Task)
	}
	return tasks, args.Error(1)
}

func (o *TaskRepository) Update(ownerID, listID, taskID string, update *transfer.TaskUpdate) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID, update)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Reorder(ownerID, listID, taskID string, position uint64) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID, position)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) SetReminder(ownerID, listID, taskID string, remindAt time.Time) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID, remindAt)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) SetPriority(ownerID, listID, taskID string, priority types.TaskPriority) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID, priority)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) SetDueDate(ownerID, listID, taskID string, dueDate time.Time) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID, dueDate)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Complete(ownerID, listID, taskID string) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Resume(ownerID, listID, taskID string) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Pin(ownerID, listID, taskID string) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Unpin(ownerID, listID, taskID string) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Move(ownerID, taskID, targetListID string) (ok bool, err error) {
	var args = o.Called(ownerID, taskID, targetListID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Today(ownerID, taskID string) (ok bool, err error) {
	var args = o.Called(ownerID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Tomorrow(ownerID, taskID string) (ok bool, err error) {
	var args = o.Called(ownerID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Defer(ownerID, taskID string) (ok bool, err error) {
	var args = o.Called(ownerID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Trash(ownerID, listID, taskID string) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) RestoreFromTrash(ownerID, listID, taskID string) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskRepository) Delete(ownerID, listID, taskID string) error {
	var args = o.Called(ownerID, listID, taskID)
	return args.Error(0)
}

type TaskServiceMock struct {
	mock.Mock
}

func NewTaskServiceMock() *TaskServiceMock {
	return new(TaskServiceMock)
}

func (o *TaskServiceMock) Save(ownerID, listID uuid.UUID, creation *transfer.TaskCreation) (insertedID uuid.UUID, err error) {
	var args = o.Called(ownerID, listID, creation)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *TaskServiceMock) Duplicate(ownerID, taskID uuid.UUID) (replicaID uuid.UUID, err error) {
	var args = o.Called(ownerID, taskID)
	var arg0 = args.Get(0)
	if nil != arg0 {
		replicaID = arg0.(uuid.UUID)
	}
	return replicaID, args.Error(1)
}

func (o *TaskServiceMock) FetchByID(ownerID, listID, taskID uuid.UUID) (task *model.Task, err error) {
	var args = o.Called(ownerID, listID, taskID)
	var arg0 = args.Get(0)
	if nil != arg0 {
		task = arg0.(*model.Task)
	}
	return task, args.Error(1)
}

func (o *TaskServiceMock) Fetch(ownerID, listID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error) {
	var args = o.Called(ownerID, listID, pagination, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		result = arg0.(*types.Result[model.Task])
	}
	return result, args.Error(1)
}

func (o *TaskServiceMock) FetchFromToday(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error) {
	var args = o.Called(ownerID, pagination, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		result = arg0.(*types.Result[model.Task])
	}
	return result, args.Error(1)
}

func (o *TaskServiceMock) FetchFromTomorrow(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error) {
	var args = o.Called(ownerID, pagination, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		result = arg0.(*types.Result[model.Task])
	}
	return result, args.Error(1)
}

func (o *TaskServiceMock) FetchFromDeferred(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error) {
	var args = o.Called(ownerID, pagination, needle, sortExpr)
	var arg0 = args.Get(0)
	if nil != arg0 {
		result = arg0.(*types.Result[model.Task])
	}
	return result, args.Error(1)
}

func (o *TaskServiceMock) Update(ownerID, listID, taskID uuid.UUID, update *transfer.TaskUpdate) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Reorder(ownerID, listID, taskID uuid.UUID, position uint64) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) SetReminder(ownerID, listID, taskID uuid.UUID, remindAt time.Time) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) SetPriority(ownerID, listID, taskID uuid.UUID, priority types.TaskPriority) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) SetDueDate(ownerID, listID, taskID uuid.UUID, dueDate time.Time) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Complete(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Resume(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Pin(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Unpin(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Move(ownerID, taskID, targetListID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Today(ownerID, taskID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Tomorrow(ownerID, taskID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Defer(ownerID, taskID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Trash(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) RestoreFromTrash(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, taskID)
	return args.Bool(0), args.Error(1)
}

func (o *TaskServiceMock) Delete(ownerID, listID, taskID uuid.UUID) error {
	var args = o.Called(ownerID, listID, taskID)
	return args.Error(0)
}
