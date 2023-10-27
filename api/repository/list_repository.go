package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/lib/pq"
	"log"
	"noda/api/data/model"
	"noda/api/data/transfer"
	"noda/failure"
	"strings"
	"time"
)

type ListRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) *ListRepository {
	return &ListRepository{db}
}

func (r *ListRepository) InsertList(
	ownerID, groupID string,
	next *transfer.ListCreation) (insertedID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT make_list ($1, $2, $3, $4);`
	var row *sql.Row
	if groupID != "" {
		row = r.db.QueryRowContext(ctx, query, ownerID, groupID, next.Name, next.Description)
	} else {
		row = r.db.QueryRowContext(ctx, query, ownerID, nil, next.Name, next.Description)
	}
	err = row.Scan(&insertedID)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrNotFound
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			}
		} else if isContextDeadlineError(err) {
			err = failure.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *ListRepository) FetchListByID(ownerID, groupID, listID string) (list *model.List, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT * FROM fetch_list_by_id ($1, $2, $3);`
	var result *sql.Row
	if strings.Trim(groupID, " ") != "" {
		result = r.db.QueryRowContext(ctx, query, ownerID, groupID, listID)
	} else {
		result = r.db.QueryRowContext(ctx, query, ownerID, nil, listID)
	}
	list = &model.List{}
	err = result.Scan(
		&list.ID, &list.OwnerID, &list.GroupID, &list.Name, &list.Description,
		&list.IsArchived, &list.ArchivedAt, &list.CreatedAt, &list.UpdatedAt)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrNotFound
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			case isNonexistentListError(pqerr):
				err = failure.ErrListNotFound
			}
		} else if isContextDeadlineError(err) {
			err = failure.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
		return nil, err
	}
	return
}

func (r *ListRepository) GetTodayListID(ownerID string) (listID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT get_today_list_id ($1);`
	result := r.db.QueryRowContext(ctx, query, ownerID)
	err = result.Scan(&listID)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrNotFound
			}
		} else if isContextDeadlineError(err) {
			err = failure.ErrDeadlineExceeded
		}
	} else {
		log.Println(err)
	}
	return
}

func (r *ListRepository) GetTomorrowListID(ownerID string) (listID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT get_tomorrow_list_id ($1);`
	result := r.db.QueryRowContext(ctx, query, ownerID)
	err = result.Scan(&listID)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrNotFound
			}
		} else if isContextDeadlineError(err) {
			err = failure.ErrDeadlineExceeded
		}
	} else {
		log.Println(err)
	}
	return
}

func (r *ListRepository) FetchLists(
	ownerID string,
	page, rpp int64,
	needle, sortExpr string) (lists []*model.List, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT "list_id" AS "id",
		       "owner_id",
		       "group_id",
		       "name",
		       "description",
		       "is_archived",
		       "archived_at",
		       "created_at",
		       "updated_at"
      FROM fetch_lists ($1, $2, $3, $4, $5);`
	result, err := r.db.QueryContext(ctx, query, ownerID, page, rpp, needle, sortExpr)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrNotFound
			}
		} else if isContextDeadlineError(err) {
			err = failure.ErrDeadlineExceeded
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

func (r *ListRepository) FetchGroupedLists(
	ownerID, groupID string,
	page, rpp int64,
	needle, sortBy string) (lists []*model.List, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
    SELECT "list_id" AS "id",
		       "owner_id",
		       "group_id",
		       "name",
		       "description",
		       "is_archived",
		       "archived_at",
		       "created_at",
		       "updated_at"
      FROM fetch_grouped_lists ($1, $2, $3, $4, $5, $6);`
	result, err := r.db.QueryContext(ctx, query, ownerID, groupID, page, rpp, needle, sortBy)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrNotFound
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			}
		} else if isContextDeadlineError(err) {
			err = failure.ErrDeadlineExceeded
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

func (r *ListRepository) FetchScatteredLists(
	ownerID, groupID string,
	page, rpp int64,
	needle, sortBy string) (lists []*model.List, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
	SELECT "list_id" AS "id",
		     "owner_id",
		     "group_id",
		     "name",
		     "description",
		     "is_archived",
		     "archived_at",
		     "created_at",
		     "updated_at"
    FROM fetch_scattered_lists ($1, $2, $3, $4, $5);`
	result, err := r.db.QueryContext(ctx, query, ownerID, groupID, page, rpp, needle, sortBy)
	if nil != err {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrNotFound
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			}
		} else if isContextDeadlineError(err) {
			err = failure.ErrDeadlineExceeded
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

func (r *ListRepository) DeleteList(ownerID, groupID, listID string) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT delete_list ($1, $2, $3);`
	var result *sql.Row
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
				err = failure.ErrNotFound
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			case isNonexistentListError(pqerr):
				err = failure.ErrListNotFound
			}
		} else if isContextDeadlineError(err) {
			err = failure.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
	}
	return
}

func (r *ListRepository) DuplicateList(ownerID, listID string) (replicaID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT duplicate_list ($1, $2);`
	result := r.db.QueryRowContext(ctx, query, ownerID, listID)
	err = result.Scan(&replicaID)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrNotFound
			case isNonexistentGroupError(pqerr):
				err = failure.ErrGroupNotFound
			case isNonexistentListError(pqerr):
				err = failure.ErrListNotFound
			}
		} else if isContextDeadlineError(err) {
			err = failure.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
	}
	return
}
