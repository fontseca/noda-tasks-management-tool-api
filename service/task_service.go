package service

import (
	"github.com/google/uuid"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/repository"
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

type taskService struct {
	r repository.TaskRepository
}

func NewTaskService(repository repository.TaskRepository) TaskService {
	return &taskService{r: repository}
}

func (t taskService) Save(ownerID, taskID uuid.UUID, creation *transfer.TaskCreation) (insertedID uuid.UUID, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Duplicate(ownerID, taskID uuid.UUID) (replicaID uuid.UUID, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) FetchByID(ownerID, listID, taskID uuid.UUID) (task *model.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Fetch(ownerID, listID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) FetchFromToday(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) FetchFromTomorrow(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) FetchFromDeferred(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Update(ownerID, listID, taskID uuid.UUID, update *transfer.TaskUpdate) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Reorder(ownerID, listID, taskID uuid.UUID, position uint64) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) SetReminder(ownerID, listID, taskID uuid.UUID, remindAt time.Time) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) SetPriority(ownerID, listID, taskID uuid.UUID, priority types.TaskPriority) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) SetDueDate(ownerID, listID, taskID uuid.UUID, dueDate time.Time) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Complete(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Resume(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Pin(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Unpin(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Move(ownerID, taskID, targetListID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Today(ownerID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Tomorrow(ownerID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Defer(ownerID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Trash(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) RestoreFromTrash(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t taskService) Delete(ownerID, listID, taskID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
