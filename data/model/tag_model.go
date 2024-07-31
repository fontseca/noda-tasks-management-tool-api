package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

/* Labels and categorizes enhance organization and searchability.  */
type Tag struct {
	ID          uuid.UUID `json:"tag_uuid"`
	OwnerID     uuid.UUID `json:"owner_uuid"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Tag) String() string {
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Printf("could not convert tag object into string: %s", err)
		return ""
	}
	return string(bytes)
}
