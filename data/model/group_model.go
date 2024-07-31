package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

/* Gathers together one or more lists.  */
type Group struct {
	UUID        uuid.UUID  `json:"group_uuid"`
	OwnerUUID   uuid.UUID  `json:"owner_uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (g *Group) String() string {
	bytes, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		log.Printf("could not convert group object into string: %s", err)
		return ""
	}
	return string(bytes)
}
