package transfer

import (
	"noda/data/types"
	"time"

	"github.com/google/uuid"
)

/* Transfers a user creation request.  */
type UserCreation struct {
	FirstName  string `json:"first_name" validate:"required"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name" validate:"required"`
	Surname    string `json:"surname"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
}

func (u *UserCreation) Validate() error {
	return validate(u)
}

/* Transfers a user update request.  */
type UserUpdate struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Surname    string `json:"surname"`
}

func (u *UserUpdate) Validate() error {
	return validate(u)
}

/* Transfers a user response without a password.  A raw user.   */
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
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

/* Transfers the credentials for a user to sign in.  */
type UserCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *UserCredentials) Validate() error {
	return validate(u)
}
