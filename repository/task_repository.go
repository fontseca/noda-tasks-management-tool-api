package repository

import (
	"database/sql"
	"errors"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"noda"
	"noda/data/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) SelectByID(id uuid.UUID) (*model.Task, error) {
	query := `
	SELECT "task_id" AS "id",
	       "group_id",
	       "owner_id",
	       "list_id",
	       "position_in_list",
	       "title",
	       "headline",
	       "description",
	       "priority",
	       "status",
	       "is_pinned",
	       "is_archived",
	       "due_date",
	       "remind_at",
	       "completed_at",
	       "archived_at",
	       "created_at",
	       "updated_at"
	  FROM "task"
	 WHERE "task_id" = $1;`
	row, err := r.db.Query(query, id.String())
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer row.Close()

	task := model.Task{}
	if err = sqlscan.ScanOne(&task, row); err != nil {
		if sqlscan.NotFound(err) {
			return nil, noda.ErrUserNotFound
		}
		return nil, err
	}
	return &task, err
}

func (r *TaskRepository) SelectAll() (*[]*model.Task, error) {
	query := `
	SELECT "task_id" AS "id",
	       "group_id",
	       "owner_id",
	       "list_id",
	       "position_in_list",
	       "title",
	       "headline",
	       "description",
	       "priority",
	       "status",
	       "is_pinned",
	       "is_archived",
	       "due_date",
	       "remind_at",
	       "completed_at",
	       "archived_at",
	       "created_at",
	       "updated_at"
	  FROM "task";`
	rows, err := r.db.Query(query)
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()

	tasks := []*model.Task{}
	if err = sqlscan.ScanAll(&tasks, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return &tasks, nil
}

func (r *TaskRepository) SelectByOwnerID(userID uuid.UUID) (*[]*model.Task, error) {
	query := `
	SELECT "task_id" AS "id",
	       "group_id",
	       "owner_id",
	       "list_id",
	       "position_in_list",
	       "title",
	       "headline",
	       "description",
	       "priority",
	       "status",
	       "is_pinned",
	       "is_archived",
	       "due_date",
	       "remind_at",
	       "completed_at",
	       "archived_at",
	       "created_at",
	       "updated_at"
	  FROM "task"
	 WHERE "owner_id" = $1;`
	rows, err := r.db.Query(query, userID.String())
	if err != nil {
		var pqerr *pq.Error
		switch {
		default:
			log.Println(err)
		case errors.As(err, &pqerr):
			log.Println(noda.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()

	tasks := []*model.Task{}
	if err = sqlscan.ScanAll(&tasks, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return &tasks, err
}
