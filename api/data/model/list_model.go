package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

/* Organizes tasks under a single unit.  */
type List struct {
	ID          uuid.UUID  `json:"list_id"`
	OwnerID     uuid.UUID  `json:"owner_id"`
	GroupID     uuid.UUID  `json:"group_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsArchived  bool       `json:"is_archived"`
	ArchivedAt  *time.Time `json:"archived_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (l *List) String() string {
	bytes, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		log.Printf("could not convert list object into string: %s", err)
		return ""
	}
	return string(bytes)
}
