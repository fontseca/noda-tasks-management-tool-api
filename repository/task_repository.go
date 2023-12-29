package repository

import (
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"time"
)

type TaskRepository interface {
	Save(ownerID, taskID string, creation *transfer.TaskCreation) (insertedID string, err error)
	Duplicate(ownerID, taskID string) (replicaID string, err error)
	FetchByID(ownerID, listID, taskID string) (task *model.Task, err error)
	Fetch(ownerID, listID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error)
	FetchFromToday(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error)
	FetchFromTomorrow(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error)
	FetchFromDeferred(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error)
	Update(ownerID, listID, taskID string, update *transfer.TaskUpdate) (ok bool, err error)
	Reorder(ownerID, listID, taskID string, position uint64) (ok bool, err error)
	SetReminder(ownerID, listID, taskID string, remindAt time.Time) (ok bool, err error)
	SetPriority(ownerID, listID, taskID string, priority types.TaskPriority) (ok bool, err error)
	SetDueDate(ownerID, listID, taskID string, dueDate time.Time) (ok bool, err error)
	Complete(ownerID, listID, taskID string) (ok bool, err error)
	Resume(ownerID, listID, taskID string) (ok bool, err error)
	Pin(ownerID, listID, taskID string) (ok bool, err error)
	Unpin(ownerID, listID, taskID string) (ok bool, err error)
	Move(ownerID, taskID, targetListID string) (ok bool, err error)
	Today(ownerID, taskID string) (ok bool, err error)
	Tomorrow(ownerID, taskID string) (ok bool, err error)
	Defer(ownerID, taskID string) (ok bool, err error)
	Trash(ownerID, listID, taskID string) (ok bool, err error)
	RestoreFromTrash(ownerID, listID, taskID string) (ok bool, err error)
	Delete(ownerID, listID, taskID string) error
}
