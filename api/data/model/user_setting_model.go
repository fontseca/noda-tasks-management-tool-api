package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

/* Represents system users with their personal information and account details.  */
type UserSetting struct {
	ID        uuid.UUID `json:"user_setting_id"`
	UserID    uuid.UUID `json:"user_id"`
	Key       string    `json:"key"`
	Value     any       `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserSetting) String() string {
	bytes, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		log.Printf("could not convert user setting object into string: %s", err)
		return ""
	}
	return string(bytes)
}
