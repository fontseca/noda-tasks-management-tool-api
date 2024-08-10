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
	query := `SELECT * FROM "groups"."fetch" (p_owner_uuid := $1,
                                            p_group_uuid := $2,
                                            p_needle := NULL,
                                            p_page := NULL,
                                            p_rpp := NULL);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	group = &model.Group{}
	err = r.db.QueryRowContext(ctx, query, ownerID, groupID).
		Scan(&group.UUID, &group.OwnerUUID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = failure.ErrGroupNotFound
		} else {
			var pqerr *pq.Error
			if errors.As(err, &pqerr) {
				switch {
				default:
					log.Println(failure.PQErrorToString(pqerr))
				case isNonexistentUserError(pqerr):
					err = failure.ErrUserNoLongerExists
				}
			} else {
				log.Println(err)
			}
		}
		return nil, err
	}
	return group, nil
}

func (r *groupRepository) Fetch(ownerID string, page, rpp int64, needle, sortBy string) (groups []*model.Group, err error) {
	query := `
	SELECT "group_uuid" AS "uuid",
         "owner_uuid",
         "name",
         "description",
         "created_at",
         "updated_at"
	  FROM "groups"."fetch" (p_owner_uuid := $1,
	                         p_group_uuid := $2,
	                         p_needle := $3,
	                         p_page := $4,
	                         p_rpp := $5);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := r.db.QueryContext(ctx, query, ownerID, nil, needle, page, rpp)
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
	query := `SELECT "groups"."update" ($1, $2, $3, $4);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
	query := `SELECT "groups"."delete" ($1, $2);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
