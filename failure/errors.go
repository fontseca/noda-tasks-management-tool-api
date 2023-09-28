package failure

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrSameEmail        = errors.New("the given email address is already registered")
	ErrIncorrectPassord = errors.New("the given password does not match with stored password")
	ErrPassordTooLong   = errors.New("the given password length exceeds 72 bytes")
	ErrUserBlocked      = errors.New("this user has been blocked")
)
