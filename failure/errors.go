package failure

import (
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)
