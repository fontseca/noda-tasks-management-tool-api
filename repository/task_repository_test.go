package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"regexp"
	"testing"
	"time"
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

var taskTableColumns = []string{
	"task_uuid",
	"owner_id",
	"list_uuid",
	"position_in_list",
	"title",
	"headline",
	"description",
	"priority",
	"status",
	"is_pinned",
	"due_date",
	"remind_at",
	"completed_at",
	"created_at",
	"updated_at"}

func TestTaskRepository_FetchByID(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT fetch_task_by_id ($1, $2, $3);`)
		res   *model.Task
		err   error
		task  = &model.Task{
			UUID:           uuid.MustParse(taskID),
			OwnerUUID:      uuid.MustParse(userID),
			ListUUID:       uuid.MustParse(listID),
			PositionInList: 1,
			Title:          "task title",
			Headline:       "task headline",
			Description:    "task description",
			Priority:       types.TaskPriorityHigh,
			Status:         types.TaskStatusComplete,
			IsPinned:       false,
			DueDate:        nil,
			RemindAt:       nil,
			CompletedAt:    nil,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID).
			WillReturnRows(sqlmock.
				NewRows(taskTableColumns).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt))
		res, err = r.FetchByID(userID, listID, taskID)
		assert.Equal(t, task, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.FetchByID(userID, listID, taskID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTaskRepository_Fetch(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT fetch_tasks ($1, $2, $3, $4, $5, $6);`)
		res   []*model.Task
		err   error
		task  = &model.Task{
			UUID:           uuid.MustParse(taskID),
			OwnerUUID:      uuid.MustParse(userID),
			ListUUID:       uuid.MustParse(listID),
			PositionInList: 1,
			Title:          "task title",
			Headline:       "task headline",
			Description:    "task description",
			Priority:       types.TaskPriorityHigh,
			Status:         types.TaskStatusComplete,
			IsPinned:       false,
			DueDate:        nil,
			RemindAt:       nil,
			CompletedAt:    nil,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		tasks = []*model.Task{task, task, task}
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, 1, 10, "", "").
			WillReturnRows(sqlmock.
				NewRows(taskTableColumns).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt))
		res, err = r.Fetch(userID, listID, 1, 10, "", "")
		assert.Equal(t, tasks, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Fetch(userID, listID, 1, 10, "", "")
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTaskRepository_FetchFromToday(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT fetch_tasks_from_today_list ($1, $2, $3, $4, $5);`)
		res   []*model.Task
		err   error
		task  = &model.Task{
			UUID:           uuid.MustParse(taskID),
			OwnerUUID:      uuid.MustParse(userID),
			ListUUID:       uuid.MustParse(listID),
			PositionInList: 1,
			Title:          "task title",
			Headline:       "task headline",
			Description:    "task description",
			Priority:       types.TaskPriorityHigh,
			Status:         types.TaskStatusComplete,
			IsPinned:       false,
			DueDate:        nil,
			RemindAt:       nil,
			CompletedAt:    nil,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		tasks = []*model.Task{task, task, task}
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, 1, 10, "", "").
			WillReturnRows(sqlmock.
				NewRows(taskTableColumns).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt))
		res, err = r.FetchFromToday(userID, 1, 10, "", "")
		assert.Equal(t, tasks, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.FetchFromToday(userID, 1, 10, "", "")
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTaskRepository_FetchFromTomorrow(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT fetch_tasks_from_tomorrow_list ($1, $2, $3, $4, $5);`)
		res   []*model.Task
		err   error
		task  = &model.Task{
			UUID:           uuid.MustParse(taskID),
			OwnerUUID:      uuid.MustParse(userID),
			ListUUID:       uuid.MustParse(listID),
			PositionInList: 1,
			Title:          "task title",
			Headline:       "task headline",
			Description:    "task description",
			Priority:       types.TaskPriorityHigh,
			Status:         types.TaskStatusComplete,
			IsPinned:       false,
			DueDate:        nil,
			RemindAt:       nil,
			CompletedAt:    nil,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		tasks = []*model.Task{task, task, task}
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, 1, 10, "", "").
			WillReturnRows(sqlmock.
				NewRows(taskTableColumns).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt))
		res, err = r.FetchFromTomorrow(userID, 1, 10, "", "")
		assert.Equal(t, tasks, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.FetchFromTomorrow(userID, 1, 10, "", "")
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTaskRepository_FetchFromDeferred(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT fetch_tasks_from_deferred_list ($1, $2, $3, $4, $5);`)
		res   []*model.Task
		err   error
		task  = &model.Task{
			UUID:           uuid.MustParse(taskID),
			OwnerUUID:      uuid.MustParse(userID),
			ListUUID:       uuid.MustParse(listID),
			PositionInList: 1,
			Title:          "task title",
			Headline:       "task headline",
			Description:    "task description",
			Priority:       types.TaskPriorityHigh,
			Status:         types.TaskStatusComplete,
			IsPinned:       false,
			DueDate:        nil,
			RemindAt:       nil,
			CompletedAt:    nil,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		tasks = []*model.Task{task, task, task}
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, 1, 10, "", "").
			WillReturnRows(sqlmock.
				NewRows(taskTableColumns).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt).
				AddRow(task.UUID, task.OwnerUUID, task.ListUUID, task.PositionInList, task.Title, task.Headline, task.Description, task.Priority, task.Status, task.IsPinned, task.DueDate, task.RemindAt, task.CompletedAt, task.CreatedAt, task.UpdatedAt))
		res, err = r.FetchFromDeferred(userID, 1, 10, "", "")
		assert.Equal(t, tasks, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.FetchFromDeferred(userID, 1, 10, "", "")
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTaskRepository_Update(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r        = NewTaskRepository(db)
		query    = regexp.QuoteMeta(`SELECT update_task ($1, $2, $3, $4);`)
		creation = &transfer.TaskUpdate{
			Title:       "task title",
			Description: "task description",
			Headline:    "task headline",
		}
		res bool
		err error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID,
				fmt.Sprintf("ROW('%s', '%s', '%s')",
					creation.Title, creation.Headline, creation.Description)).
			WillReturnRows(sqlmock.
				NewRows([]string{"update_task"}).
				AddRow(true))
		res, err = r.Update(userID, listID, taskID, creation)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Update(userID, listID, taskID, creation)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Reorder(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT reorder_task_in_list ($1, $2, $3, $4);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID, 5).
			WillReturnRows(sqlmock.
				NewRows([]string{"reorder_task_in_list"}).
				AddRow(true))
		res, err = r.Reorder(userID, listID, taskID, 5)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Reorder(userID, listID, taskID, 5)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_SetReminder(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT set_task_reminder_date ($1, $2, $3, $4);`)
		res   bool
		err   error
		tm    = time.Now().Add(5 * time.Hour)
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID, tm).
			WillReturnRows(sqlmock.
				NewRows([]string{"set_task_reminder_date"}).
				AddRow(true))
		res, err = r.SetReminder(userID, listID, taskID, tm)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.SetReminder(userID, listID, taskID, tm)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_SetPriority(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r        = NewTaskRepository(db)
		query    = regexp.QuoteMeta(`SELECT set_task_priority ($1, $2, $3, $4);`)
		res      bool
		err      error
		priority = types.TaskPriorityHigh
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID, priority).
			WillReturnRows(sqlmock.
				NewRows([]string{"set_task_reminder_date"}).
				AddRow(true))
		res, err = r.SetPriority(userID, listID, taskID, priority)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.SetPriority(userID, listID, taskID, priority)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_SetDueDate(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT set_task_due_date ($1, $2, $3, $4);`)
		res   bool
		err   error
		tm    = time.Now().Add(5 * time.Hour)
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID, tm).
			WillReturnRows(sqlmock.
				NewRows([]string{"set_task_due_date"}).
				AddRow(true))
		res, err = r.SetDueDate(userID, listID, taskID, tm)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.SetDueDate(userID, listID, taskID, tm)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Complete(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT set_task_as_completed ($1, $2, $3);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"set_task_due_date"}).
				AddRow(true))
		res, err = r.Complete(userID, listID, taskID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Complete(userID, listID, taskID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Resume(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT set_task_as_uncompleted ($1, $2, $3);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"set_task_as_uncompleted"}).
				AddRow(true))
		res, err = r.Resume(userID, listID, taskID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Resume(userID, listID, taskID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Pin(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT pin_task ($1, $2, $3);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"pin_task"}).
				AddRow(true))
		res, err = r.Pin(userID, listID, taskID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Pin(userID, listID, taskID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Unpin(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT unpin_task ($1, $2, $3);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"unpin_task"}).
				AddRow(true))
		res, err = r.Unpin(userID, listID, taskID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Unpin(userID, listID, taskID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Move(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r            = NewTaskRepository(db)
		query        = regexp.QuoteMeta(`SELECT move_task_from_list ($1, $2, $3);`)
		res          bool
		err          error
		targetListID = uuid.New().String()
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, taskID, targetListID).
			WillReturnRows(sqlmock.
				NewRows([]string{"move_task_from_list"}).
				AddRow(true))
		res, err = r.Move(userID, taskID, targetListID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Move(userID, taskID, targetListID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Today(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT move_task_to_today_list ($1, $2);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"move_task_to_today_list"}).
				AddRow(true))
		res, err = r.Today(userID, taskID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Today(userID, taskID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Tomorrow(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT move_task_to_tomorrow_list ($1, $2);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"move_task_to_tomorrow_list"}).
				AddRow(true))
		res, err = r.Tomorrow(userID, taskID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Tomorrow(userID, taskID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Defer(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT move_task_to_deferred_list ($1, $2);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"move_task_to_deferred_list"}).
				AddRow(true))
		res, err = r.Defer(userID, taskID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Defer(userID, taskID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Trash(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT trash_task ($1, $2, $3);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"trash_task"}).
				AddRow(true))
		res, err = r.Trash(userID, listID, taskID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.Trash(userID, listID, taskID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_RestoreFromTrash(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT restore_task_from_trash ($1, $2, $3);`)
		res   bool
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(userID, listID, taskID).
			WillReturnRows(sqlmock.
				NewRows([]string{"restore_task_from_trash"}).
				AddRow(true))
		res, err = r.RestoreFromTrash(userID, listID, taskID)
		assert.True(t, res)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		res, err = r.RestoreFromTrash(userID, listID, taskID)
		assert.False(t, res)
		assert.Error(t, err)
	})
}

func TestTaskRepository_Delete(t *testing.T) {
	defer beQuiet()()
	db, mock := newMock()
	defer db.Close()
	var (
		r     = NewTaskRepository(db)
		query = regexp.QuoteMeta(`SELECT delete_task ($1, $2, $3);`)
		err   error
	)

	t.Run("success", func(t *testing.T) {
		mock.
			ExpectExec(query).
			WithArgs(userID, listID, taskID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		err = r.Delete(userID, listID, taskID)
		assert.NoError(t, err)
	})

	t.Run("unexpected database error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WillReturnError(&pq.Error{})
		err = r.Delete(userID, listID, taskID)
		assert.Error(t, err)
	})
}
