package failure

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrGroupNotFound     = errors.New("group not found")
	ErrListNotFound      = errors.New("list not found")
	ErrSettingNotFound   = errors.New("user setting not found")
	ErrSameEmail         = errors.New("the given email address is already registered")
	ErrIncorrectPassword = errors.New("the given password does not match with stored password")
	ErrPasswordTooLong   = errors.New("the given password length exceeds 72 bytes")
	ErrUserBlocked       = errors.New("this user has been blocked")
	ErrDeadlineExceeded  = errors.New("context deadline exceeded")
)
