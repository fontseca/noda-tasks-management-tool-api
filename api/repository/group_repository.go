package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"noda/api/data/transfer"
	"noda/failure"
	"time"

	"github.com/lib/pq"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db}
}

type gr = GroupRepository

func (r *gr) InsertGroup(
	ownerID string,
	newGroup *transfer.GroupCreation) (insertedID string, err error) {
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
				log.Println(failure.PQErrorToString(pqerr))
			case isNonexistentUserError(pqerr):
				err = failure.ErrNotFound
			}
		} else {
			log.Println(err)
		}
	}
	return
}
