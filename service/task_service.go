package service

import (
	"github.com/google/uuid"
	"log"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/repository"
	"time"
)

type TaskService interface {
	Save(ownerID, listID uuid.UUID, creation *transfer.TaskCreation) (insertedID uuid.UUID, err error)
	Duplicate(ownerID, taskID uuid.UUID) (replicaID uuid.UUID, err error)
	FetchByID(ownerID, listID, taskID uuid.UUID) (task *model.Task, err error)
	Fetch(ownerID, listID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks *types.Result[model.Task], err error)
	FetchFromToday(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks *types.Result[model.Task], err error)
	FetchFromTomorrow(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks *types.Result[model.Task], err error)
	FetchFromDeferred(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks *types.Result[model.Task], err error)
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

func (t *taskService) Save(ownerID, listID uuid.UUID, creation *transfer.TaskCreation) (insertedID uuid.UUID, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("Save", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("Save", "listID")
		log.Println(err)
		return uuid.Nil, err
	case nil == creation:
		err = noda.NewNilParameterError("Save", "creation")
		log.Println(err)
		return uuid.Nil, err
	case 128 < len(creation.Title):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("Title", "creation", 128)
	case 64 < len(creation.Headline):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("Headline", "creation", 64)
	case 512 < len(creation.Description):
		return uuid.Nil, noda.ErrTooLong.Clone().FormatDetails("Description", "creation", 512)
	}
	doTrim(&creation.Title, &creation.Headline, &creation.Description)
	if "" == creation.Title {
		creation.Title = "Untitled"
	}
	if "" == creation.Priority {
		creation.Priority = types.TaskPriorityMedium
	}
	if "" == creation.Status {
		creation.Status = types.TaskStatusIncomplete
	}
	inserted, err := t.r.Save(ownerID.String(), listID.String(), creation)
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(inserted)
}

func (t *taskService) Duplicate(ownerID, taskID uuid.UUID) (replicaID uuid.UUID, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("Duplicate", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	case uuid.Nil == taskID:
		err = noda.NewNilParameterError("Duplicate", "taskID")
		log.Println(err)
		return uuid.Nil, err
	}
	replica, err := t.r.Duplicate(ownerID.String(), taskID.String())
	if nil != err {
		return uuid.Nil, err
	}
	return uuid.Parse(replica)
}

func (t *taskService) FetchByID(ownerID, listID, taskID uuid.UUID) (task *model.Task, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = noda.NewNilParameterError("FetchByID", "ownerID")
		log.Println(err)
		return nil, err
	case uuid.Nil == listID:
		err = noda.NewNilParameterError("FetchByID", "listID")
		log.Println(err)
		return nil, err
	case uuid.Nil == taskID:
		err = noda.NewNilParameterError("FetchByID", "taskID")
		log.Println(err)
		return nil, err
	}
	return t.r.FetchByID(ownerID.String(), listID.String(), taskID.String())
}

func (t *taskService) Fetch(ownerID, listID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks *types.Result[model.Task], err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) FetchFromToday(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks *types.Result[model.Task], err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) FetchFromTomorrow(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks *types.Result[model.Task], err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) FetchFromDeferred(ownerID uuid.UUID, page, rpp int64, needle, sortExpr string) (tasks *types.Result[model.Task], err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Update(ownerID, listID, taskID uuid.UUID, update *transfer.TaskUpdate) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Reorder(ownerID, listID, taskID uuid.UUID, position uint64) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) SetReminder(ownerID, listID, taskID uuid.UUID, remindAt time.Time) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) SetPriority(ownerID, listID, taskID uuid.UUID, priority types.TaskPriority) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) SetDueDate(ownerID, listID, taskID uuid.UUID, dueDate time.Time) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Complete(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Resume(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Pin(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Unpin(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Move(ownerID, taskID, targetListID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Today(ownerID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Tomorrow(ownerID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Defer(ownerID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Trash(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) RestoreFromTrash(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *taskService) Delete(ownerID, listID, taskID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
