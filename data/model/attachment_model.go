package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

/* Abstracts a file attached to a task.  */
type Attachment struct {
	ID        uuid.UUID `json:"attachment_id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	TaskID    uuid.UUID `json:"task_id"`
	FileName  string    `json:"file_name"`
	FileURL   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Attachment) String() string {
	bytes, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		log.Printf("could not convert attachment object into string: %s", err)
		return ""
	}
	return string(bytes)
}
