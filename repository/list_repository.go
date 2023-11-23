package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/lib/pq"
	"log"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
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
	query := `SELECT make_list ($1, $2, $3, $4);`
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

func (r *listRepository) FetchByID(ownerID, groupID, listID string) (list *model.List, err error) {
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
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = noda.ErrGroupNotFound
			case isNonexistentListError(pqerr):
				err = noda.ErrListNotFound
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = noda.ErrDeadlineExceeded
		} else {
			log.Println(err)
		}
		return nil, err
	}
	return
}

func (r *listRepository) GetTodayListID(ownerID string) (listID string, err error) {
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
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = noda.ErrDeadlineExceeded
		}
	} else {
		log.Println(err)
	}
	return
}

func (r *listRepository) GetTomorrowListID(ownerID string) (listID string, err error) {
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
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			}
		} else if isContextDeadlineError(err) {
			log.Println(err)
			err = noda.ErrDeadlineExceeded
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
	result, err := r.db.QueryContext(ctx, query, ownerID, groupID, page, rpp, needle, sortExpr)
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
	result, err := r.db.QueryContext(ctx, query, ownerID, page, rpp, needle, sortExpr)
	if nil != err {
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
	lists = make([]*model.List, 0)
	err = sqlscan.ScanAll(&lists, result)
	if nil != err {
		log.Println(err)
		lists = nil
	}
	return
}

func (r *listRepository) Remove(ownerID, groupID, listID string) (ok bool, err error) {
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
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = noda.ErrGroupNotFound
			case isNonexistentListError(pqerr):
				err = noda.ErrListNotFound
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

func (r *listRepository) Duplicate(ownerID, listID string) (replicaID string, err error) {
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
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = noda.ErrGroupNotFound
			case isNonexistentListError(pqerr):
				err = noda.ErrListNotFound
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

func (r *listRepository) Scatter(ownerID, listID string) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT convert_to_scattered_list ($1, $2);`
	result := r.db.QueryRowContext(ctx, query, ownerID, listID)
	err = result.Scan(&ok)
	if err != nil {
		var pqerr *pq.Error
		if errors.As(err, &pqerr) {
			switch {
			default:
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			case isNonexistentListError(pqerr):
				err = noda.ErrListNotFound
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

func (r *listRepository) Move(ownerID, listID, targetGroupID string) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT move_list ($1, $2, $3);`
	result := r.db.QueryRowContext(ctx, query, ownerID, listID, targetGroupID)
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
			case isNonexistentListError(pqerr):
				err = noda.ErrListNotFound
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

func (r *listRepository) Update(ownerID, groupID, listID string, update *transfer.ListUpdate) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `SELECT update_list ($1, $2, $3, $4, $5);`
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
				log.Println(noda.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = noda.ErrUserNoLongerExists
			case isNonexistentGroupError(pqerr):
				err = noda.ErrGroupNotFound
			case isNonexistentListError(pqerr):
				err = noda.ErrListNotFound
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
