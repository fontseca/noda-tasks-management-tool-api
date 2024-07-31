package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

/* Organizes tasks under a single unit.  */
type List struct {
	UUID        uuid.UUID `json:"list_uuid"`
	OwnerUUID   uuid.UUID `json:"owner_uuid"`
	GroupUUID   uuid.UUID `json:"group_uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (l *List) String() string {
	bytes, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		log.Printf("could not convert list object into string: %s", err)
		return ""
	}
	return string(bytes)
}
