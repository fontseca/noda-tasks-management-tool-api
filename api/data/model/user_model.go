package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"user_id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	Surname    string    `json:"surname"`
	PictureUrl *string   `json:"picture_url"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (u *User) String() string {
	bytes, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		log.Printf("could not convert person object into string: %s", err)
		return ""
	}
	return string(bytes)
}
