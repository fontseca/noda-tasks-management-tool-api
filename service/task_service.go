package service

import (
	"github.com/google/uuid"
	"log"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/failure"
	"noda/repository"
	"time"
)

type TaskService interface {
	Save(ownerID, listID uuid.UUID, creation *transfer.TaskCreation) (insertedID uuid.UUID, err error)
	Duplicate(ownerID, taskID uuid.UUID) (replicaID uuid.UUID, err error)
	FetchByID(ownerID, listID, taskID uuid.UUID) (task *model.Task, err error)
	Fetch(ownerID, listID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error)
	FetchFromToday(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error)
	FetchFromTomorrow(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error)
	FetchFromDeferred(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error)
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
		err = failure.NewNilParameterError("Save", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Save", "listID")
		log.Println(err)
		return uuid.Nil, err
	case nil == creation:
		err = failure.NewNilParameterError("Save", "creation")
		log.Println(err)
		return uuid.Nil, err
	case 128 < len(creation.Title):
		return uuid.Nil, failure.ErrTooLong.Clone().FormatDetails("Title", "creation", 128)
	case 64 < len(creation.Headline):
		return uuid.Nil, failure.ErrTooLong.Clone().FormatDetails("Headline", "creation", 64)
	case 512 < len(creation.Description):
		return uuid.Nil, failure.ErrTooLong.Clone().FormatDetails("Description", "creation", 512)
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
		err = failure.NewNilParameterError("Duplicate", "ownerID")
		log.Println(err)
		return uuid.Nil, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Duplicate", "taskID")
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
		err = failure.NewNilParameterError("FetchByID", "ownerID")
		log.Println(err)
		return nil, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("FetchByID", "listID")
		log.Println(err)
		return nil, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("FetchByID", "taskID")
		log.Println(err)
		return nil, err
	}
	return t.r.FetchByID(ownerID.String(), listID.String(), taskID.String())
}

func (t *taskService) Fetch(ownerID, listID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Fetch", "ownerID")
		log.Println(err)
		return nil, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Fetch", "listID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = failure.NewNilParameterError("Fetch", "pagination")
		log.Println(err)
		return nil, err
	}
	doDefaultPagination(pagination)
	doTrim(&needle, &sortExpr)
	tasks, err := t.r.Fetch(ownerID.String(), listID.String(), pagination.Page, pagination.RPP, needle, sortExpr)
	if nil != err {
		return nil, err
	}
	result = &types.Result[model.Task]{
		Page:      pagination.Page,
		RPP:       pagination.RPP,
		Retrieved: int64(len(tasks)),
		Payload:   tasks,
	}
	return result, nil
}

func (t *taskService) FetchFromToday(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("FetchFromToday", "ownerID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = failure.NewNilParameterError("FetchFromToday", "pagination")
		log.Println(err)
		return nil, err
	}
	doDefaultPagination(pagination)
	doTrim(&needle, &sortExpr)
	tasks, err := t.r.FetchFromToday(ownerID.String(), pagination.Page, pagination.RPP, needle, sortExpr)
	if nil != err {
		return nil, err
	}
	result = &types.Result[model.Task]{
		Page:      pagination.Page,
		RPP:       pagination.RPP,
		Retrieved: int64(len(tasks)),
		Payload:   tasks,
	}
	return result, nil
}

func (t *taskService) FetchFromTomorrow(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("FetchFromTomorrow", "ownerID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = failure.NewNilParameterError("FetchFromTomorrow", "pagination")
		log.Println(err)
		return nil, err
	}
	doDefaultPagination(pagination)
	doTrim(&needle, &sortExpr)
	tasks, err := t.r.FetchFromTomorrow(ownerID.String(), pagination.Page, pagination.RPP, needle, sortExpr)
	if nil != err {
		return nil, err
	}
	result = &types.Result[model.Task]{
		Page:      pagination.Page,
		RPP:       pagination.RPP,
		Retrieved: int64(len(tasks)),
		Payload:   tasks,
	}
	return result, nil
}

func (t *taskService) FetchFromDeferred(ownerID uuid.UUID, pagination *types.Pagination, needle, sortExpr string) (result *types.Result[model.Task], err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("FetchFromDeferred", "ownerID")
		log.Println(err)
		return nil, err
	case nil == pagination:
		err = failure.NewNilParameterError("FetchFromDeferred", "pagination")
		log.Println(err)
		return nil, err
	}
	doDefaultPagination(pagination)
	doTrim(&needle, &sortExpr)
	tasks, err := t.r.FetchFromDeferred(ownerID.String(), pagination.Page, pagination.RPP, needle, sortExpr)
	if nil != err {
		return nil, err
	}
	result = &types.Result[model.Task]{
		Page:      pagination.Page,
		RPP:       pagination.RPP,
		Retrieved: int64(len(tasks)),
		Payload:   tasks,
	}
	return result, nil
}

func (t *taskService) Update(ownerID, listID, taskID uuid.UUID, update *transfer.TaskUpdate) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Update", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Update", "listID")
		log.Println(err)
		return false, err
	case nil == update:
		err = failure.NewNilParameterError("Update", "update")
		log.Println(err)
		return false, err
	case 128 < len(update.Title):
		return false, failure.ErrTooLong.Clone().FormatDetails("Title", "update", 128)
	case 64 < len(update.Headline):
		return false, failure.ErrTooLong.Clone().FormatDetails("Headline", "update", 64)
	case 512 < len(update.Description):
		return false, failure.ErrTooLong.Clone().FormatDetails("Description", "update", 512)
	}
	doTrim(&update.Title, &update.Headline, &update.Description)
	return t.r.Update(ownerID.String(), listID.String(), taskID.String(), update)
}

func (t *taskService) Reorder(ownerID, listID, taskID uuid.UUID, position uint64) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Reorder", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Reorder", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Reorder", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.Reorder(ownerID.String(), listID.String(), taskID.String(), position)
}

func (t *taskService) SetReminder(ownerID, listID, taskID uuid.UUID, remindAt time.Time) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("SetReminder", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("SetReminder", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("SetReminder", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.SetReminder(ownerID.String(), listID.String(), taskID.String(), remindAt)
}

func (t *taskService) SetPriority(ownerID, listID, taskID uuid.UUID, priority types.TaskPriority) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("SetPriority", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("SetPriority", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("SetPriority", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.SetPriority(ownerID.String(), listID.String(), taskID.String(), priority)
}

func (t *taskService) SetDueDate(ownerID, listID, taskID uuid.UUID, dueDate time.Time) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("SetDueDate", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("SetDueDate", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("SetDueDate", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.SetDueDate(ownerID.String(), listID.String(), taskID.String(), dueDate)
}

func (t *taskService) Complete(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Complete", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Complete", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Complete", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.Complete(ownerID.String(), listID.String(), taskID.String())
}

func (t *taskService) Resume(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Resume", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Resume", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Resume", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.Resume(ownerID.String(), listID.String(), taskID.String())
}

func (t *taskService) Pin(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Pin", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Pin", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Pin", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.Pin(ownerID.String(), listID.String(), taskID.String())
}

func (t *taskService) Unpin(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Unpin", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Unpin", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Unpin", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.Unpin(ownerID.String(), listID.String(), taskID.String())
}

func (t *taskService) Move(ownerID, taskID, targetListID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Move", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Move", "taskID")
		log.Println(err)
		return false, err
	case uuid.Nil == targetListID:
		err = failure.NewNilParameterError("Move", "targetListID")
		log.Println(err)
		return false, err
	}
	return t.r.Move(ownerID.String(), taskID.String(), targetListID.String())
}

func (t *taskService) Today(ownerID, taskID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Today", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Today", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.Today(ownerID.String(), taskID.String())
}

func (t *taskService) Tomorrow(ownerID, taskID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Tomorrow", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Tomorrow", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.Tomorrow(ownerID.String(), taskID.String())
}

func (t *taskService) Defer(ownerID, taskID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Defer", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Defer", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.Defer(ownerID.String(), taskID.String())
}

func (t *taskService) Trash(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Trash", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Trash", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Trash", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.Trash(ownerID.String(), listID.String(), taskID.String())
}

func (t *taskService) RestoreFromTrash(ownerID, listID, taskID uuid.UUID) (ok bool, err error) {
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("RestoreFromTrash", "ownerID")
		log.Println(err)
		return false, err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("RestoreFromTrash", "listID")
		log.Println(err)
		return false, err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("RestoreFromTrash", "taskID")
		log.Println(err)
		return false, err
	}
	return t.r.RestoreFromTrash(ownerID.String(), listID.String(), taskID.String())
}

func (t *taskService) Delete(ownerID, listID, taskID uuid.UUID) error {
	var err error
	switch {
	case uuid.Nil == ownerID:
		err = failure.NewNilParameterError("Delete", "ownerID")
		log.Println(err)
		return err
	case uuid.Nil == listID:
		err = failure.NewNilParameterError("Delete", "listID")
		log.Println(err)
		return err
	case uuid.Nil == taskID:
		err = failure.NewNilParameterError("Delete", "taskID")
		log.Println(err)
		return err
	}
	return t.r.Delete(ownerID.String(), listID.String(), taskID.String())
}
