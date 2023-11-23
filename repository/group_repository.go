package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"time"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/lib/pq"
)

type GroupRepository interface {
	Save(ownerID string, creation *transfer.GroupCreation) (insertedID string, err error)
	FetchByID(ownerID, groupID string) (group *model.Group, err error)
	Fetch(ownerID string, page, rpp int64, needle, sortExpr string) (groups []*model.Group, err error)
	Update(ownerID, groupID string, update *transfer.GroupUpdate) (ok bool, err error)
	Remove(ownerID, groupID string) (ok bool, err error)
}

type groupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) GroupRepository {
	return &groupRepository{db}
}

func (r *groupRepository) Save(ownerID string, newGroup *transfer.GroupCreation) (insertedID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := r.db.QueryRowContext(ctx, "SELECT make_group ($1, $2, $3);",
		ownerID, newGroup.Name, newGroup.Description)
	err = result.Scan(&insertedID)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = noda.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *groupRepository) FetchByID(ownerID, groupID string) (group *model.Group, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT * FROM fetch_group_by_id ($1, $2);`
	result := r.db.QueryRowContext(ctx, query, ownerID, groupID)
	err = result.Err()
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = noda.ErrGroupNotFound
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = noda.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
		return
	}
	group = &model.Group{}
	result.Scan(
		&group.ID, &group.OwnerID, &group.Name, &group.Description,
		&group.IsArchived, &group.ArchivedAt, &group.CreatedAt, &group.UpdatedAt)
	return
}

func (r *groupRepository) Fetch(ownerID string, page, rpp int64, needle, sortBy string) (groups []*model.Group, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
	SELECT "group_id" AS "id",
         "owner_id",
         "name",
         "description",
         "is_archived",
         "archived_at",
         "created_at",
         "updated_at"
	  FROM fetch_groups ($1, $2, $3, $4, $5);`
	result, err := r.db.QueryContext(ctx, query, ownerID, page, rpp, needle, sortBy)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = noda.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
		return
	}
	defer result.Close()
	groups = []*model.Group{}
	err = sqlscan.ScanAll(&groups, result)
	if err != nil {
		log.Println(err)
		groups = nil
	}
	return
}

func (r *groupRepository) Update(ownerID, groupID string, up *transfer.GroupUpdate) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT update_group ($1, $2, $3, $4);`
	result := r.db.QueryRowContext(ctx, query, ownerID, groupID, up.Name, up.Description)
	err = result.Scan(&ok)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = noda.ErrGroupNotFound
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = noda.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *groupRepository) Remove(ownerID, groupID string) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT delete_group ($1, $2);`
	result := r.db.QueryRowContext(ctx, query, ownerID, groupID)
	err = result.Scan(&ok)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = noda.ErrGroupNotFound
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = noda.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
	}
	return
}
