package model

import (
	"encoding/json"
	"log"
	"noda/data/types"
	"time"

	"github.com/google/uuid"
)

/* Represents system users with their personal information and account details.  */
type User struct {
	ID         uuid.UUID  `json:"user_id"`
	Role       types.Role `json:"role"`
	FirstName  string     `json:"first_name"`
	MiddleName string     `json:"middle_name"`
	LastName   string     `json:"last_name"`
	Surname    string     `json:"surname"`
	PictureUrl *string    `json:"picture_url"`
	Email      string     `json:"email"`
	IsBlocked  bool       `json:"is_blocked"`
	Password   string     `json:"password"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (u *User) String() string {
	bytes, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		log.Printf("could not convert user object into string: %s", err)
		return ""
	}
	return string(bytes)
}
