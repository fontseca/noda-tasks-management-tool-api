package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"noda/data/model"
	"noda/data/transfer"
	"noda/failure"
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
	result := r.db.QueryRowContext(ctx, `SELECT "groups"."make" ($1, $2, $3);`,
		ownerID, newGroup.Name, newGroup.Description)
	err = result.Scan(&insertedID)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrUserNoLongerExists
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = failure.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *groupRepository) FetchByID(ownerID, groupID string) (group *model.Group, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT * FROM "groups"."fetch_by_id" ($1, $2);`
	result := r.db.QueryRowContext(ctx, query, ownerID, groupID)
	err = result.Err()
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = failure.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
		return
	}
	group = &model.Group{}
	result.Scan(&group.UUID, &group.OwnerUUID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt)
	return
}

func (r *groupRepository) Fetch(ownerID string, page, rpp int64, needle, sortBy string) (groups []*model.Group, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
	SELECT "group_uuid" AS "uuid",
         "owner_uuid",
         "name",
         "description",
         "created_at",
         "updated_at"
	  FROM "groups"."fetch" ($1, $2, $3, $4, $5);`
	result, err := r.db.QueryContext(ctx, query, ownerID, page, rpp, needle, sortBy)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrUserNoLongerExists
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = failure.ErrDeadlineExceeded
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
	query := `SELECT "groups"."update" ($1, $2, $3, $4);`
	result := r.db.QueryRowContext(ctx, query, ownerID, groupID, up.Name, up.Description)
	err = result.Scan(&ok)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = failure.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *groupRepository) Remove(ownerID, groupID string) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT "groups"."delete" ($1, $2);`
	result := r.db.QueryRowContext(ctx, query, ownerID, groupID)
	err = result.Scan(&ok)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = failure.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
	}
	return
}
