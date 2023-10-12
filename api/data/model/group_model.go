package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

/* Gathers together one or more lists.  */
type Group struct {
	ID          uuid.UUID `json:"group_id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsArchived  bool      `json:"is_archived"`
	ArchivedAt  time.Time `json:"archived_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (g *Group) String() string {
	bytes, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		log.Printf("could not convert group object into string: %s", err)
		return ""
	}
	return string(bytes)
}
