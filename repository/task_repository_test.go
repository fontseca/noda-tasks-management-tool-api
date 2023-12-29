package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"noda/data/transfer"
	"noda/data/types"
	"regexp"
	"testing"
)

const taskID = "f8d5b3a2-80f0-4460-bc40-2762141ffc06"

func TestTaskRepository_Save(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r        = NewTaskRepository(db)
		query    = regexp.QuoteMeta(`SELECT make_task ($1, $2, $3);`)
		creation = &transfer.TaskCreation{
			Title:       "task title",
			Description: "task description",
			Headline:    "task headline",
			Priority:    types.TaskPriorityMedium,
			Status:      types.TaskStatusIncomplete,
		}
		res string
		err error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID,
				fmt.Sprintf("ROW('%s', '%s', '%s', '%s', '%s', %s, %s)",
					creation.Title, creation.Headline, creation.Description, creation.Priority, creation.Status, "NULL", "NULL")).
			WillReturnRows(sqlmock.
				NewRows([]string{"make_task"}).
				AddRow(taskID))
		res, err = r.Save(userID, listID, creation)
		assert.Equal(t, taskID, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Save(userID, listID, creation)
		assert.Error(t, err)
		assert.Equal(t, "", res)
	})
}

func TestTaskRepository_Duplicate(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r         = NewTaskRepository(db)
		query     = regexp.QuoteMeta(`SELECT duplicate_task ($1, $2);`)
		res       string
		err       error
		replicaID = uuid.New().String()
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"duplicate_task"}).
				AddRow(replicaID))
		res, err = r.Duplicate(userID, taskID)
		assert.Equal(t, replicaID, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Duplicate(userID, taskID)
		assert.Error(t, err)
		assert.Equal(t, "", res)
	})
}
