package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/lib/pq"
	"log"
	"noda/data/model"
	"noda/data/transfer"
	"noda/failure"
	"strings"
	"time"
)

type ListRepository interface {
	Save(ownerID, groupID string, creation *transfer.ListCreation) (insertedID string, err error)
	GetTodayListID(ownerID string) (listID string, err error)
	GetTomorrowListID(ownerID string) (listID string, err error)
	FetchByID(ownerID, groupID, listID string) (list *model.List, err error)
	Fetch(ownerID string, page, rpp int64, needle, sortExpr string) (lists []*model.List, err error)
	FetchGrouped(ownerID, groupID string, page, rpp int64, needle, sortExpr string) (lists []*model.List, err error)
	FetchScattered(ownerID string, page, rpp int64, needle, sortExpr string) (lists []*model.List, err error)
	Update(ownerID, groupID, listID string, update *transfer.ListUpdate) (ok bool, err error)
	Duplicate(ownerID, listID string) (replicaID string, err error)
	Move(ownerID, listID, targetGroupID string) (ok bool, err error)
	Scatter(ownerID, listID string) (ok bool, err error)
	Remove(ownerID, groupID, listID string) (ok bool, err error)
}

type listRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) ListRepository {
	return &listRepository{db}
}

func (r *listRepository) Save(ownerID, groupID string, creation *transfer.ListCreation) (insertedID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT "lists"."make" ($1, $2, $3, $4);`
	var row *sql.Row
	if groupID != "" {
		row = r.db.QueryRowContext(ctx, query, ownerID, groupID, creation.Name, creation.Description)
	} else {
		row = r.db.QueryRowContext(ctx, query, ownerID, nil, creation.Name, creation.Description)
	}
	err = row.Scan(&insertedID)
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
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *listRepository) FetchByID(ownerID, groupID, listID string) (list *model.List, err error) {
	query := `SELECT * FROM "lists"."fetch" (p_owner_uuid := $1,
                                           p_group_uuid := $2,
                                           p_list_uuid := $3,
                                           p_needle := NULL,
                                           p_page := NULL,
                                           p_rpp := NULL);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := r.db.QueryRowContext(ctx, query, ownerID, groupID, listID)
	list = &model.List{}
	err = result.Scan(
		&list.UUID, &list.OwnerUUID, &list.GroupUUID, &list.Name, &list.Description, &list.CreatedAt, &list.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = failure.ErrListNotFound
		} else {
			var pqerr *pq.Error
			if errors.As(err, &pqerr) {
				switch {
				default:
					log.Println(failure.PQErrorToString(pqerr))
				case isNonexistentUserError(pqerr):
					err = failure.ErrUserNoLongerExists
				case isNonexistentGroupError(pqerr):
					err = failure.ErrGroupNotFound
				case isNonexistentListError(pqerr):
				}
			} else {
				log.Println(err)
			}
		}
		return nil, err
	}
	return list, nil
}

func (r *listRepository) GetTodayListID(ownerID string) (listID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT "lists"."get_today_list_uuid" ($1);`
	result := r.db.QueryRowContext(ctx, query, ownerID)
	err = result.Scan(&listID)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrUserNoLongerExists
			}
		}
	} else {
		log.Println(err)
	}
	return
}

func (r *listRepository) GetTomorrowListID(ownerID string) (listID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT "lists"."get_tomorrow_list_uuid" ($1);`
	result := r.db.QueryRowContext(ctx, query, ownerID)
	err = result.Scan(&listID)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrUserNoLongerExists
			}
		}
	} else {
		log.Println(err)
	}
	return
}

func (r *listRepository) Fetch(
	ownerID string,
	page, rpp int64,
	needle, sortExpr string,
) (lists []*model.List, err error) {
	query := `SELECT "list_uuid" AS "uuid",
		               "owner_uuid",
		               "group_uuid",
		               "name",
		               coalesce ("description", '') AS "description",
		               "created_at",
		               "updated_at"
              FROM "lists"."fetch" (p_owner_uuid := $1,
                                    p_group_uuid := NULL,
                                    p_list_uuid := NULL,
                                    p_needle := $3,
                                    p_page := $4,
                                    p_rpp := $5);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := r.db.QueryContext(ctx, query, ownerID, needle, page, rpp)
	if err != nil {
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
		return
	}
	defer result.Close()
	lists = make([]*model.List, 0)
	err = sqlscan.ScanAll(&lists, result)
	if err != nil {
		log.Println(err)
		lists = nil
	}
	return
}

func (r *listRepository) FetchGrouped(
	ownerID, groupID string,
	page, rpp int64,
	needle, sortExpr string,
) (lists []*model.List, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
    SELECT "list_uuid" AS "uuid",
           "owner_uuid",
           "group_uuid",
		       "name",
		       coalesce ("description", '') AS "description",
		       "created_at",
		       "updated_at"
      FROM "lists"."fetch" (p_owner_uuid := $1,
                            p_group_uuid := $2,
                            p_list_uuid := NULL,
                            p_needle := $3,
                            p_page := $4,
                            p_rpp := $5);`
	result, err := r.db.QueryContext(ctx, query, ownerID, groupID, needle, page, rpp)
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
		} else {
			log.Println(err)
		}
		return
	}
	defer result.Close()
	lists = make([]*model.List, 0)
	err = sqlscan.ScanAll(&lists, result)
	if nil != err {
		log.Println(err)
		lists = nil
	}
	return
}

func (r *listRepository) FetchScattered(
	ownerID string,
	page, rpp int64,
	needle, sortExpr string,
) (lists []*model.List, err error) {
	query := `
	SELECT "list_uuid" AS "uuid",
         "owner_uuid",
         "group_uuid",
		     "name",
		     coalesce ("description", '') AS "description",
		     "created_at",
		     "updated_at"
    FROM "lists"."fetch" (p_owner_uuid := $1,
                          p_group_uuid := NULL,
                          p_list_uuid := NULL,
                          p_needle := $2,
                          p_page := $3,
                          p_rpp := $4);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := r.db.QueryContext(ctx, query, ownerID, needle, page, rpp)
	if nil != err {
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
		return
	}
	defer result.Close()
	lists = make([]*model.List, 0)
	err = sqlscan.ScanAll(&lists, result)
	if nil != err {
		log.Println(err)
		lists = nil
	}
	return
}

func (r *listRepository) Remove(ownerID, groupID, listID string) (ok bool, err error) {
	query := `SELECT "lists"."delete" ($1, $2, $3);`
	var result *sql.Row
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if "" == strings.Trim(groupID, " ") {
		result = r.db.QueryRowContext(ctx, query, ownerID, nil, listID)
	} else {
		result = r.db.QueryRowContext(ctx, query, ownerID, groupID, listID)
	}
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
			case isNonexistentListError(pqerr):
				err = failure.ErrListNotFound
			}
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *listRepository) Duplicate(ownerID, listID string) (replicaID string, err error) {
	query := `SELECT "lists"."duplicate" ($1, $2);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := r.db.QueryRowContext(ctx, query, ownerID, listID)
	err = result.Scan(&replicaID)
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
			case isNonexistentListError(pqerr):
				err = failure.ErrListNotFound
			}
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *listRepository) Scatter(ownerID, listID string) (ok bool, err error) {
	query := `SELECT "lists"."convert_to_scattered_list" ($1, $2);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := r.db.QueryRowContext(ctx, query, ownerID, listID)
	err = result.Scan(&ok)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				err = failure.ErrListNotFound
			}
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *listRepository) Move(ownerID, listID, targetGroupID string) (ok bool, err error) {
	query := `SELECT "lists"."move" ($1, $2, $3);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := r.db.QueryRowContext(ctx, query, ownerID, listID, targetGroupID)
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
			case isNonexistentListError(pqerr):
				err = failure.ErrListNotFound
			}
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *listRepository) Update(ownerID, groupID, listID string, update *transfer.ListUpdate) (ok bool, err error) {
	query := `SELECT "lists"."update" ($1, $2, $3, $4, $5);`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var row *sql.Row
	if "" != strings.Trim(groupID, " ") {
		row = r.db.QueryRowContext(ctx, query, ownerID, groupID, listID, update.Name, update.Description)
	} else {
		row = r.db.QueryRowContext(ctx, query, ownerID, nil, listID, update.Name, update.Description)
	}
	err = row.Scan(&ok)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			case isNonexistentListError(pqerr):
				err = failure.ErrListNotFound
			}
		} else {
			log.Println(err)
		}
	}
	return
}
