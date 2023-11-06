package noda

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"net/http"
	"strings"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrGroupNotFound     = errors.New("group not found")
	ErrListNotFound      = errors.New("list not found")
	ErrSettingNotFound   = errors.New("user setting not found")
	ErrSameEmail         = errors.New("the given email address is already registered")
	ErrIncorrectPassword = errors.New("the given password does not match with stored password")
	ErrPasswordTooLong   = errors.New("the given password length exceeds 72 bytes")
	ErrUserBlocked       = errors.New("this user has been blocked")
	ErrDeadlineExceeded  = errors.New("context deadline exceeded")
)

type ErrorCode string

type Error struct {
	status  int
	code    ErrorCode
	message string
	details string
	hint    string
}

func (e *Error) Error() string {
	return e.details
}

func (e *Error) Clone() *Error {
	return &Error{
		code:    e.code,
		message: e.message,
		details: e.details,
		hint:    e.hint,
		status:  e.status,
	}
}

func (e *Error) Details() string {
	return e.details
}

func (e *Error) SetDetails(details string) *Error {
	e.details = strings.Trim(details, " \n\t")
	return e
}

func (e *Error) FormatDetails(a ...any) *Error {
	e.details = fmt.Sprintf(strings.Trim(e.details, " \n\t"), a...)
	return e
}

func (e *Error) Status() int {
	return e.status
}

func (e *Error) SetStatus(status int) *Error {
	e.status = status
	return e
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) SetMessage(message string) *Error {
	e.message = strings.Trim(message, " \n\t")
	return e
}

func (e *Error) Hint() string {
	return e.hint
}

func (e *Error) SetHint(hint string) *Error {
	e.hint = strings.Trim(hint, " \n\t")
	return e
}

func (a *Aggregation) Has() bool {
	return len(a.errors) > 0
}

func PQErrorToString(err *pq.Error) string {
	if err.Hint == "" {
		err.Hint = "(none)"
	}
	if err.Detail == "" {
		err.Detail = "(none)"
	}
	return fmt.Sprintf("postgres driver failed with error \033[1;31m%s\033[0m (%s):\n"+
		"  message: \033[0;33m%s\033[0m\n"+
		"   detail: %s\n"+
		"     hint: %s",
		err.Code, err.Code.Name(), err.Message, err.Detail, err.Hint)
}

type Response struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func NewResponse(message string, details any) *Response {
	err, ok := details.(error)
	if ok {
		return &Response{
			Message: message,
			Details: err.Error(),
		}
	}

	return &Response{
		Message: message,
		Details: details,
	}
}

func Emit(
	w http.ResponseWriter,
	status int,
	message string,
	details any,
) {
	response := NewResponse(message, details)
	res, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(res)
}
