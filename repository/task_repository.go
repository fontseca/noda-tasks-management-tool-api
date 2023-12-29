package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"noda"
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

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Save(ownerID, listID string, creation *transfer.TaskCreation) (insertedID string, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT make_task ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID,
		fmt.Sprintf("ROW('%s', '%s', '%s', '%s', '%s', %s, %s)",
			creation.Title, creation.Headline, creation.Description, creation.Priority, creation.Status, "NULL", "NULL"))
	err = row.Scan(&insertedID)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return "", noda.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return "", noda.ErrListNotFound
			}
		} else {
			log.Println(err)
		}
		return "", err
	}
	return insertedID, nil
}

func (r *taskRepository) Duplicate(ownerID, taskID string) (replicaID string, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT duplicate_task ($1, $2);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, taskID)
	err = row.Scan(&replicaID)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return "", noda.ErrUserNoLongerExists
			case isNonexistentTaskError(pqerr):
				return "", noda.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return "", err
	}
	return replicaID, nil
}

func (r *taskRepository) FetchByID(ownerID, listID, taskID string) (task *model.Task, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT fetch_task_by_id ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID)
	task = new(model.Task)
	err = row.Scan(
		&task.ID,
		&task.OwnerID,
		&task.ListID,
		&task.PositionInList,
		&task.Title,
		&task.Headline,
		&task.Description,
		&task.Priority,
		&task.Status,
		&task.IsPinned,
		&task.DueDate,
		&task.RemindAt,
		&task.CompletedAt,
		&task.CreatedAt,
		&task.UpdatedAt)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return nil, noda.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return nil, noda.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return nil, noda.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return nil, err
	}
	return task, nil
}

func (r *taskRepository) Fetch(ownerID, listID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT fetch_tasks ($1, $2, $3, $4, $5, $6);`
	rows, err := r.db.QueryContext(ctx, query, ownerID, listID, page, rpp, needle, sortExpr)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return nil, noda.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return nil, noda.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return nil, noda.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return nil, err
	}
	tasks = make([]*model.Task, 0)
	for rows.Next() {
		var task = new(model.Task)
		err = rows.Scan(
			&task.ID,
			&task.OwnerID,
			&task.ListID,
			&task.PositionInList,
			&task.Title,
			&task.Headline,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.IsPinned,
			&task.DueDate,
			&task.RemindAt,
			&task.CompletedAt,
			&task.CreatedAt,
			&task.UpdatedAt)
		if nil != err {
			break
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *taskRepository) FetchFromToday(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) FetchFromTomorrow(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) FetchFromDeferred(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Update(ownerID, listID, taskID string, update *transfer.TaskUpdate) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Reorder(ownerID, listID, taskID string, position uint64) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) SetReminder(ownerID, listID, taskID string, remindAt time.Time) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) SetPriority(ownerID, listID, taskID string, priority types.TaskPriority) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) SetDueDate(ownerID, listID, taskID string, dueDate time.Time) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Complete(ownerID, listID, taskID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Resume(ownerID, listID, taskID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Pin(ownerID, listID, taskID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Unpin(ownerID, listID, taskID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Move(ownerID, taskID, targetListID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Today(ownerID, taskID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Tomorrow(ownerID, taskID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Defer(ownerID, taskID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Trash(ownerID, listID, taskID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) RestoreFromTrash(ownerID, listID, taskID string) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *taskRepository) Delete(ownerID, listID, taskID string) error {
	//TODO implement me
	panic("implement me")
}
