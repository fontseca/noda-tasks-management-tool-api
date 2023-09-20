package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

/* Represents a logical steps to follow to complete a task.  */
type Step struct {
	ID          uuid.UUID  `json:"step_id"`
	TaskID      uuid.UUID  `json:"task_id"`
	Order       uint64     `json:"order"`
	Description string     `json:"description"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (s *Step) String() string {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Printf("could not convert step object into string: %s", err)
		return ""
	}
	return string(bytes)
}
