package service

import (
	"github.com/google/uuid"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"time"
)

type TaskService interface {
	Save(ownerID, taskID uuid.UUID, creation *transfer.TaskCreation) (insertedID uuid.UUID, err error)
	Duplicate(ownerID, taskID uuid.UUID) (replicaID uuid.UUID, err error)
	FetchByID(ownerID, listID, taskID uuid.UUID) (task *model.Task, err error)
	Fetch(ownerID, listID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error)
	FetchFromToday(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error)
	FetchFromTomorrow(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error)
	FetchFromDeferred(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error)
	Update(ownerID, listID, taskID uuid.UUID, update *transfer.TaskUpdate) (ok bool, err error)
	Reorder(ownerID, listID, taskID uuid.UUID, position uint64) (ok bool, err error)
	SetReminder(ownerID, listID, taskID uuid.UUID, remindAt time.Time) (ok bool, err error)
	SetPriority(ownerID, listID, taskID uuid.UUID, priority types.TaskPriority) (ok bool, err error)
	SetDueDate(ownerID, listID, taskID uuid.UUID, dueDate time.Time) (ok bool, err error)
	Complete(ownerID, listID, taskID uuid.UUID) (ok bool, err error)
	Resume(ownerID, listID, taskID uuid.UUID) (ok bool, err error)
	Pin(ownerID, listID, taskID uuid.UUID) (ok bool, err error)
	Unpin(ownerID, listID, taskID uuid.UUID) (ok bool, err error)
	Move(ownerID, taskID, targetListID uuid.UUID) (ok bool, err error)
	Today(ownerID, taskID uuid.UUID) (ok bool, err error)
	Tomorrow(ownerID, taskID uuid.UUID) (ok bool, err error)
	Defer(ownerID, taskID uuid.UUID) (ok bool, err error)
	Trash(ownerID, listID, taskID uuid.UUID) (ok bool, err error)
	RestoreFromTrash(ownerID, listID, taskID uuid.UUID) (ok bool, err error)
	Delete(ownerID, listID, taskID uuid.UUID) error
}
