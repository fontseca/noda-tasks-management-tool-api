package model

import (
	"encoding/json"
	"log"
	"noda/data/types"
	"time"

	"github.com/google/uuid"
)

/* Manages individual tasks, including titles, descriptions, statuses, etc.  */
type Task struct {
	UUID           uuid.UUID          `json:"task_uuid"`
	OwnerUUID      uuid.UUID          `json:"owner_uuid"`
	ListUUID       uuid.UUID          `json:"list_uuid"`
	PositionInList types.Position     `json:"position_in_list"`
	Title          string             `json:"title"`
	Headline       string             `json:"headline"`
	Description    string             `json:"description"`
	Priority       types.TaskPriority `json:"priority"`
	Status         types.TaskStatus   `json:"status"`
	IsPinned       bool               `json:"is_pinned"`
	DueDate        *time.Time         `json:"due_date"`
	RemindAt       *time.Time         `json:"remind_at"`
	CompletedAt    *time.Time         `json:"completed_at"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

func (t *Task) String() string {
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Printf("could not convert task object into string: %s", err)
		return ""
	}
	return string(bytes)
}
