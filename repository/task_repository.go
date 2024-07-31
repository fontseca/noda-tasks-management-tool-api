package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"noda/failure"
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
	var query = `SELECT "tasks"."make" ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID,
		fmt.Sprintf("ROW('%s', '%s', '%s', '%s', '%s', %s, %s)",
			creation.Title, creation.Headline, creation.Description, creation.Priority, creation.Status, "NULL", "NULL"))
	err = row.Scan(&insertedID)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return "", failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return "", failure.ErrListNotFound
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
	var query = `SELECT "tasks"."duplicate" ($1, $2);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, taskID)
	err = row.Scan(&replicaID)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return "", failure.ErrUserNoLongerExists
			case isNonexistentTaskError(pqerr):
				return "", failure.ErrTaskNotFound
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
	var query = `SELECT "tasks"."fetch_by_uuid" ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID)
	task = new(model.Task)
	err = row.Scan(
		&task.UUID,
		&task.OwnerUUID,
		&task.ListUUID,
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
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return nil, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return nil, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return nil, failure.ErrTaskNotFound
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
	var query = `SELECT "tasks"."fetch" ($1, $2, $3, $4, $5, $6);`
	rows, err := r.db.QueryContext(ctx, query, ownerID, listID, page, rpp, needle, sortExpr)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return nil, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return nil, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return nil, failure.ErrTaskNotFound
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
			&task.UUID,
			&task.OwnerUUID,
			&task.ListUUID,
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
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."fetch_from_today_list" ($1, $2, $3, $4, $5);`
	rows, err := r.db.QueryContext(ctx, query, ownerID, page, rpp, needle, sortExpr)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return nil, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return nil, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return nil, failure.ErrTaskNotFound
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
			&task.UUID,
			&task.OwnerUUID,
			&task.ListUUID,
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

func (r *taskRepository) FetchFromTomorrow(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."fetch_from_tomorrow_list" ($1, $2, $3, $4, $5);`
	rows, err := r.db.QueryContext(ctx, query, ownerID, page, rpp, needle, sortExpr)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return nil, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return nil, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return nil, failure.ErrTaskNotFound
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
			&task.UUID,
			&task.OwnerUUID,
			&task.ListUUID,
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

func (r *taskRepository) FetchFromDeferred(ownerID string, page, rpp int64, needle, sortExpr string) (tasks []*model.Task, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."fetch_from_deferred_list" ($1, $2, $3, $4, $5);`
	rows, err := r.db.QueryContext(ctx, query, ownerID, page, rpp, needle, sortExpr)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return nil, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return nil, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return nil, failure.ErrTaskNotFound
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
			&task.UUID,
			&task.OwnerUUID,
			&task.ListUUID,
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

func (r *taskRepository) Update(ownerID, listID, taskID string, update *transfer.TaskUpdate) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."update" ($1, $2, $3, $4);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID,
		fmt.Sprintf("ROW('%s', '%s', '%s')", update.Title, update.Headline, update.Description))
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Reorder(ownerID, listID, taskID string, position uint64) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."reorder_task_in_list" ($1, $2, $3, $4);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID, position)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) SetReminder(ownerID, listID, taskID string, remindAt time.Time) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."set_reminder_date" ($1, $2, $3, $4);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID, remindAt)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) SetPriority(ownerID, listID, taskID string, priority types.TaskPriority) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."set_priority" ($1, $2, $3, $4);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID, priority)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) SetDueDate(ownerID, listID, taskID string, dueDate time.Time) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."set_due_date" ($1, $2, $3, $4);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID, dueDate)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Complete(ownerID, listID, taskID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."set_as_completed" ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Resume(ownerID, listID, taskID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."set_as_uncompleted" ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Pin(ownerID, listID, taskID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."pin" ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Unpin(ownerID, listID, taskID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."unpin" ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Move(ownerID, taskID, targetListID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."move_one_from_list" ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, taskID, targetListID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Today(ownerID, taskID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."move_one_to_today_list" ($1, $2);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, taskID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Tomorrow(ownerID, taskID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."move_one_to_tomorrow_list" ($1, $2);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, taskID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Defer(ownerID, taskID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."move_one_to_deferred_list" ($1, $2);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, taskID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Trash(ownerID, listID, taskID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."trash" ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) RestoreFromTrash(ownerID, listID, taskID string) (ok bool, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."restore_from_trash" ($1, $2, $3);`
	var row = r.db.QueryRowContext(ctx, query, ownerID, listID, taskID)
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return false, failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return false, failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return false, failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return false, err
	}
	return ok, nil
}

func (r *taskRepository) Delete(ownerID, listID, taskID string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var query = `SELECT "tasks"."delete" ($1, $2, $3);`
	_, err := r.db.ExecContext(ctx, query, ownerID, listID, taskID)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				return failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				return failure.ErrListNotFound
			case isNonexistentTaskError(pqerr):
				return failure.ErrTaskNotFound
			}
		} else {
			log.Println(err)
		}
		return err
	}
	return nil
}
