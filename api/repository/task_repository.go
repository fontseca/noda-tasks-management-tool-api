package repository

import (
	"database/sql"
	"errors"
	"log"
	"noda/api/data/model"
	"noda/failure"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) GetByID(id uuid.UUID) (*model.Task, error) {
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
			log.Println(failure.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer row.Close()

	task := model.Task{}
	if err = sqlscan.ScanOne(&task, row); err != nil {
		if sqlscan.NotFound(err) {
			return nil, failure.ErrNotFound
		}
		return nil, err
	}
	return &task, err
}

func (r *TaskRepository) GetAll() (*[]*model.Task, error) {
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
			log.Println(failure.PQErrorToString(pqerr))
		}
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task
	if err = sqlscan.ScanAll(&tasks, rows); err != nil {
		log.Println(err)
		return nil, err
	}
	return &tasks, nil
}
